//Package mqttuser for manager the mqtt user
package mqttuser

import (
	"github.com/team4yf/fpm-iot-go-middleware/internal/model"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/utils"
	"github.com/team4yf/yf-fpm-server-go/ctx"
	"github.com/team4yf/yf-fpm-server-go/fpm"
)

//CreateHandler create a new mqtt user
func CreateHandler() func(*ctx.Ctx, *fpm.Fpm) {

	return func(c *ctx.Ctx, fpmApp *fpm.Fpm) {
		var req model.MQTTUser
		err := c.ParseBody(&req)
		if err != nil {
			c.Fail(err)
			return
		}
		req.Salt = utils.GenShortID()
		req.Status = 0
		req.Password = utils.Sha256Encode(req.Password + req.Salt)
		dbclient, _ := fpmApp.GetDatabase("pg")
		count := 0
		err = dbclient.Model(req).Condition("username = ? and app_id = ?", req.Username, req.AppID).Count(&count).Error()
		if err != nil {
			c.Fail(err)
			return
		}
		if count > 0 {
			c.Fail(map[string]string{
				`err`: `Username exists`,
			})
			return
		}
		err = dbclient.Model(req).Create(&req).Error()
		if err != nil {
			c.Fail(err)
			return
		}

		c.JSON(req)
	}

}
