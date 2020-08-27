package service

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/team4yf/fpm-iot-go-middleware/pkg/utils"
	"github.com/team4yf/yf-fpm-server-go/fpm"
	"github.com/team4yf/yf-fpm-server-go/pkg/cache"

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
	cache cache.Cache
}

//NewSimpleDeviceService create a new service
func NewSimpleDeviceService(c cache.Cache) DeviceService {
	mux.Lock()
	defer mux.Unlock()
	if instance != nil {
		return instance
	}
	service := &SimpleDeviceService{
		cache: c,
	}
	instance = service
	return service
}

//RegisterDevice insert the device into the db
func (s *SimpleDeviceService) RegisterDevice(device *m.Device) (err error) {
	//get the device, return if exists
	db, _ := fpm.Default().GetDatabase("pg")
	total := 0
	err = db.Model(device).Condition("sn = ? and status=1 ", device.SN).Count(&total).Error()
	if err != nil {
		return
	}
	if total > 0 {
		return
	}
	return db.Create(device).Error()
}

//GetDeviceInfo get the device info from the db
func (s *SimpleDeviceService) GetDeviceInfo(sn string) (device *m.Device, err error) {
	db, _ := fpm.Default().GetDatabase("pg")
	device = &m.Device{}
	err = db.Model(device).Condition("sn = ? and status = 1", sn).First(device).Error()
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
	db, _ := fpm.Default().GetDatabase("pg")
	proj := &m.Project{}

	err = db.Model(proj).Condition("app_id = ? and project_id = ?", appID, projectID).First(proj).Error()
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
