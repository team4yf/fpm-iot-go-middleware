package cache

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/pool"
)

type SyncLocker interface {
	GetLock(key string, expire int64) bool
	ReleaseLock(key string) error
}

type RedisLocker struct {
	cli *redis.Client
}

func NewRedisLocker() SyncLocker {

	return &RedisLocker{
		cli: pool.Get(),
	}
}

func (r *RedisLocker) GetLock(key string, expire int64) bool {
	ok, err := r.cli.SetNX(TIMEOUT_CTX, key, 1, time.Duration(0)).Result()
	if err != nil {
		return false
	}
	// 设置过期时间，防止进程挂了之后，无法释放
	ok, err = r.cli.Expire(TIMEOUT_CTX, key, time.Duration(expire)*time.Second).Result()
	if err != nil {
		r.cli.Del(TIMEOUT_CTX, key).Err()
		return false
	}
	return ok
}
func (r *RedisLocker) ReleaseLock(key string) error {
	err := r.cli.Del(TIMEOUT_CTX, key).Err()
	if err != nil {
		return err
	}
	return nil
}
