package service

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/team4yf/fpm-iot-go-middleware/pkg/utils"
	"github.com/team4yf/yf-fpm-server-go/fpm"
	"github.com/team4yf/yf-fpm-server-go/pkg/cache"
	"github.com/team4yf/yf-fpm-server-go/pkg/db"

	m "github.com/team4yf/fpm-iot-go-middleware/internal/model"
)

const (
	cacheDuration = time.Duration(6) * time.Hour
)

var (
	DEVICE_NOT_EXISTS = errors.New("not exists device")
	instance          DeviceService
	mux               sync.Mutex
)

type cachedDevice struct {
	AppID     string `json:"appid"`
	ProjectID int64  `json:"projid"`
	DeviceID  string `json:"deviceid"`
}

//DeviceService the interface defination
type DeviceService interface {
	RegisterDevice(*m.Device) error
	GetDeviceInfo(string) (*m.Device, error)
	Receive(deviceType, brand, event, deviceID string) (string, int64, error)

	GetSetting(appID string, projectID int64) (map[string]interface{}, error)
}

//SimpleDeviceService the service object todo implement the interface
type SimpleDeviceService struct {
	cache    cache.Cache
	dbclient db.Database
}

//NewSimpleDeviceService create a new service
func NewSimpleDeviceService(c cache.Cache) DeviceService {
	mux.Lock()
	defer mux.Unlock()
	if instance != nil {
		return instance
	}
	dbclient, ok := fpm.Default().GetDatabase("pg")
	if !ok {
		panic("db connect error")
	}
	service := &SimpleDeviceService{
		cache:    c,
		dbclient: dbclient,
	}
	instance = service
	return service
}

//RegisterDevice insert the device into the db
func (s *SimpleDeviceService) RegisterDevice(device *m.Device) (err error) {
	var total int64
	q := db.NewQuery()
	q.SetTable(device.TableName())
	q.SetCondition("sn = ? and status=1 ", device.SN)
	err = s.dbclient.Count(q.BaseData, &total)
	if err != nil {
		return
	}
	if total > 0 {
		return
	}
	return s.dbclient.Create(q.BaseData, device)
}

//GetDeviceInfo get the device info from the db
func (s *SimpleDeviceService) GetDeviceInfo(sn string) (device *m.Device, err error) {
	device = &m.Device{}
	q := db.NewQuery()
	q.SetTable(device.TableName())
	q.SetCondition("sn = ? and status = 1", sn)
	err = s.dbclient.First(q, device)
	return
}

//GetSetting get setting from the table project
func (s *SimpleDeviceService) GetSetting(appID string, projectID int64) (setting map[string]interface{}, err error) {
	key := fmt.Sprintf("setting:%s:%d", appID, projectID)
	exists, err := s.cache.IsSet(key)
	if err != nil {
		return
	}
	if exists {
		setting = make(map[string]interface{})
		_, err = s.cache.GetObject(key, &setting)
		if err != nil {
			return
		}
		return
	}
	proj := &m.Project{}
	q := db.NewQuery()
	q.SetTable(proj.TableName())
	q.SetCondition("app_id = ? and project_id = ?", appID, projectID)
	err = s.dbclient.First(q, proj)
	if err != nil {
		return
	}
	setting = make(map[string]interface{})
	err = utils.StringToStruct(proj.Setting, &setting)
	err = s.cache.SetObject(key, setting, cacheDuration)
	if err != nil {
		return
	}
	return
}

//Receive 处理获取到数据之后的逻辑
// 主要逻辑就是通过和redis中的数据进行对比，读取保存的信息，返回出设备对应的项目ID
func (s *SimpleDeviceService) Receive(deviceType, brand, event, deviceID string) (appid string, projid int64, err error) {
	key := fmt.Sprintf("device:%s:%s:%s", deviceType, brand, deviceID)

	exists, err := s.cache.IsSet(key)
	if err != nil {
		return
	}

	data := &cachedDevice{}
	if exists {
		_, err = s.cache.GetObject(key, data)
		if err != nil {
			return
		}
		//从缓存中转换出数据
		appid = data.AppID
		projid = data.ProjectID
		return
	}

	// 从DB中获取
	device, err := s.GetDeviceInfo(deviceID)
	if err != nil {
		return
	}
	if device == nil {
		err = DEVICE_NOT_EXISTS
		return
	}
	// 缓存起来
	appid = device.AppID
	projid = device.ProjectID
	data.AppID = device.AppID
	data.ProjectID = device.ProjectID
	err = s.cache.SetObject(key, data, cacheDuration)
	return
}
