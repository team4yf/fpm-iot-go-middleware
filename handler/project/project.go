//Package project for manager the project
package project

import (
	"errors"

	"github.com/team4yf/fpm-iot-go-middleware/internal/model"
	"github.com/team4yf/yf-fpm-server-go/fpm"
	"github.com/team4yf/yf-fpm-server-go/pkg/db"
)

//InitBiz 初始化相关的业务函数
func InitBiz(fpmApp *fpm.Fpm) {
	bizModule := make(fpm.BizModule, 0)
	bizModule["create"] = func(param *fpm.BizParam) (data interface{}, err error) {
		var req model.Project
		err = param.Convert(&req)
		if err != nil {
			return
		}
		req.Status = 1
		dbclient, _ := fpmApp.GetDatabase("pg")
		var count int64
		q := db.NewQuery()
		q.SetTable(req.TableName()).SetCondition("project_id = ? and app_id = ?", req.ProjectID, req.AppID)
		err = dbclient.Count(q.BaseData, &count)
		if err != nil {
			return
		}
		if count > 0 {
			err = errors.New(`Project exists`)
			return
		}

		err = dbclient.Create(q.BaseData, &req)
		if err != nil {
			return
		}

		return req, nil
	}
	bizModule["update"] = func(param *fpm.BizParam) (data interface{}, err error) {
		var req model.Project
		err = param.Convert(&req)
		if err != nil {
			return
		}
		req.Status = 1
		dbclient, _ := fpmApp.GetDatabase("pg")
		var count int64
		q := db.NewQuery()
		q.SetTable(req.TableName()).SetCondition("project_id = ? and app_id = ?", req.ProjectID, req.AppID)
		err = dbclient.Count(q.BaseData, &count)
		if err != nil {
			return
		}
		if count < 1 {
			err = errors.New(`Project not exists`)
			return
		}

		err = dbclient.Updates(q.BaseData, db.CommonMap{
			"setting": req.Setting,
		}, &count)
		if err != nil {
			return
		}

		return req, nil
	}
	fpmApp.AddBizModule("project", &bizModule)
}
