package service

import (
	"github.com/pkg/errors"
	"github.com/team4yf/fpm-iot-go-middleware/internal/model"
	"github.com/team4yf/yf-fpm-server-go/fpm"
	"github.com/team4yf/yf-fpm-server-go/pkg/cache"
)

//ClientService for client data manager
type ClientService interface {
	Get(appid, enviroment string) (*model.Client, error)
	ListByCondition(expression string, conditions ...interface{}) ([]*model.Client, error)
}

type simpleClientService struct {
	app *fpm.Fpm
	c   cache.Cache
}

//NewSimpleClientService Create a new simpleClientService
func NewSimpleClientService(app *fpm.Fpm, c cache.Cache) ClientService {
	service := &simpleClientService{
		c:   c,
		app: app,
	}
	return service
}

func (s *simpleClientService) Get(appid, enviroment string) (client *model.Client, err error) {
	client = &model.Client{}
	db, ok := s.app.GetDatabase("pg")
	if !ok {
		return nil, errors.New("[clientService] Get database interface ")
	}
	err = db.Model(client).Condition("app_id = ? and enviroment = ?", appid, enviroment).First(client).Error()
	if err != nil {
		return nil, errors.Wrap(err, "[clientService] Get error")
	}
	return
}

func (s *simpleClientService) ListByCondition(expression string, conditions ...interface{}) (clients []*model.Client, err error) {
	clients = make([]*model.Client, 0)
	db, ok := s.app.GetDatabase("pg")
	if !ok {
		return nil, errors.New("[clientService] Get database interface ")
	}
	err = db.Model(model.Client{}).Condition(expression, conditions...).Find(&clients).Error()
	return
}
