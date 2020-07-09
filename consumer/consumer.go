//Package consumer the consumer of the mq, normally receive the iot-app message to device
package consumer

import (
	"strings"

	"github.com/team4yf/fpm-iot-go-middleware/external/device/light"
	"github.com/team4yf/fpm-iot-go-middleware/internal/core"
	"github.com/team4yf/fpm-iot-go-middleware/internal/message"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/log"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/utils"
)

//DefaultMqttConsumer the default subscriber of the mqtt mq server
//See detail: https://shimo.im/docs/bJaoNiMc4yEfkRSt#anchor-MFbv
func DefaultMqttConsumer(app *core.App) func(interface{}, interface{}) {

	light.Init()

	return func(topic, datastream interface{}) {
		str := (string)(datastream.([]byte))
		log.Debugf("received: %s", str)
		var msg message.S2DMessage
		if err := utils.StringToStruct(str, &msg); err != nil {
			log.Errorf("convert payload to s2d message fail: error -> %v", err)
			return
		}
		header, bind, payload := msg.Header, msg.Bind, msg.Payload
		log.Debugf("header: %v, bind: %v, payload: %v", header, bind, payload)
		if header.Version != 10 {
			log.Errorf("only support version: 10")
			return
		}

		namespaces := strings.Split(header.NameSpace, ".")
		if len(namespaces) < 3 {
			log.Errorf("error namespace: %s! should like: FPM.Lamp.Light!", header.NameSpace)
			return
		}
		org, industry, spec := strings.ToUpper(namespaces[0]), strings.ToUpper(namespaces[1]), strings.ToUpper(namespaces[2])
		if org != "FPM" {
			log.Errorf("only support namespace: FPM, more namespace support plz contact the admin")
			return
		}
		if industry != "LAMP" {
			log.Errorf("only support namespace: FPM.LAMP, more namespace support plz contact the admin")
			return
		}
		if spec != "LIGHT" {
			log.Errorf("unsupported namespace: FPM.LAMP.%s, plz contact the admin", spec)
			return
		}

		if "Control" != header.Name && "GroupControl" != header.Name {
			log.Errorf("only support name: Control or GroupControl")
			return
		}

		appID, projID := header.AppID, header.ProjID
		log.Debugf("the appID %s, the projID %d", appID, projID)
		setting, err := app.Service.GetSetting(appID, projID)
		if err != nil {
			log.Errorf("Get setting of the appID %s, projID %d error: %v", appID, projID, err)
			return
		}
		log.Debugf("Get setting of the appID %s, projID %d setting: %v", appID, projID, setting)

		switch spec {
		case "LIGHT":
			lightSetting, ok := setting["light"].(map[string]interface{})
			if !ok {
				log.Errorf("no light setting defined of appid: %s, project id: %d ", appID, projID)
				return
			}
			controlLight(lightSetting, header, payload)

		}

	}
}

func feedback() {

}

func controlLight(lightSetting map[string]interface{}, header *message.Header, payloads []*message.S2DPayload) {

	client, err := light.NewAPIClient(lightSetting["brand"].(string), lightSetting["appid"].(string))
	if err != nil {
		log.Errorf("cant get the client of brand %s, id %s, error: %+v", lightSetting["brand"].(string), lightSetting["appid"].(string), err)
	}
	//use the appID & projID to find the setting
	//
	// control single deivce
	if header.Name == "Control" {
		payload := payloads[0]
		cmd, args := payload.Cmd, payload.Argument
		rsp, err := client.Execute(cmd, args)
		if err != nil {
			log.Errorf("execute api control error %+v", err)
			return
		}
		log.Infof("Payload:%+v Response: %+v", payload, rsp)
	}
}
