//Package consumer the consumer of the mq, normally receive the iot-app message to device
package consumer

import (
	"strings"
	"time"

	"github.com/team4yf/fpm-iot-go-middleware/external/device/light"
	"github.com/team4yf/fpm-iot-go-middleware/external/rest"
	"github.com/team4yf/fpm-iot-go-middleware/internal/message"
	"github.com/team4yf/fpm-iot-go-middleware/internal/service"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/utils"
	"github.com/team4yf/yf-fpm-server-go/fpm"
)

//DefaultMqttConsumer the default subscriber of the mqtt mq server
//See detail: https://shimo.im/docs/bJaoNiMc4yEfkRSt#anchor-MFbv
func DefaultMqttConsumer(fpmApp *fpm.Fpm) func(interface{}, interface{}) {
	light.Init()
	c, _ := fpmApp.GetCacher()
	deviceService := service.NewSimpleDeviceService(c)
	return func(topic, datastream interface{}) {
		str := (string)(datastream.([]byte))
		fpmApp.Logger.Debugf("received: %s", str)
		var msg message.S2DMessage
		if err := utils.StringToStruct(str, &msg); err != nil {
			fpmApp.Logger.Errorf("convert payload to s2d message fail: error -> %v", err)
			return
		}
		header, bind, payload := msg.Header, msg.Bind, msg.Payload
		fpmApp.Logger.Debugf("header: %v, bind: %v, payload: %v", header, bind, payload)
		if header.Version != 10 {
			fpmApp.Logger.Errorf("only support version: 10")
			return
		}

		namespaces := strings.Split(header.NameSpace, ".")
		if len(namespaces) < 3 {
			fpmApp.Logger.Errorf("error namespace: %s! should like: FPM.Lamp.Light!", header.NameSpace)
			return
		}
		org, industry, spec := strings.ToUpper(namespaces[0]), strings.ToUpper(namespaces[1]), strings.ToUpper(namespaces[2])
		if org != "FPM" {
			fpmApp.Logger.Errorf("only support namespace: FPM, more namespace support plz contact the admin")
			return
		}
		if industry != "LAMP" {
			fpmApp.Logger.Errorf("only support namespace: FPM.LAMP, more namespace support plz contact the admin")
			return
		}
		if spec != "LIGHT" {
			fpmApp.Logger.Errorf("unsupported namespace: FPM.LAMP.%s, plz contact the admin", spec)
			return
		}

		if "Control" != header.Name && "GroupControl" != header.Name {
			fpmApp.Logger.Errorf("only support name: Control or GroupControl")
			return
		}

		appID, projID := header.AppID, header.ProjID
		fpmApp.Logger.Debugf("the appID %s, the projID %d", appID, projID)
		//Cache project exists status.
		setting, err := deviceService.GetSetting(appID, projID)
		if err != nil {
			fpmApp.Logger.Errorf("Get setting of the appID %s, projID %d error: %v", appID, projID, err)
			return
		}
		fpmApp.Logger.Debugf("Get setting of the appID %s, projID %d setting: %v", appID, projID, setting)

		switch spec {
		case "LIGHT":
			lightSetting, ok := setting["light"].(map[string]interface{})
			if !ok {
				fpmApp.Logger.Errorf("no light setting defined of appid: %s, project id: %d ", appID, projID)
				return
			}
			controlLight(lightSetting, header, payload)

		}

	}
}

func feedback(header *message.Header, payload *message.S2DPayload, rsp *rest.APIResponse, err error) {
	fpmApp := fpm.Default()
	fpmApp.Logger.Infof("feedbackStart")
	var result interface{}
	if err != nil {
		result = err.Error
	} else {
		result = rsp
	}

	feedbackBody := &message.D2SFeedback{
		MsgID:     payload.MsgID,
		Timestamp: time.Now().Unix(),
		Cgi:       payload.Cgi,
		Result:    result,
	}

	feedbackMessage := &message.D2SFeedbackMessage{
		Header:   header,
		Feedback: feedbackBody,
	}
	message := utils.JSON2String(feedbackMessage)
	fpmApp.Logger.Infof("feedback: %+v", message)

	fpmApp.Execute("mqttclient.publish", &fpm.BizParam{
		"topic":   "$d2s/" + header.AppID + "/partner/feedback",
		"payload": ([]byte)(message),
	})
}

func controlLight(lightSetting map[string]interface{}, header *message.Header, payloads []*message.S2DPayload) {
	fpmApp := fpm.Default()
	client, err := light.NewAPIClient(lightSetting["brand"].(string), lightSetting["appid"].(string))
	if err != nil {
		fpmApp.Logger.Errorf("cant get the client of brand %s, id %s, error: %+v", lightSetting["brand"].(string), lightSetting["appid"].(string), err)
	}
	//use the appID & projID to find the setting
	//
	// control single deivce
	if header.Name == "Control" {
		payload := payloads[0]
		needFeedback, cmd, args := payload.Feedback, payload.Cmd, payload.Argument
		rsp, err := client.Execute(cmd, args)
		if needFeedback != 0 {
			feedback(header, payload, rsp, err)
		}
		if err != nil {
			fpmApp.Logger.Errorf("execute api control error %+v", err)
			return
		}
		fpmApp.Logger.Infof("Payload:%+v Response: %+v", payload, rsp)
	}
}

//DevicePushConsumer the device push consumer
func DevicePushConsumer(fpmApp *fpm.Fpm) func(interface{}, interface{}) {
	return func(topic, datastream interface{}) {
		str := (string)(datastream.([]byte))
		fpmApp.Logger.Debugf("topic: %v, datastream: %v", topic, str)
	}
}
