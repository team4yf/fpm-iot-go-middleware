//  核心的结构体
// 用于响应所有的请求
// 处理 设备平台推送的数据，进行整合之后，将消息转发给 MQTT
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

type App struct {
	Config     *Config
	Router     *mux.Router
	Middleware *Middleware
	pubSub     PubSub
	service    Service
}

func (app *App) Init(pubSub PubSub, service Service) {
	app.Router = mux.NewRouter()
	app.pubSub = pubSub
	app.service = service
	app.Middleware = &Middleware{}
	m := alice.New(app.Middleware.LoggerMiddleware, app.Middleware.RecoverMiddleware)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app.Router.Handle("/push/{device}/{brand}/{event}", m.ThenFunc(app.pushHandler)).Methods("POST")
}

func (app *App) Run(addr string) {

	log.Printf("startup %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, app.Router))
}

func (app *App) pushHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	device := params["device"]
	brand := params["brand"]
	event := params["event"]
	body, err := GetBodyString(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	deviceSpecificName := device + "-" + brand
	if !app.Config.IsSet("notify." + deviceSpecificName) {
		log.Printf("event type: %s not set in config.json", deviceSpecificName)
		http.Error(w, "unKnown data source", http.StatusBadRequest)
		return
	}

	go func() {
		var deviceID string
		deviceSpecificName := device + "-" + brand
		devicePath := app.Config.GetConfigOrDefault("notify."+deviceSpecificName+".devicePath", "$.data")
		if res, err := GetJsonPathData(body, devicePath); err != nil {
			log.Printf("device id not: %v", err)
			return
		} else {
			deviceID = res.(string)
		}
		if uuid, err := app.service.Receive(device, brand, event, deviceID); err != nil {
			log.Printf("Error reading body: %v", err)
			return
		} else {
			wrapper := make(map[string]interface{})
			bind := app.Config.GetMapOrDefault("notify."+deviceSpecificName+".bind", nil)

			wrapper["payload"] = body
			wrapper["event"] = event
			wrapper["uuid"] = uuid
			wrapper["deviceID"] = deviceID
			wrapper["device"] = device
			wrapper["brand"] = brand
			wrapper["bind"] = bind

			j, _ := json.Marshal(wrapper)

			app.pubSub.Publish(fmt.Sprintf("^push/%s/event", uuid), j)
		}
	}()
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
