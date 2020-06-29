package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	repo "github.com/team4yf/fpm-iot-go-middleware/internal/repository"

	"github.com/go-redis/redis/v8"
	m "github.com/team4yf/fpm-iot-go-middleware/internal/model"
)

var DEVICE_NOT_EXISTS = errors.New("not exists device")

var Ctx = context.Background()

type DeviceService interface {
	RegisterDevice(*m.Device) error
	GetDeviceInfo(string) (*m.Device, error)
	Receive(deviceType, brand, event, deviceID string) (string, int64, error)
}

type SimpleDeviceService struct {
	cli             *redis.Client
	applicationRepo repo.ApplictionRepo
	deviceRepo      repo.DeviceRepo
	projectRepo     repo.ProjectRepo
	clientRepo      repo.ClientRepo
}

func NewSimpleDeviceService(addr, passwd string, db int) DeviceService {
	opt := &redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	}
	service := &SimpleDeviceService{
		cli:             redis.NewClient(opt),
		applicationRepo: repo.NewApplictionRepo(),
		deviceRepo:      repo.NewDeviceRepo(),
		projectRepo:     repo.NewProjectRepo(),
		clientRepo:      repo.NewClientRepo(),
	}
	_, err := service.cli.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("redis cant connect ", err)
	}
	return service
}
func (s *SimpleDeviceService) RegisterDevice(device *m.Device) (err error) {
	return s.deviceRepo.Create(device)
}
func (s *SimpleDeviceService) GetDeviceInfo(sn string) (*m.Device, error) {
	return s.deviceRepo.Get(sn)
}

// 处理获取到数据之后的逻辑
// 主要逻辑就是通过和redis中的数据进行对比，读取保存的信息，返回出设备对应的项目ID
func (s *SimpleDeviceService) Receive(deviceType, brand, event, deviceID string) (string, int64, error) {
	key := fmt.Sprintf("device:%s:%s:%s", deviceType, brand, deviceID)
	data, err := s.cli.HGetAll(Ctx, key).Result()

	if err != nil {
		return "", -1, err
	}
	if len(data) == 0 {
		// 从DB中获取
		device, err := s.deviceRepo.Get(deviceID)
		if err != nil {
			return "", -1, err
		}
		if device == nil {
			return "", -1, DEVICE_NOT_EXISTS
		}
		// 缓存到redis中
		s.cli.HSet(Ctx, key, "appid", device.AppID)
		s.cli.HSet(Ctx, key, "projid", device.ProjectID)
		return device.AppID, device.ProjectID, nil
	}

	projid, err := strconv.ParseInt(data["projid"], 10, 64)
	if err != nil {
		return "", -1, err
	}
	return data["appid"], projid, nil
}
