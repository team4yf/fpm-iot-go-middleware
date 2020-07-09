package service

import (
	"github.com/team4yf/fpm-iot-go-middleware/config"
	"github.com/team4yf/fpm-iot-go-middleware/internal/model"
	"github.com/team4yf/fpm-iot-go-middleware/internal/repository"
	repo "github.com/team4yf/fpm-iot-go-middleware/internal/repository"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/cache"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/cache/rds"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/pool"
)

type cacher interface {
	GetCache() cache.Cache
}

//ClientService for client data manager
type ClientService interface {
	cacher
	Get(appid, enviroment string) (*model.Client, error)
	ListByCondition(expression string, conditions ...string) ([]*model.Client, error)
}

type simpleClientService struct {
	c    cache.Cache
	repo repository.ClientRepo
}

//NewSimpleClientService Create a new simpleClientService
func NewSimpleClientService() ClientService {
	rdsClient := pool.GetRedis()
	service := &simpleClientService{
		c:    rds.NewRedisCache(config.AppName, rdsClient),
		repo: repo.NewClientRepo(),
	}
	return service
}

func (s *simpleClientService) GetCache() cache.Cache {
	return s.c
}

func (s *simpleClientService) Get(appid, enviroment string) (client *model.Client, err error) {

	client, err = s.repo.Get(appid, enviroment)

	return
}

func (s *simpleClientService) ListByCondition(expression string, conditions ...string) (clients []*model.Client, err error) {
	clients, err = s.repo.ListByCondition(expression, conditions...)
	return
}
