package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/test"
)

func TestRedisIsSet(t *testing.T) {

	test.InitTestConfig("../../conf/config.test.json")

	redisCache := NewRedisCache()

	var err error
	err = redisCache.SetInt("a", 1, 100*time.Second)
	assert.Nil(t, err, "setInt should not occur error")

	var success bool
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
