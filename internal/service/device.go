package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/team4yf/fpm-iot-go-middleware/config"
	repo "github.com/team4yf/fpm-iot-go-middleware/internal/repository"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/cache"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/cache/rds"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/pool"

	m "github.com/team4yf/fpm-iot-go-middleware/internal/model"
)

const (
	cacheDuration = time.Duration(6) * time.Hour
)

var (
	DEVICE_NOT_EXISTS = errors.New("not exists device")
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
	cache           cache.Cache
	applicationRepo repo.ApplictionRepo
	deviceRepo      repo.DeviceRepo
	projectRepo     repo.ProjectRepo
	clientRepo      repo.ClientRepo
}

//NewSimpleDeviceService create a new service
func NewSimpleDeviceService() DeviceService {
	rdsClient, _ := pool.Get("redis")
	service := &SimpleDeviceService{
		cache:           rds.NewRedisCache(config.AppName, rdsClient.(*redis.Client)),
		applicationRepo: repo.NewApplictionRepo(),
		deviceRepo:      repo.NewDeviceRepo(),
		projectRepo:     repo.NewProjectRepo(),
		clientRepo:      repo.NewClientRepo(),
	}
	return service
}

//RegisterDevice insert the device into the db
func (s *SimpleDeviceService) RegisterDevice(device *m.Device) (err error) {
	return s.deviceRepo.Create(device)
}

//GetDeviceInfo get the device info from the db
func (s *SimpleDeviceService) GetDeviceInfo(sn string) (*m.Device, error) {
	return s.deviceRepo.Get(sn)
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
	setting, err = s.projectRepo.GetSetting(appID, projectID)
	if err != nil {
		return
	}
	err = s.cache.SetObject(key, setting, cacheDuration)
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
	device, err := s.deviceRepo.Get(deviceID)
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
