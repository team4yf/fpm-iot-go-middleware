//Package lt10 瓴泰科技智慧路灯
package lt10

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/team4yf/fpm-iot-go-middleware/external/rest"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/utils"
	"github.com/team4yf/yf-fpm-server-go/pkg/cache"
)

//DeviceType 设备类型，1：回路  2：控制器  默认为1
type DeviceType int

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

//defaultClient 默认的客户端接口实现
type defaultClient struct {
	options    *rest.Options //初始化参数
	cacher     cache.Cache   //缓存实例
	token      string        //获取到的token
	expireTime time.Time     //过期时间
	lock       sync.RWMutex  //同步锁
	inited     bool          //初始化标示位

}

func refreshToken(client *defaultClient, force bool) (token string, err error) {

	opts := client.options
	key := fmt.Sprintf("token:light:lt10:%s", opts.AppID)
	exists := false
	if exists, err = client.cacher.IsSet(key); err != nil {
		return
	}
	//如果force为true，则需要重置，将exists设置为false
	exists = exists && !force
	//存在则直接返回
	if exists {
		token, err = client.cacher.GetString(key)
		return
	}
	rspWrapper := utils.GetWithHeader(opts.BaseURL+apiURL["getAccessToken"], map[string]string{
		"appId": opts.AppID,
	}, 120)
	if !rspWrapper.Success {
		err = rspWrapper.Err
		return
	}

	var data rest.Data
	err = rspWrapper.ConvertBody(&data)
	if err != nil {
		return
	}
	// log.Infof("getToken: %+v", apiRsp.Data)
	token = data["data"].(string)
	err = client.cacher.SetString(key, token, client.options.TokenExpire*time.Millisecond)
	return
}

//NewClient 新建一个终端
func NewClient(opts *rest.Options, cacher cache.Cache) rest.Client {
	client := &defaultClient{
		options: opts,
		inited:  false,
		cacher:  cacher,
	}
	return client
}

func (cli *defaultClient) Init() (err error) {
	cli.lock.Lock()
	defer cli.lock.Unlock()
	token, err := refreshToken(cli, cli.inited)
	if err != nil {
		return
	}
	cli.inited = true
	cli.token = token
	cli.expireTime = time.Unix(time.Now().Unix()+(int64)(time.Millisecond*cli.options.TokenExpire), 0)
	return
}

func (cli *defaultClient) Execute(api string, body interface{}) (rsp *rest.APIResponse, err error) {
	//check the key expired time, refresh token
	if cli.expireTime.Before(time.Now()) {
		if err = cli.Init(); err != nil {
			return
		}
	}

	data, _ := json.Marshal(body)
	// log.Infof("Execute: %s", (string)(data))
	opts := cli.options
	rspWrapper := utils.PostJSONWithHeader(opts.BaseURL+apiURL[api], map[string]string{
		"accessToken": cli.token,
	}, data, 120)
	if !rspWrapper.Success {
		err = rspWrapper.Err
		return
	}
	var rspData rest.Data
	err = rspWrapper.ConvertBody(&rspData)
	if err != nil {
		return
	}
	rsp = &rest.APIResponse{
		HTTPStatus: rspWrapper.StatusCode,
		Data:       rspData,
	}
	return
}
func (cli *defaultClient) GetAppKey() string {
	return cli.options.AppID
}
