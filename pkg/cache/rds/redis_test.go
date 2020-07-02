package rds

import (
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/team4yf/fpm-iot-go-middleware/config"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/pool"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/test"
)

func TestRedisIsSet(t *testing.T) {
	var success bool
	test.InitTestConfig("../../../conf/config.test.json")
	pool.InitRedis(config.RedisConfig)

	client, success := pool.Get("redis")
	assert.Equal(t, true, success, "Get should return true")
	cli := client.(*redis.Client)
	redisCache := NewRedisCache(cli)

	var err error
	err = redisCache.SetInt("a", 1, 100*time.Second)
	assert.Nil(t, err, "setInt should not occur error")

	success, err = redisCache.IsSet("a")
	assert.Nil(t, err, "IsSet should not occur error")
	assert.Equal(t, true, success, "IsSet should return true")

	success, err = redisCache.Remove("a")
	assert.Nil(t, err, "Remove should not occur error")
	assert.Equal(t, true, success, "Remove should return true")

	success, err = redisCache.IsSet("a")
	assert.Nil(t, err, "IsSet should not occur error")
	assert.Equal(t, false, success, "IsSet should return false")

}
