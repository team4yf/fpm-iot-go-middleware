package rds

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	errs "github.com/pkg/errors"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/cache"
)

var (
	errNotDoneYet = errors.New("Not Done Yet!")
	errNoData     = errors.New("Find Nothing!")
	timeoutCtx    = context.Background()
	Cache         cache.Cache
)

type redisCache struct {
	cli       *redis.Client
	perfixKey string
}

//NewRedisCache 创建一个新的基于Redis实现的服务
// 需要传入配置的信息
func NewRedisCache(prefix string, c *redis.Client) cache.Cache {
	cache := &redisCache{
		cli:       c,
		perfixKey: prefix,
	}
	Cache = cache
	return cache
}

func (r *redisCache) PaddingKey(key string) string {
	return fmt.Sprintf("%s:%s", r.perfixKey, key)
}

func (r *redisCache) SetString(key, val string, duration time.Duration) error {
	if err := r.cli.Set(timeoutCtx, r.PaddingKey(key), val, duration).Err(); err != nil {
		return errs.Wrap(err, "set data to redis set err")
	}
	return nil
}
func (r *redisCache) SetObject(key string, val interface{}, duration time.Duration) error {
	raw, err := json.Marshal(val)
	if err != nil {
		return errs.Wrap(err, "marshal data err")
	}
	if err = r.cli.Set(timeoutCtx, r.PaddingKey(key), string(raw), duration).Err(); err != nil {
		return errs.Wrap(err, "set data to redis set err")
	}
	return nil
}

func (r *redisCache) Set(key string, val interface{}, duration time.Duration) error {
	return errNotDoneYet
}

func (r *redisCache) SetInt(key string, val int64, duration time.Duration) error {
	if err := r.cli.Set(timeoutCtx, r.PaddingKey(key), val, duration).Err(); err != nil {
		return errs.Wrap(err, "set data to redis set err")
	}
	return nil
}

func (r *redisCache) GetString(key string) (val string, err error) {
	val, err = r.cli.Get(timeoutCtx, r.PaddingKey(key)).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
			return
		}
		err = errs.Wrap(err, "redis do get error:"+r.PaddingKey(key))
		return
	}
	return

}

func (r *redisCache) Get(key string) (interface{}, error) {
	return r.GetString(key)
}
func (r *redisCache) GetInt(key string) (int64, error) {
	val, err := r.GetString(key)
	if err != nil {
		return -1, err
	}
	var i int64
	i, err = strconv.ParseInt(val, 10, 64)
	if err != nil {
		return -1, errs.Wrap(err, "parseInt error:"+key)
	}
	return i, nil
}

func (r *redisCache) GetObject(key string, obj interface{}) (bool, error) {
	val, err := r.GetString(key)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal([]byte(val), obj)
	if err != nil {
		return false, errs.Wrap(err, "Unmarshal error")
	}
	return true, nil

}

func (r *redisCache) GetMap(key string) (map[string]interface{}, error) {
	val, err := r.GetString(key)
	if err != nil {
		return nil, err
	}
	obj := make(map[string]interface{})
	err = json.Unmarshal([]byte(val), obj)
	if err != nil {
		return nil, errs.Wrap(err, "Unmarshal error")
	}
	return obj, nil
}
func (r *redisCache) IsSet(key string) (bool, error) {
	cmd := r.cli.Exists(timeoutCtx, r.PaddingKey(key))
	value, err := cmd.Result()
	if err != nil {
		return false, err
	}
	return value == 1, nil
}

func (r *redisCache) Remove(key string) (bool, error) {
	err := r.cli.Del(timeoutCtx, r.PaddingKey(key)).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *redisCache) SafetyIncr(key string, step int64) (bool, error) {
	err := r.cli.IncrBy(timeoutCtx, r.PaddingKey(key), step).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}
