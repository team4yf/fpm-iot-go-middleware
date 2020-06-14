package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	m "github.com/team4yf/fpm-iot-go-middleware/internal/model"
)

type ApplictionRepo interface {
	Create(*m.Application) error
	Get(string) (*m.Application, error)
}
type applictionRepo struct {
	db *gorm.DB
}

func NewApplictionRepo() ApplictionRepo {
	return &applictionRepo{
		db: m.Conn,
	}
}

func (r *applictionRepo) Create(entity *m.Application) error {
	err := r.db.Create(&entity).Error
	if err != nil {
		return errors.Wrap(err, "[ApplictionRepo] Create error")
	}

	return nil
}

func (r *applictionRepo) Get(appId string) (*m.Application, error) {
	entity := &m.Application{}
	err := r.db.Where("app_id = ?", appId).First(entity).Error
	if err != nil {
		return nil, errors.Wrap(err, "[ApplictionRepo] Get error")
	}

	return entity, nil
}
