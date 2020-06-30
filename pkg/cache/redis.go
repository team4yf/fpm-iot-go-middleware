package cache

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	errs "github.com/pkg/errors"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/pool"
)

var NOT_DONE_YET = errors.New("not done yet!")

var NOT_DATA = errors.New("find nothing!")

type RedisCache struct {
	cli *redis.Client
}

// 创建一个新的基于Redis实现的服务
// 需要传入配置的信息
func NewRedisCache() Cache {
	cache := &RedisCache{
		cli: pool.Get(),
	}

	return cache
}

func (r *RedisCache) SetString(key, val string, duration time.Duration) error {
	if err := r.cli.Set(TIMEOUT_CTX, key, val, duration).Err(); err != nil {
		return errs.Wrap(err, "set data to redis set err")
	}
	return nil
}
func (r *RedisCache) SetObject(key string, val interface{}, duration time.Duration) error {
	raw, err := json.Marshal(val)
	if err != nil {
		return errs.Wrap(err, "marshal data err")
	}
	if err = r.cli.Set(TIMEOUT_CTX, key, string(raw), duration).Err(); err != nil {
		return errs.Wrap(err, "set data to redis set err")
	}
	return nil
}

func (r *RedisCache) Set(key string, val interface{}, duration time.Duration) error {
	return NOT_DONE_YET
}

func (r *RedisCache) SetInt(key string, val int64, duration time.Duration) error {
	if err := r.cli.Set(TIMEOUT_CTX, key, val, duration).Err(); err != nil {
		return errs.Wrap(err, "set data to redis set err")
	}
	return nil
}

func (r *RedisCache) Get(key string) (interface{}, error) {
	return r.GetString(key)
}
func (r *RedisCache) GetInt(key string) (int64, error) {
	val, err := r.cli.Get(TIMEOUT_CTX, key).Result()
	if err != nil {
		if err == redis.Nil {
			return -1, nil
		}
		return -1, errs.Wrap(err, "redis do get error:"+key)
	}
	var i int64
	i, err = strconv.ParseInt(val, 10, 64)
	if err != nil {
		return -1, errs.Wrap(err, "parseInt error:"+key)
	}
	return i, nil
}

func (r *RedisCache) GetObject(key string, obj interface{}) (bool, error) {
	val, err := r.cli.Get(TIMEOUT_CTX, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, errs.Wrap(err, "redis do get error:"+key)
	}

	err = json.Unmarshal([]byte(val), obj)
	if err != nil {
		return false, errs.Wrap(err, "Unmarshal error")
	}
	return true, nil

}

func (r *RedisCache) GetString(key string) (string, error) {
	if val, err := r.cli.Get(TIMEOUT_CTX, key).Result(); err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", errs.Wrap(err, "redis do get error:"+key)
	} else {
		return val, nil
	}
}
func (r *RedisCache) GetMap(key string) (map[string]interface{}, error) {
	return nil, NOT_DONE_YET
}
func (r *RedisCache) IsSet(key string) (bool, error) {
	cmd := r.cli.Exists(TIMEOUT_CTX, key)
	value, err := cmd.Result()
	if err != nil {
		return false, err
	}
	return value == 1, nil
}

func (r *RedisCache) Remove(key string) (bool, error) {
	err := r.cli.Del(TIMEOUT_CTX, key).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}
func (r *RedisCache) SafetyIncr(key string, step int64) (bool, error) {
	err := r.cli.IncrBy(TIMEOUT_CTX, key, step).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}
