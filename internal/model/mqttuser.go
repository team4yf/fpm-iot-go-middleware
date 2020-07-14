package model

import (
	"github.com/jinzhu/gorm"
)

//MQTTUser mqtt user struct
type MQTTUser struct {
	gorm.Model  `json:"-"`
	IsSuperuser bool   `json:"isSuperuser"`      //是否是超级用户
	Username    string `json:"username"`         //用户名
	Password    string `json:"password"`         //加密后的密码
	Salt        string `json:"salt"`             //加密用的盐信息
	Status      int    `json:"status,omitempty"` //mqtt 用户状态，0正常，非0异常
	AppID       string `json:"appid"`            //项目对应的AppID
}

//TableName 对应表名
func (MQTTUser) TableName() string {
	return "fim_mqtt_user"
}
