package rds

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/cache"
)

type redisLocker struct {
	cli *redis.Client
}

//NewRedisLocker create a redisLocker
func NewRedisLocker(c *redis.Client) cache.SyncLocker {

	return &redisLocker{
		cli: c,
	}
}

func (r *redisLocker) GetLock(key string, expired time.Duration) bool {
	ok, err := r.cli.SetNX(timeoutCtx, key, 1, time.Duration(0)).Result()
	if err != nil {
		return false
	}
	// 设置过期时间，防止进程挂了之后，无法释放
	ok, err = r.cli.Expire(timeoutCtx, key, expired*time.Second).Result()
	if err != nil {
		r.cli.Del(timeoutCtx, key).Err()
		return false
	}
	return ok
}
func (r *redisLocker) ReleaseLock(key string) error {
	err := r.cli.Del(timeoutCtx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
