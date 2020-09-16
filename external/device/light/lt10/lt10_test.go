package lt10

import (
	"testing"

	"github.com/stretchr/testify/assert"
	_ "github.com/team4yf/fpm-go-plugin-cache-redis/plugin"
	_ "github.com/team4yf/fpm-go-plugin-mqtt-client/plugin"
	_ "github.com/team4yf/fpm-go-plugin-orm/plugins/pg"
	"github.com/team4yf/fpm-iot-go-middleware/external/rest"
	"github.com/team4yf/yf-fpm-server-go/fpm"
)

var client rest.Client

func Setup(t *testing.T) {

	fpmApp := fpm.NewWithConfig("./conf/config.test.json")
	fpmApp.Init()
	options := &rest.Options{
		AppID:       "LT0314fbf27a4d2986",
		AppSecret:   "1bc7b874c74623298a6",
		Username:    "18796664408",
		TokenExpire: 60 * 1000 * 24 * 7,

		Enviroment: "prod",
		BaseURL:    "http://101.132.142.5:8088/api",
	}
	cache, _ := fpmApp.GetCacher()
	client = NewClient(options, cache)

	err := client.Init()

	assert.Nil(t, err, "should not be err")
}

func TestApis(t *testing.T) {
	Setup(t)
	imei := "861050049029237"
	req := []map[string]interface{}{
		{
			"circuit":      1,
			"imei":         imei,
			"commandType":  "BRIGHTNESS",
			"commandValue": 30,
			"type":         LightControlType,
		},
	}
	rsp, err := client.Execute("command", req)
	assert.Nil(t, err, "Execute command should not be err")

	assert.Equal(t, 200, rsp.HTTPStatus, "Execute command should = 200")

}
