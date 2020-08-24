package service

import (
	"github.com/team4yf/fpm-iot-go-middleware/internal/model"
	"github.com/team4yf/fpm-iot-go-middleware/internal/repository"
	repo "github.com/team4yf/fpm-iot-go-middleware/internal/repository"
	"github.com/team4yf/yf-fpm-server-go/pkg/cache"
)

//ClientService for client data manager
type ClientService interface {
	Get(appid, enviroment string) (*model.Client, error)
	ListByCondition(expression string, conditions ...string) ([]*model.Client, error)
}

type simpleClientService struct {
	c    cache.Cache
	repo repository.ClientRepo
}

//NewSimpleClientService Create a new simpleClientService
func NewSimpleClientService(c cache.Cache) ClientService {
	service := &simpleClientService{
		c:    c,
		repo: repo.NewClientRepo(),
	}
	return service
}

func (s *simpleClientService) Get(appid, enviroment string) (client *model.Client, err error) {

	client, err = s.repo.Get(appid, enviroment)

	return
}

func (s *simpleClientService) ListByCondition(expression string, conditions ...string) (clients []*model.Client, err error) {
	clients, err = s.repo.ListByCondition(expression, conditions...)
	return
}
