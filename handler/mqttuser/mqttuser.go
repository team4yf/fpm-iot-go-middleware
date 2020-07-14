//Package mqttuser for manager the mqtt user
package mqttuser

import (
	"net/http"

	"github.com/team4yf/fpm-iot-go-middleware/internal/core"
	"github.com/team4yf/fpm-iot-go-middleware/internal/model"
	"github.com/team4yf/fpm-iot-go-middleware/internal/repository"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/utils"
)

var mqttUserRep repository.MQTTUserRepo

//CreateHandler create a new mqtt user
func CreateHandler(app *core.App) func(http.ResponseWriter, *http.Request) {
	mqttUserRep = repository.NewMQTTUserRepo()

	return func(w http.ResponseWriter, r *http.Request) {
		var req model.MQTTUser
		err := app.ParseBody(r, &req)
		if err != nil {
			app.FailWithError(w, err)
			return
		}
		req.Salt = utils.GenShortID()
		req.Status = 0
		req.Password = utils.Sha256Encode(req.Password + req.Salt)
		err = mqttUserRep.Create(&req)
		if err != nil {
			app.FailWithError(w, err)
			return
		}

		app.SendOk(w, req)
	}

}
