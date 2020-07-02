package rest

import "time"

//Client 调用接口的客户端
type Client interface {
	//Init() 客户端初始化函数，通常是验证本地的token，没有则刷新token
	Init() error
	//Execute()执行具体的接口参数
	Execute(api string, body interface{}) (*APIResponse, error)
	//GetAppKey() get the appkey for the client
	GetAppKey() string
}

//Options 相关的Api初始化选项
type Options struct {
	IOTID       string        //iot id
	AppID       string        //App ID
	AppSecret   string        //App Secret
	Username    string        //用户名
	TokenExpire time.Duration //Token有效期，通常是一个数字，用于存放在redis缓存的时间
	Enviroment  string        //服务的调用环境，生产/测试，Prod/Test
	BaseURL     string        //服务的基础URL地址
}

//Data the body of the rest server
type Data map[string]interface{}

//APIResponse api返回值
type APIResponse struct {
	HTTPStatus int  `json:"httpStatus"`
	Data       Data `json:"data"`
}
