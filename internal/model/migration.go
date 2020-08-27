package model

import (
	"github.com/team4yf/yf-fpm-server-go/pkg/db"
)

//Migration auto create or upgrade table
type Migration struct {
	DS db.Database
}

//Install create the tables
func (migration *Migration) Install() error {

	tables := []interface{}{
		&Device{},
		&Client{},
		&Application{},
		&Project{},
		&MQTTUser{},
	}
	migration.DS.AutoMigrate(tables...)
	mock()
	return nil
}

//TODO: mock some data
func mock() error {
	return nil
}
