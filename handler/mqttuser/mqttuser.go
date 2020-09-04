//Package mqttuser for manager the mqtt user
package mqttuser

import (
	"errors"

	"github.com/team4yf/fpm-iot-go-middleware/internal/model"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/utils"
	"github.com/team4yf/yf-fpm-server-go/fpm"
	"github.com/team4yf/yf-fpm-server-go/pkg/db"
)

//InitBiz 初始化相关的业务函数
func InitBiz(fpmApp *fpm.Fpm) {
	bizModule := make(fpm.BizModule, 0)
	bizModule["create"] = func(param *fpm.BizParam) (data interface{}, err error) {
		var req model.MQTTUser
		err = param.Convert(&req)
		if err != nil {
			return
		}
		req.Salt = utils.GenShortID()
		req.Status = 0
		req.Password = utils.Sha256Encode(req.Password + req.Salt)
		dbclient, _ := fpmApp.GetDatabase("pg")
		var count int64
		q := db.NewQuery()
		q.SetTable(req.TableName()).SetCondition("username = ? and app_id = ?", req.Username, req.AppID)
		err = dbclient.Count(q.BaseData, &count)
		if err != nil {
			return
		}
		if count > 0 {
			err = errors.New(`Username exists`)
			return
		}

		err = dbclient.Create(q.BaseData, &req)
		if err != nil {
			return
		}

		return req, nil
	}
	fpmApp.AddBizModule("mqttuser", &bizModule)
}
