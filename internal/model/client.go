package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Client 第三方应用客户端对应的实体类
type Client struct {
	gorm.Model  `json:"-"`
	AppKey      string        `json:"appKey" gorm:"index"` // 应用ID
	SecretKey   string        `json:"secretKey"`
	Expired     time.Duration `json:"expired"`    // 对应的token有效期，单位秒
	Username    string        `json:"username"`   //对应的用户名
	Brand       string        `json:"brand"`      //对应的厂家的品牌
	Name        string        `json:"name"`       // 对应的client的名称
	APIBaseURL  string        `json:"apiBaseURL"` // 接口对应的url地址
	Environment string        `json:"enviroment"` // 接口对应的环境： TEST/SANDBOX/PRODUCTION
	EnableSSL   bool          `json:"enableSSL"`  // 是否启用SSL通信
	CertPath    string        `json:"certPath"`   // 证书路径
	Type        string        `json:"type"`       // client 对应的类型：Env, Light, Camer
	Status      int           `json:"status"`     // client状态
}

//TableName 对应表名
func (Client) TableName() string {
	return "fim_client"
}
