package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/team4yf/fpm-iot-go-middleware/internal/core"
	"github.com/team4yf/fpm-iot-go-middleware/internal/model"
	"github.com/team4yf/fpm-iot-go-middleware/pkg"
)

type reqCreateDevice struct {
	sn        string `json:"sn"`
	appId     string `json:"appId"`
	projectId int64  `json:"projectId"`
}

func Load(app *core.App) {
	app.Post("/push/{device}/{brand}/{event}", func(w http.ResponseWriter, r *http.Request) {
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
		// 记录下获取到的数据，用来进行日志查询
		log.Printf("Receive: %s, body: %s\n", r.URL, body)
		// 获取设备的类型+品牌，从配置文件中获取对应的参数信息
		deviceSpecificName := device + "-" + brand
		if !app.Config.IsSet("notify." + deviceSpecificName) {
			log.Printf("event type: %s not set in config.json", deviceSpecificName)
			http.Error(w, "unKnown data source", http.StatusBadRequest)
			return
		}

		// 开启一个新的 goroutine 来获取设备对应的应用id，并发送到 mqtt 消息服务器上
		go func() {
			// 从配置文件中获取到设备平台推送的消息体中的 设备ID 的JsonPath
			devicePath := app.Config.GetConfigOrDefault("notify."+deviceSpecificName+".devicePath", "$.data").(string)
			res, err := pkg.GetJsonPathData(body, devicePath)
			if err != nil {
				log.Printf("device id not: %v", err)
				return
			}
			deviceID := res.(string)
			// 通过设备和id获取到具体对应的项目信息，如果设备不存在或者设备状态不对的话，会抛出异常信息
			uuid, projid, err := app.Service.Receive(device, brand, event, deviceID)
			if err != nil {
				log.Printf("Device Not Exists Or Not Actived: %v", err)
				return
			}
			wrapper := make(map[string]interface{})
			// 添加固定的静态数据，用于应用平台使用
			bind := app.Config.GetMapOrDefault("notify."+deviceSpecificName+".bind", nil)

			wrapper["origin"] = body // 源消息体
			wrapper["event"] = event // 设备事件
			wrapper["aid"] = uuid    // 设备对应的应用服务平台id
			wrapper["pid"] = projid  // 设备对应的在服务中的项目id
			wrapper["sn"] = deviceID // 设备的编码
			wrapper["type"] = device // 设备对应的类型
			wrapper["brand"] = brand // 设备对应的品牌
			wrapper["bind"] = bind   // 设备绑定的静态数据

			j, _ := json.Marshal(wrapper)

			app.PubSub.Publish(fmt.Sprintf("^push/%s/event", uuid), j)

		}()
		// 响应配置文件中的内容
		response := app.Config.GetMapOrDefault("notify."+deviceSpecificName+".response", nil)
		// log.Printf("device: %s brand:%s event:%s body:%s response:%s\n", device, brand, event, body, response)
		app.WriteJSON(w, 200, response)
	})

	app.Post("/device/{type}/{brand}/create", func(w http.ResponseWriter, r *http.Request) {
		// 从接口路径中获取参数
		params := mux.Vars(r)
		device := params["type"]
		brand := params["brand"]

		bodyMap, err := pkg.GetBodyMap(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		log.Printf("body %s", bodyMap)

		data := &model.Device{}
		data.AppID = bodyMap["appId"].(string)
		data.Brand = brand
		data.Type = device
		data.ProjectID = int64(bodyMap["projectId"].(float64))
		data.SN = bodyMap["sn"].(string)
		data.Status = 1
		data.RegisterAt = time.Now()
		data.LastUpdateAt = time.Now()

		if err := app.Service.RegisterDevice(data); err != nil {
			log.Printf("Error register device: %v", err)
			http.Error(w, "can't register device", http.StatusBadRequest)
			return
		}
		app.WriteJSON(w, 200, `{"code": 1}`)

	})
}
