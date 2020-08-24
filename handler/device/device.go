package device

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/oliveagle/jsonpath"
	"github.com/spf13/viper"
	msg "github.com/team4yf/fpm-iot-go-middleware/internal/message"
	"github.com/team4yf/fpm-iot-go-middleware/internal/model"
	"github.com/team4yf/fpm-iot-go-middleware/internal/service"
	"github.com/team4yf/yf-fpm-server-go/ctx"
	"github.com/team4yf/yf-fpm-server-go/fpm"
)

var deviceService service.DeviceService
var mux sync.Mutex
var isInited bool

//Init some init work here
func Init() {
	mux.Lock()
	defer mux.Unlock()
	if isInited {
		return
	}
	app := fpm.Default()
	c, exists := app.GetCacher()
	if !exists {
		panic(`Cacher Not Inited!`)
	}
	deviceService = service.NewSimpleDeviceService(c)
	isInited = true
}

//PushHandler 推送相关的处理函数
func PushHandler() func(*ctx.Ctx, *fpm.Fpm) {
	return func(c *ctx.Ctx, fpmApp *fpm.Fpm) {
		// 从接口路径中获取参数
		device := c.Param("device")
		brand := c.Param("brand")
		event := c.Param("event")

		body := &fpm.BizParam{}
		if err := c.ParseBody(&body); err != nil {
			c.Fail(err)
			return
		}
		// 获取设备的类型+品牌，从配置文件中获取对应的参数信息
		deviceSpecificName := device + "-" + brand
		// 记录下获取到的数据，用来进行日志查询
		fpmApp.Logger.Debugf("Receive body: %+v\n deviceSpecificName: %s\n", body, deviceSpecificName)
		if !viper.IsSet("notify." + deviceSpecificName) {
			c.Fail("unKnown data source")
			return
		}

		// 开启一个新的 goroutine 来获取设备对应的应用id，并发送到 mqtt 消息服务器上
		go func() {
			// 从配置文件中获取到设备平台推送的消息体中的 设备ID 的JsonPath
			devicePath := fpmApp.GetConfig("notify." + deviceSpecificName + ".devicePath").(string)
			fpmApp.Logger.Debugf("jsonPath: %s", devicePath)
			res, err := getJSONPathDataFromBiz(body, devicePath)
			if err != nil {
				fpmApp.Logger.Errorf("device id not: %v", err)
				return
			}
			deviceID := res.(string)
			// 通过设备和id获取到具体对应的项目信息，如果设备不存在或者设备状态不对的话，会抛出异常信息
			uuid, projid, err := deviceService.Receive(device, brand, event, deviceID)
			if err != nil {
				fpmApp.Logger.Errorf("Device Not Exists Or Not Actived: %v", err)
				return
			}

			msgHeader := &msg.Header{
				Version:   10,
				NameSpace: "FPM.Lamp." + device,
				Name:      event,
				AppID:     uuid,
				ProjID:    projid,
				Source:    "HTTP",
			}

			// 添加固定的静态数据，用于应用平台使用
			bind := viper.GetStringMap("notify." + deviceSpecificName + ".bind")
			msgPayloadDevice := &msg.Device{
				ID:      deviceID,
				Type:    device,
				Name:    "-",
				Brand:   brand,
				Version: "v10",
				Extra:   bind,
			}

			msgPayload := &msg.D2SPayload{
				Device:    msgPayloadDevice,
				Data:      body,
				Cgi:       deviceID,
				Timestamp: time.Now().Unix(),
			}

			msg := msg.D2SMessage{
				Header:  msgHeader,
				Payload: msgPayload,
			}

			j, _ := json.Marshal(msg)

			//执行发布消息的函数
			fpmApp.Execute("mqttclient.publish", &fpm.BizParam{
				"topic":   fmt.Sprintf("$d2s/%s/partner/push", uuid),
				"payload": j,
			})

		}()
		// 响应配置文件中的内容
		response := fpmApp.GetConfig("notify." + deviceSpecificName + ".response")
		fpmApp.Logger.Debugf("device: %s brand:%s event:%s body:%s response:%s\n", device, brand, event, body, response)
		c.JSON(response)
	}
}

func getJSONPathDataFromBiz(param *fpm.BizParam, jp string) (interface{}, error) {
	bytes, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	var jsonData interface{}
	json.Unmarshal(bytes, &jsonData)
	res, err := jsonpath.JsonPathLookup(jsonData, jp)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//CreateHandler 创建设备相关的
func CreateHandler() func(*ctx.Ctx, *fpm.Fpm) {
	return func(c *ctx.Ctx, fpmApp *fpm.Fpm) {
		// 从接口路径中获取参数
		device := c.Param("type")
		brand := c.Param("brand")

		body := &fpm.BizParam{}
		if err := c.ParseBody(&body); err != nil {
			c.Fail(err)
			return
		}

		data := &model.Device{}
		data.AppID = (*body)["appId"].(string)
		data.Brand = brand
		data.Type = model.DeviceType(device)
		data.ProjectID = int64((*body)["projectId"].(float64))
		data.SN = (*body)["sn"].(string)
		data.Status = 1
		data.RegisterAt = time.Now()
		data.LastUpdateAt = time.Now()

		// TODO: save the device
		if err := deviceService.RegisterDevice(data); err != nil {
			c.Fail("can't register device")
			return
		}
		c.JSON(map[string]int{
			"errno": 0,
		})

	}
}
