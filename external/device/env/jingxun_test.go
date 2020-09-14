package env

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/team4yf/fpm-go-pkg/utils"
	_ "github.com/team4yf/fpm-go-plugin-cache-redis/plugin"
	_ "github.com/team4yf/fpm-go-plugin-mqtt-client/plugin"
	_ "github.com/team4yf/fpm-go-plugin-orm/plugins/pg"
	"github.com/team4yf/yf-fpm-server-go/fpm"
)

func TestParse(t *testing.T) {
	fpmApp := fpm.NewWithConfig("conf/config.test.json")
	fpmApp.Init()
	strData := "FEDC01171CB740E30F000001DE03001C0000010F000001CB000000340000004B000000870000007A0000001400"
	buf := stringToHex(strData)

	device := NewEnvDevice("jingxun")
	data, err := device.Parse(buf)
	assert.Nil(t, err, "should not be error")

	fpmApp.Execute("mqttclient.publish", &fpm.BizParam{
		"topic":   fmt.Sprintf(`$d2s/%s/partner/push`, data.Header.AppID),
		"payload": ([]byte)(utils.JSON2String(data)),
	})
	fmt.Printf("data: %v", data)
}
