package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	//StatusDeviceInit init status, basiclly when create the device
	StatusDeviceInit = iota
	//StatusDeviceActived active the device, can communicate with it after that
	StatusDeviceActived
	//StatusDeviceOffline offline
	StatusDeviceOffline
	//StatusDeviceDenied device can't be used
	StatusDeviceDenied
)

//DeviceStatus the status of the device, it's defined in the model.StatusDevicexxx
type DeviceStatus int

//DeviceType the type of the device
type DeviceType string

//Device 设备对应的实体类
type Device struct {
	gorm.Model   `json:"-"`
	SN           string       `json:"sn" gorm:"index"`
	RegisterAt   time.Time    `json:"registerAt"`   // 注册时间
	LastUpdateAt time.Time    `json:"lastUpdateAt"` //上一次交互的时间
	Status       DeviceStatus `json:"status"`       // 设备状态
	ProjectID    int64        `json:"projectId"`    //对应的项目ID
	AppID        string       `json:"appId"`        //对应的应用ID
	Type         DeviceType   `json:"type"`         // 设备对应的类型 camera, env, light
	Brand        string       `json:"brand"`        // 设备对应的供应商品牌
}

//TableName 对应表名
func (Device) TableName() string {
	return "fim_device"
}
