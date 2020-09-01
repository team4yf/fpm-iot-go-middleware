package service

import (
	"github.com/pkg/errors"
	"github.com/team4yf/fpm-iot-go-middleware/internal/model"
	"github.com/team4yf/yf-fpm-server-go/fpm"
	"github.com/team4yf/yf-fpm-server-go/pkg/cache"
	"github.com/team4yf/yf-fpm-server-go/pkg/db"
)

//ClientService for client data manager
type ClientService interface {
	Get(appid, enviroment string) (*model.Client, error)
	ListByCondition(expression string, conditions ...interface{}) ([]*model.Client, error)
}

type simpleClientService struct {
	app      *fpm.Fpm
	c        cache.Cache
	dbclient db.Database
}

//NewSimpleClientService Create a new simpleClientService
func NewSimpleClientService(app *fpm.Fpm, c cache.Cache) ClientService {
	dbclient, ok := app.GetDatabase("pg")
	if !ok {
		panic("db connect error")
	}
	service := &simpleClientService{
		c:        c,
		app:      app,
		dbclient: dbclient,
	}
	return service
}

func (s *simpleClientService) Get(appid, enviroment string) (client *model.Client, err error) {
	client = &model.Client{}

	q := db.NewQuery()
	q.SetTable(client.TableName())
	q.SetCondition("app_id = ? and enviroment = ?", appid, enviroment)

	err = s.dbclient.First(q, client)
	if err != nil {
		return nil, errors.Wrap(err, "[clientService] Get error")
	}
	return
}

func (s *simpleClientService) ListByCondition(expression string, conditions ...interface{}) (clients []*model.Client, err error) {
	clients = make([]*model.Client, 0)
	q := db.NewQuery()
	q.SetTable((model.Client{}).TableName())
	q.SetCondition(expression, conditions...)
	err = s.dbclient.First(q, &clients)
	return
}
