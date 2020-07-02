//Package rest the rest client
package rest

import "time"

//LinTaiOptions 相关的Api初始化选项
type LinTaiOptions struct {
	IOTID       string        //iot id
	AppID       string        //App ID
	AppSecret   string        //App Secret
	Username    string        //用户名
	TokenExpire time.Duration //Token有效期，通常是一个数字，用于存放在redis缓存的时间
	Enviroment  string        //服务的调用环境，生产/测试，Prod/Test
	BaseURL     string        //服务的基础URL地址
}

//LinTaiAPIResponse api返回值
type LinTaiAPIResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}
