package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	m "github.com/team4yf/fpm-iot-go-middleware/internal/model"
)

type ClientRepo interface {
	Create(*m.Client) error
	Get(appid string, enviroment string) (*m.Client, error)
}
type clientRepo struct {
	db *gorm.DB
}

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
