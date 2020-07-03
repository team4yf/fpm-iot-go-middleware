// 用于缓存的接口
package cache

import (
	"context"
	"time"
)

//TIMEOUT_CTX 默认的超时上下文
var TIMEOUT_CTX = context.Background()

type Cache interface {
	Set(key string, val interface{}, duration time.Duration) error
	SetString(key, val string, duration time.Duration) error
	SetInt(key string, val int64, duration time.Duration) error
	SetObject(key string, val interface{}, duration time.Duration) error

	Get(key string) (interface{}, error)
	IsSet(key string) (bool, error)
	Remove(key string) (bool, error)
	GetString(key string) (string, error)
	GetInt(key string) (int64, error)
	GetObject(key string, val interface{}) (bool, error)

	GetMap(key string) (map[string]interface{}, error)

	SafetyIncr(key string, step int64) (bool, error)
}