//Package lintaiv10 瓴泰科技智慧路灯
package lintaiv10

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/utils"
)

//Options 相关的Api初始化选项
type Options struct {
	AppID       string //App ID
	AppSecret   string //App Secret
	Username    string //用户名
	TokenExpire int64  //Token有效期，通常是一个数字，用于存放在redis缓存的时间
	Enviroment  string //服务的调用环境，生产/测试，Prod/Test
	BaseURL     string //服务的基础URL地址
}

//DeviceType 设备类型，1：回路  2：控制器  默认为1
type DeviceType int

//APIResponse api返回值
type APIResponse struct {
	Code uint8       `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

//LightControlType  2 对应控制器类型
const LightControlType = 2

var apiURL map[string]string

func init() {
	apiURL = map[string]string{
		"getAccessToken":       "/auth/getAccessToken",
		"command":              "/lightControl/command",
		"setLocalStrategy":     "/lightControl/setLocalStrategy",
		"batchUpdateFrequency": "/lightControl/batchUpdateFrequency",
		"setEarlyWarn":         "/lightControl/setEarlyWarn",
		"cancelEarlyWarn":      "/lightControl/cancelEarlyWarn",
		"getDevice":            "/lightControl/list",
	}
}

//Client 调用接口的客户端
type Client interface {
	//Init() 客户端初始化函数，通常是验证本地的token，没有则刷新token
	Init() error
	//Execute()执行具体的接口参数
	Execute(api string, body interface{}) (*APIResponse, error)
}

//defaultClient 默认的客户端接口实现
type defaultClient struct {
	options     *Options      //初始化参数
	redisClient *redis.Client //redis客户端链接实例
	token       string        //获取到的token
	expireTime  int64         //过期时间
	lock        sync.RWMutex  //同步锁
	inited      bool          //初始化标示位

}

func refreshToken(opts *Options) string {
	rspWrapper := utils.GetWithHeader(opts.BaseURL+apiURL["getAccessToken"], map[string]string{
		"appId": opts.AppID,
	}, 120)
	log.Printf("getToken: %+v", rspWrapper)
	return rspWrapper.Body
}

//NewClient 新建一个终端
func NewClient(opts *Options) Client {
	client := &defaultClient{
		options: opts,
		inited:  false,
	}
	return client
}

var errFoo = errors.New("stub")

func (cli *defaultClient) Init() (err error) {
	cli.lock.Lock()
	defer cli.lock.Unlock()
	token := refreshToken(cli.options)
	cli.inited = true
	cli.token = token
	cli.expireTime = time.Now().Unix() + cli.options.TokenExpire
	return
}

func (cli *defaultClient) Execute(api string, body interface{}) (rsp *APIResponse, err error) {
	err = errFoo
	return
}
