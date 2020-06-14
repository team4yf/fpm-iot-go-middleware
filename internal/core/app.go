//  核心的结构体
// 用于响应所有的请求
// 处理 设备平台推送的数据，进行整合之后，将消息转发给 MQTT
package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/team4yf/fpm-iot-go-middleware/config"
	s "github.com/team4yf/fpm-iot-go-middleware/internal/service"
	"github.com/team4yf/fpm-iot-go-middleware/pkg"
	"github.com/team4yf/fpm-iot-go-middleware/router/middleware"
)

type App struct {
	Config     *config.Config
	Router     *mux.Router
	Middleware *middleware.Middleware
	pubSub     pkg.PubSub
	service    s.DeviceService
}

func (app *App) Init(pubSub pkg.PubSub, service s.DeviceService) {
	app.Router = mux.NewRouter()
	app.pubSub = pubSub
	app.service = service
	app.Middleware = &middleware.Middleware{}
	m := alice.New(app.Middleware.LoggerMiddleware, app.Middleware.RecoverMiddleware)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app.Router.Handle("/push/{device}/{brand}/{event}", m.ThenFunc(app.pushHandler)).Methods("POST")
}

func (app *App) Run(addr string) {

	log.Printf("startup %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, app.Router))
}

func (app *App) pushHandler(w http.ResponseWriter, r *http.Request) {
	// 从接口路径中获取参数
	params := mux.Vars(r)
	device := params["device"]
	brand := params["brand"]
	event := params["event"]
	// 获取post的消息体
	body, err := pkg.GetBodyString(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	// 获取设备的类型+品牌，从配置文件中获取对应的参数信息
	deviceSpecificName := device + "-" + brand
	if !app.Config.IsSet("notify." + deviceSpecificName) {
		log.Printf("event type: %s not set in config.json", deviceSpecificName)
		http.Error(w, "unKnown data source", http.StatusBadRequest)
		return
	}

	// 开启一个新的 goroutine 来获取设备对应的应用id，并发送到 mqtt 消息服务器上
	go func() {
		var deviceID string
		devicePath := app.Config.GetConfigOrDefault("notify."+deviceSpecificName+".devicePath", "$.data").(string)
		if res, err := pkg.GetJsonPathData(body, devicePath); err != nil {
			log.Printf("device id not: %v", err)
			return
		} else {
			deviceID = res.(string)
		}
		if uuid, projid, err := app.service.Receive(device, brand, event, deviceID); err != nil {
			log.Printf("Error reading body: %v", err)
			return
		} else {
			wrapper := make(map[string]interface{})
			bind := app.Config.GetMapOrDefault("notify."+deviceSpecificName+".bind", nil)

			wrapper["payload"] = body
			wrapper["event"] = event
			wrapper["uuid"] = uuid
			wrapper["projid"] = projid
			wrapper["deviceID"] = deviceID
			wrapper["device"] = device
			wrapper["brand"] = brand
			wrapper["bind"] = bind

			j, _ := json.Marshal(wrapper)

			app.pubSub.Publish(fmt.Sprintf("^push/%s/event", uuid), j)
		}
	}()
	// 响应配置文件中的内容
	response := app.Config.GetMapOrDefault("notify."+deviceSpecificName+".response", nil)
	// log.Printf("device: %s brand:%s event:%s body:%s response:%s\n", device, brand, event, body, response)
	writeJSON(w, 200, response)
}

func writeJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
