package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/team4yf/fpm-iot-go-middleware/config"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	//Db the connection
	Db *gorm.DB
)

//CommonMap the common map for database
type CommonMap map[string]interface{}

type BaseDB struct {
	gorm.DB
}

//CreateDb create new instance
func CreateDb() *gorm.DB {
	//use the config for the app
	dsn := getDbEngineDSN(&config.Db)
	fmt.Printf("db setting: %+v\n", dsn)
	db, err := gorm.Open(config.Db.Engine, dsn)
	if err != nil {
		panic(err.Error())
	}

	db.DB().SetConnMaxLifetime(time.Minute * 5)
	db.DB().SetMaxIdleConns(20)
	db.DB().SetMaxOpenConns(500)

	db.LogMode(config.Db.ShowSQL)
	Db = db
	return db
}

//CreateTmpDb get templary db
func CreateTmpDb(db *config.DBSetting) (*gorm.DB, error) {
	dsn := getDbEngineDSN(db)

	return gorm.Open(db.Engine, dsn)
}

// 获取数据库引擎DSN  mysql,sqlite,postgres
func getDbEngineDSN(db *config.DBSetting) string {
	engine := strings.ToLower(db.Engine)
	dsn := ""
	switch engine {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&allowNativePasswords=true",
			db.User,
			db.Password,
			db.Host,
			db.Port,
			db.Database,
			db.Charset)
	case "postgres":
		dsn = fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
			db.User,
			db.Password,
			db.Host,
			db.Port,
			db.Database)
	}

	return dsn
}
