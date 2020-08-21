//Package mqttuser for manager the mqtt user
package mqttuser

import (
	"github.com/team4yf/fpm-iot-go-middleware/internal/model"
	"github.com/team4yf/fpm-iot-go-middleware/internal/repository"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/utils"
	"github.com/team4yf/yf-fpm-server-go/ctx"
	"github.com/team4yf/yf-fpm-server-go/fpm"
)

var mqttUserRep repository.MQTTUserRepo

//CreateHandler create a new mqtt user
func CreateHandler() func(*ctx.Ctx, *fpm.Fpm) {
	mqttUserRep = repository.NewMQTTUserRepo()

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
		err = mqttUserRep.Create(&req)
		if err != nil {
			c.Fail(err)
			return
		}

		c.JSON(req)
	}

}
