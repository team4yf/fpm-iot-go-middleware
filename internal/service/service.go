package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

var DEVICE_NOT_EXISTS = errors.New("not exists device")

type Service interface {
	Receive(deviceType, brand, event, deviceID string) (string, error)
}

type RedisService struct {
	cli *redis.Client
}

// 创建一个新的基于Redis实现的服务
// 需要传入配置的信息
func NewRedisService(addr, passwd string, db int) Service {
	opt := &redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	}
	service := &RedisService{
		cli: redis.NewClient(opt),
	}
	_, err := service.cli.Ping().Result()
	if err != nil {
		log.Fatal("redis cant connect ", err)
	}
	return service
}

// 处理获取到数据之后的逻辑
// 主要逻辑就是通过和redis中的数据进行对比，读取保存的信息，返回出设备对应的项目ID
func (s *RedisService) Receive(deviceType, brand, event, deviceID string) (string, error) {
	key := fmt.Sprintf("device:%s:%s", deviceType, brand)
	if val, err := s.cli.HGet(key, deviceID).Result(); err != nil {
		if err == redis.Nil {
			return "", DEVICE_NOT_EXISTS
		}
		return "", err
	} else {
		return val, nil
	}

}
