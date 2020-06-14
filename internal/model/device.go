package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// 设备对应的实体类
type Device struct {
	gorm.Model   `json:"-"`
	SN           string    `json:"sn" gorm:"index"`
	RegisterAt   time.Time `json:"registerAt"`   // 注册时间
	LastUpdateAt time.Time `json:"lastUpdateAt"` //上一次交互的时间
	Status       int       `json:"status"`       // 设备状态
	ProjectID    int64     `json:"projectId"`    //对应的项目ID
	AppID        string    `json:"appId"`        //对应的应用ID
}

// 对应表名
func (Device) TableName() string {
	return "fim_device"
}
