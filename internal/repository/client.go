package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	m "github.com/team4yf/fpm-iot-go-middleware/internal/model"
)

//ClientRepo the client repository
type ClientRepo interface {
	Create(*m.Client) error
	Get(appid string, enviroment string) (*m.Client, error)
	ListByCondition(expression string, conditions ...string) ([]*m.Client, error)
}
type clientRepo struct {
	db *gorm.DB
}

//NewClientRepo create the default repo
func NewClientRepo() ClientRepo {
	return &clientRepo{
		db: m.Db,
	}
}

func (r *clientRepo) Create(entity *m.Client) error {
	err := r.db.Create(&entity).Error
	if err != nil {
		return errors.Wrap(err, "[ClientRepo] Create error")
	}

	return nil
}

func (r *clientRepo) Get(appid string, enviroment string) (*m.Client, error) {
	entity := &m.Client{}
	err := r.db.Where("app_id = ? and enviroment = ?", appid, enviroment).First(entity).Error
	if err != nil {
		return nil, errors.Wrap(err, "[ClientRepo] Get error")
	}

	return entity, nil
}

func (r *clientRepo) ListByCondition(expression string, conditions ...string) (clients []*m.Client, err error) {
	clients = make([]*m.Client, 0, 0)
	err = r.db.Where(expression, conditions).Find(&clients).Error

	return
}
