package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	config "github.com/team4yf/fpm-iot-go-middleware/config"
)

var (
	Conn = connect()
)

func init() {
	Conn.AutoMigrate(
		&Device{},
		&Client{},
		&Application{},
		&Project{},
	)
}

func connect() *gorm.DB {

	dbEndPoint := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable connect_timeout=2",
		config.PG_HOST,
		config.PG_PORT,
		config.PG_USER,
		config.PG_PASS,
		config.PG_DB)
	c, err := gorm.Open("postgres", dbEndPoint)
	if err != nil {
		panic(err)
	}

	c.DB().SetConnMaxLifetime(time.Minute * 5)
	c.DB().SetMaxIdleConns(20)
	c.DB().SetMaxOpenConns(500)

	c.LogMode(config.PG_SHOWSQL)
	return c
}
