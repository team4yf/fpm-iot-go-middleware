package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	m "github.com/team4yf/fpm-iot-go-middleware/internal/model"
)

type DeviceRepo interface {
	Create(*m.Device) error
	Get(string) (*m.Device, error)
	Check(string) (bool, error)
}
type deviceRepo struct {
	db *gorm.DB
}

func NewDeviceRepo() DeviceRepo {
	return &deviceRepo{
		db: m.Db,
	}
}

func (r *deviceRepo) Create(entity *m.Device) error {
	err := r.db.Create(&entity).Error
	if err != nil {
		return errors.Wrap(err, "[DeviceRepo] Create error")
	}

	return nil
}

func (r *deviceRepo) Get(sn string) (*m.Device, error) {
	entity := &m.Device{}
	err := r.db.Where("sn = ? and status=1 ", sn).First(entity).Error
	if err != nil {
		return nil, errors.Wrap(err, "[DeviceRepo] Get error")
	}

	return entity, nil
}

func (r *deviceRepo) Check(sn string) (exists bool, err error) {
	entity := &m.Device{}
	count := 0
	err = r.db.Model(entity).Where("sn = ? and status=1 ", sn).Count(&count).Error
	if err != nil {
		return
	}
	exists = count > 0
	return
}
