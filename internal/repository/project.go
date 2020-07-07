package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	m "github.com/team4yf/fpm-iot-go-middleware/internal/model"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/utils"
)

type ProjectRepo interface {
	Create(m.Project) error
	GetSetting(appID string, projectID int64) (map[string]interface{}, error)
}

type projectRepo struct {
	db *gorm.DB
}

func NewProjectRepo() ProjectRepo {
	return &projectRepo{
		db: m.Db,
	}
}

func (r *projectRepo) Create(entity m.Project) error {
	err := r.db.Create(&entity).Error
	if err != nil {
		return errors.Wrap(err, "[ProjectRepo] Create error")
	}

	return nil
}

func (r *projectRepo) GetSetting(appID string, projectID int64) (setting map[string]interface{}, err error) {
	var proj m.Project
	err = r.db.Where(&m.Project{AppID: appID, ProjectID: projectID}).First(&proj).Error
	if err != nil {
		return
	}
	setting = make(map[string]interface{})
	err = utils.StringToStruct(proj.Setting, &setting)
	return
}
