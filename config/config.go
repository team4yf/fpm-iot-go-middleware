package config

import (
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/log"
)

var (
	//Db the db config
	Db          DBSetting
	RedisConfig *redis.Options
)

// Config 读取配置
type Config struct {
	Name string
}

//DBSetting the config about the db
type DBSetting struct {
	Engine   string
	User     string
	Password string
	Host     string
	Port     int
	Database string
	Charset  string
	ShowSQL  bool
}

var (
	REDIS_HOST  string
	REDIS_PORT  string
	REDIS_DB    int
	REDIS_PASS  string
	MQTT_URL    string
	MQTT_USER   string
	MQTT_PASS   string
	MQTT_QOS    int
	MQTT_RETAIN bool
	PG_HOST     string
	PG_PORT     int
	PG_USER     string
	PG_DB       string
	PG_PASS     string
	PG_SHOWSQL  bool
)

func loadConfig() {
	viper.SetDefault("redis", map[string]interface{}{"host": "localhost", "port": 6379, "db": 0, "password": "admin123"})
	viper.SetDefault("mqtt", map[string]interface{}{"url": "mqtt.yunplus.io:1883", "user": "admin", "qos": 0, "retain": false, "pass": "123123123"})
	viper.SetDefault("postgresql", map[string]interface{}{"host": "localhost", "port": 5432, "db": "fim", "user": "postgres", "password": "Fim741235896", "showsql": true})
	viper.SetConfigType("json")
	viper.AddConfigPath(".") // 设置配置文件和可执行二进制文件在用一个目录
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Info("no such config file")
		} else {
			// Config file was found but another error was produced
			log.Info("read config error")
		}
		log.Fatal(err) // 读取配置文件失败致命错误
	}
	REDIS_HOST = viper.GetString("redis.host")
	REDIS_PORT = viper.GetString("redis.port")
	REDIS_DB = viper.GetInt("redis.db")
	REDIS_PASS = viper.GetString("redis.password")

	PG_HOST = viper.GetString("postgresql.host")
	PG_PORT = viper.GetInt("postgresql.port")
	PG_USER = viper.GetString("postgresql.user")
	PG_DB = viper.GetString("postgresql.db")
	PG_PASS = viper.GetString("postgresql.password")
	PG_SHOWSQL = viper.GetBool("postgresql.showsql")

	MQTT_URL = viper.GetString("mqtt.url")
	MQTT_USER = viper.GetString("mqtt.user")
	MQTT_PASS = viper.GetString("mqtt.pass")
	MQTT_QOS = viper.GetInt("mqtt.qos")
	MQTT_RETAIN = viper.GetBool("mqtt.retain")
	log.Debug(REDIS_HOST, REDIS_PORT, REDIS_DB, REDIS_PASS, MQTT_URL, MQTT_USER, MQTT_PASS, MQTT_QOS, MQTT_RETAIN)
}

func (c *Config) GetConfigOrDefault(key string, dft interface{}) interface{} {
	if viper.IsSet(key) {
		return viper.Get(key)
	}
	return dft

}
func (c *Config) IsSet(key string) bool {
	return viper.IsSet(key)
}
func (c *Config) GetMapOrDefault(key string, dft map[string]interface{}) map[string]interface{} {
	if viper.IsSet(key) {
		return viper.GetStringMap(key)
	}
	return dft
}

// Init 初始化配置，默认读取config.local.yaml
func Init(cfg string) error {
	c := Config{
		Name: cfg,
	}

	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}
	// 初始化日志包
	c.initLog()

	loadConfig()
	return nil
}

func (cfg *Config) initConfig() error {
	viper.AutomaticEnv()      // 读取匹配的环境变量
	viper.SetEnvPrefix("FPM") // 读取环境变量的前缀为 BS
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	deployMode := viper.GetString("deploy.mode")
	if deployMode == "" {
		deployMode = "local"
	}
	deployMode = strings.ToLower(deployMode)
	//read config file by BS_DEPLOY_MODE=PROD
	fmt.Println("DEPLOY_MODE:" + viper.GetString("deploy.mode"))
	if cfg.Name != "" {
		viper.SetConfigFile(cfg.Name) // 如果指定了配置文件，则解析指定的配置文件
	} else {
		viper.AddConfigPath("conf") // 如果没有指定配置文件，则解析默认的配置文件
		viper.SetConfigName("config." + deployMode)
	}
	viper.SetConfigType("json") // 设置配置文件格式为json

	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		return errors.WithStack(err)
	}

	Db = DBSetting{
		Engine:   viper.GetString("db.engine"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetInt("db.port"),
		Database: viper.GetString("db.database"),
		Charset:  viper.GetString("db.charset"),
		ShowSQL:  viper.GetBool("db.showSql"),
	}
	RedisConfig = &redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.passwd"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.pool"),
	}
	return nil
}

func (cfg *Config) initLog() {
	config := log.Config{
		Writers:         viper.GetString("log.writers"),
		LoggerLevel:     viper.GetString("log.level"),
		LoggerFile:      viper.GetString("log.logger_file"),
		LoggerWarnFile:  viper.GetString("log.logger_warn_file"),
		LoggerErrorFile: viper.GetString("log.logger_error_file"),
		LogFormatText:   viper.GetBool("log.log_format_text"),
		RollingPolicy:   viper.GetString("log.rollingPolicy"),
		LogRotateDate:   viper.GetInt("log.log_rotate_date"),
		LogRotateSize:   viper.GetInt("log.log_rotate_size"),
		LogBackupCount:  viper.GetInt("log.log_backup_count"),
	}
	err := log.NewLogger(&config, log.InstanceZapLogger)
	if err != nil {
		fmt.Printf("InitWithConfig err: %v", err)
	}
}
