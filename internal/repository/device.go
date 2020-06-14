package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	m "github.com/team4yf/fpm-iot-go-middleware/internal/model"
)

type DeviceRepo interface {
	Create(*m.Device) error
	Get(string) (*m.Device, error)
}
type deviceRepo struct {
	db *gorm.DB
}

func NewDeviceRepo() DeviceRepo {
	return &deviceRepo{
		db: m.Conn,
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
	err := r.db.Where("sn = ?", sn).First(entity).Error
	if err != nil {
		return nil, errors.Wrap(err, "[DeviceRepo] Get error")
	}

	return entity, nil
}
