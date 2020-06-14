package model

import (
	"github.com/jinzhu/gorm"
)

// 项目对应的实体类
type Project struct {
	gorm.Model `json:"-"`
	AppID      string `json:"appId" gorm:"index"` // 项目对应的应用ID
	Name       string `json:"name"`               // 项目名称
	Status     int    `json:"status"`             // 项目状态
	ProjectID  int64  `json:"projectId"`          //对应的项目ID
}

// 对应表名
func (Project) TableName() string {
	return "fim_project"
}
