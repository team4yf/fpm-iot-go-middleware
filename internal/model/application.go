package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Application 项目对应的实体类
//通常对应着一个部署的应用平台
type Application struct {
	gorm.Model  `json:"-"`
	AppID       string    `json:"appId" gorm:"index"` // 应用ID
	Name        string    `json:"name"`               // 应用名称
	Status      int       `json:"status"`             // 应用状态
	ActiveAt    time.Time `json:"activeAt"`           // 激活时间
	Company     string    `json:"company"`            // 应用管理的公司
	Contact     string    `json:"contact"`            // 联系人
	Email       string    `json:"email"`              // Email
	Mobile      string    `json:"mobile"`             // Mobile
	ActiveCode  string    `json:"activeCode"`         // 激活码
	HomePage    string    `json:"homePage"`           // 主页地址
	Health      string    `json:"health"`             // 健康检查的地址
	Description string    `json:"description"`        // 应用详情
	DeployMode  string    `json:"deplyMode"`          // 应用服务部署的模式，Singleton/Multi
}

//TableName 对应表名
func (Application) TableName() string {
	return "fim_application"
}
