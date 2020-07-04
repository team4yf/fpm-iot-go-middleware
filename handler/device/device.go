package device

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/team4yf/fpm-iot-go-middleware/config"
	"github.com/team4yf/fpm-iot-go-middleware/errno"
	"github.com/team4yf/fpm-iot-go-middleware/internal/core"
	"github.com/team4yf/fpm-iot-go-middleware/internal/model"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/log"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/utils"
)

type messageHeader struct {
	Version   int    `json:"v"`
	NameSpace string `json:"ns"`
	Name      string `json:"name"`
	AppID     string `json:"appId"`
	ProjID    int64  `json:"projId"`
	Source    string `json:"source"`
}

type payloadDevice struct {
	ID      string                 `json:"id"`
	Type    string                 `json:"type"`
	Name    string                 `json:"name"`
	Brand   string                 `json:"brand"`
	Version string                 `json:"v"`
	Extra   map[string]interface{} `json:"x,omitempty"`
}

type messagePayload struct {
	Device    payloadDevice `json:"device"`
	Data      interface{}   `json:"data"`
	Cgi       string        `json:"cgi"`
	Timestamp int64         `json:"timestamp"`
}

type message struct {
	Header  messageHeader  `json:"header"`
	Payload messagePayload `json:"payload"`
}

//PushHandler 推送相关的处理函数
func PushHandler(app *core.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从接口路径中获取参数
		params := mux.Vars(r)
		device := params["device"]
		brand := params["brand"]
		event := params["event"]
		// 获取post的消息体
		body, err := utils.GetBodyString(r.Body)
		if err != nil {
			app.FailWithError(w, err)
			return
		}
		// 记录下获取到的数据，用来进行日志查询
		log.Debugf("Receive: %s, body: %s\n", r.URL, body)
		// 获取设备的类型+品牌，从配置文件中获取对应的参数信息
		deviceSpecificName := device + "-" + brand
		if !config.IsSet("notify." + deviceSpecificName) {
			app.Fail(w, "unKnown data source")
			return
		}

		// 开启一个新的 goroutine 来获取设备对应的应用id，并发送到 mqtt 消息服务器上
		go func() {
			// 从配置文件中获取到设备平台推送的消息体中的 设备ID 的JsonPath
			devicePath := config.GetConfigOrDefault("notify."+deviceSpecificName+".devicePath", "$.data").(string)
			log.Debugf("jsonPath: %s", devicePath)
			res, err := utils.GetJsonPathData(body, devicePath)
			if err != nil {
				log.Errorf("device id not: %v", err)
				return
			}
			deviceID := res.(string)
			// 通过设备和id获取到具体对应的项目信息，如果设备不存在或者设备状态不对的话，会抛出异常信息
			uuid, projid, err := app.Service.Receive(device, brand, event, deviceID)
			if err != nil {
				log.Errorf("Device Not Exists Or Not Actived: %v", err)
				return
			}

			msgHeader := messageHeader{
				Version:   10,
				NameSpace: "FPM.Lamp." + device,
				Name:      event,
				AppID:     uuid,
				ProjID:    projid,
				Source:    "HTTP",
			}

			// 添加固定的静态数据，用于应用平台使用
			bind := config.GetMapOrDefault("notify."+deviceSpecificName+".bind", nil)
			msgPayloadDevice := payloadDevice{
				ID:      deviceID,
				Type:    device,
				Name:    "-",
				Brand:   brand,
				Version: "v10",
				Extra:   bind,
			}

			msgPayload := messagePayload{
				Device:    msgPayloadDevice,
				Data:      body,
				Cgi:       deviceID,
				Timestamp: time.Now().Unix(),
			}

			msg := message{
				Header:  msgHeader,
				Payload: msgPayload,
			}

			j, _ := json.Marshal(msg)

			app.PubSub.Publish(fmt.Sprintf("$d2s/%s/partner/push", uuid), j)

		}()
		// 响应配置文件中的内容
		response := config.GetMapOrDefault("notify."+deviceSpecificName+".response", nil)
		log.Debugf("device: %s brand:%s event:%s body:%s response:%s\n", device, brand, event, body, response)
		app.SendOk(w, response)
	}
}

//CreateHandler 创建设备相关的
func CreateHandler(app *core.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从接口路径中获取参数
		params := mux.Vars(r)
		device := params["type"]
		brand := params["brand"]

		bodyMap, err := utils.GetBodyMap(r.Body)
		if err != nil {
			app.Fail(w, "can't read body")
			return
		}
		log.Debugf("body %s", bodyMap)

		data := &model.Device{}
		data.AppID = bodyMap["appId"].(string)
		data.Brand = brand
		data.Type = model.DeviceType(device)
		data.ProjectID = int64(bodyMap["projectId"].(float64))
		data.SN = bodyMap["sn"].(string)
		data.Status = 1
		data.RegisterAt = time.Now()
		data.LastUpdateAt = time.Now()

		if err := app.Service.RegisterDevice(data); err != nil {
			log.Errorf("Error register device: %v", err)
			app.Fail(w, "can't register device")
			return
		}
		app.SendOk(w, errno.OK)

	}
}
