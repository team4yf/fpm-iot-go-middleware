package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	m "github.com/team4yf/fpm-iot-go-middleware/internal/model"
)

type ProjectRepo interface {
	Create(m.Project) error
}

type projectRepo struct {
	db *gorm.DB
}

func NewProjectRepo() ProjectRepo {
	return &projectRepo{
		db: m.Conn,
	}
}

func (r *projectRepo) Create(entity m.Project) error {
	err := r.db.Create(&entity).Error
	if err != nil {
		return errors.Wrap(err, "[ProjectRepo] Create error")
	}

	return nil
}
