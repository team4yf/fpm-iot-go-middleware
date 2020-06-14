package model

import (
	"github.com/jinzhu/gorm"
)

// 项目对应的实体类
type Application struct {
	gorm.Model `json:"-"`
	AppID      string `json:"appId" gorm:"index"` // 应用ID
	Name       string `json:"name"`               // 应用名称
	Status     int    `json:"status"`             // 应用状态
}

// 对应表名
func (Application) TableName() string {
	return "fim_application"
}
