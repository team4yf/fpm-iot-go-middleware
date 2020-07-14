package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	m "github.com/team4yf/fpm-iot-go-middleware/internal/model"
)

//MQTTUserRepo the mqttuser repository
type MQTTUserRepo interface {
	Create(*m.MQTTUser) error
	ListByCondition(expression string, conditions ...string) ([]*m.MQTTUser, error)
}
type mqttUserRepo struct {
	db *gorm.DB
}

//NewMQTTUserRepo create the default repo
func NewMQTTUserRepo() MQTTUserRepo {
	return &mqttUserRepo{
		db: m.Db,
	}
}

func (r *mqttUserRepo) Create(entity *m.MQTTUser) error {
	err := r.db.Create(&entity).Error
	if err != nil {
		return errors.Wrap(err, "[mqttUserRepo] Create error")
	}

	return nil
}

func (r *mqttUserRepo) ListByCondition(expression string, conditions ...string) (users []*m.MQTTUser, err error) {
	users = make([]*m.MQTTUser, 0, 0)
	err = r.db.Where(expression, conditions).Find(&users).Error

	return
}
