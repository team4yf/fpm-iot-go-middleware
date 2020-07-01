package config

import (
	"fmt"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	lintaiv10 "github.com/team4yf/fpm-iot-go-middleware/external/device/light/lintai/v10"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/log"
)

var (
	//Db the db config
	Db *DBSetting
	//RedisConfig the redis connect config
	RedisConfig *redis.Options
	//MqttConfig the mqtt connect config
	MqttConfig *MqttSetting
	//LintaiAppConfig lintai v10 config
	LintaiAppConfig *lintaiv10.Options
)

// Config 读取配置
type Config struct {
	Name string
}

// Init 初始化配置，默认读取config.local.json
func Init(configFilePath string) error {
	c := Config{
		Name: configFilePath,
	}

	// 初始化配置文件
	if err := c.loadConfig(); err != nil {
		return err
	}
	// 初始化日志包
	c.initLog()
	return nil
}

func (cfg *Config) loadConfig() error {
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

	Db = &DBSetting{
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
	MqttConfig = &MqttSetting{
		Options:  &MQTT.ClientOptions{},
		Retained: viper.GetBool("mqttserver.retain"),
		Qos:      (byte)(viper.GetInt("mqttserver.qos")),
	}

	MqttConfig.Options.AddBroker(fmt.Sprintf("tcp://%s:%d",
		viper.GetString("mqttserver.host"),
		viper.GetInt("mqttserver.port")))
	// opts.SetClientID("go-simple")
	MqttConfig.Options.SetUsername(viper.GetString("mqttserver.username"))
	MqttConfig.Options.SetPassword(viper.GetString("mqttserver.password"))

	LintaiAppConfig = &lintaiv10.Options{
		AppID:       viper.GetString("lintaiv10.appid"),
		AppSecret:   viper.GetString("lintaiv10.appsecret"),
		Username:    viper.GetString("lintaiv10.username"),
		TokenExpire: time.Duration(viper.GetInt("lintaiv10.expired")) * time.Second,

		Enviroment: viper.GetString("lintaiv10.env"),
		BaseURL:    viper.GetString("lintaiv10.baseurl"),
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

//GetConfigOrDefault get config from the config file or env, return default value if not set
func GetConfigOrDefault(key string, dft interface{}) interface{} {
	if viper.IsSet(key) {
		return viper.Get(key)
	}
	return dft
}

//IsSet judge if the  config setted if the config file or env.
func IsSet(key string) bool {
	return viper.IsSet(key)
}

//GetMapOrDefault get a map config from the config file or env, return default value if not set.
func GetMapOrDefault(key string, dft map[string]interface{}) map[string]interface{} {
	if viper.IsSet(key) {
		return viper.GetStringMap(key)
	}
	return dft
}
