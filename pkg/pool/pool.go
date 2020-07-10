package pool

import (
	"context"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
)

var pools map[string]sync.Pool

var keyRedis = "redis"

//DefaultCtx the default context for redis session.
var DefaultCtx = context.Background()

func init() {
	pools = make(map[string]sync.Pool)
}

//InitRedis initionlize the resouce pool
func InitRedis(opt *redis.Options) {
	pools[keyRedis] = sync.Pool{

		New: func() interface{} {
			cli := redis.NewClient(opt)
			_, err := cli.Ping(DefaultCtx).Result()
			if err != nil {
				log.Fatal("redis cant connect ", err)
			}
			return cli
		},
	}
}

//Get get a resource from the pool
func Get(key string) (interface{}, bool) {
	pool, ok := pools[key]
	if !ok {
		return nil, false
	}
	return pool.Get(), true
}

func GetRedis() *redis.Client {
	pool, ok := pools[keyRedis]
	if !ok {
		return nil
	}
	return pool.Get().(*redis.Client)
}

//Put return back the resource to the pool
func Put(key string, resource interface{}) bool {
	pool, ok := pools[key]
	if !ok {
		return false
	}
	pool.Put(resource)
	return true
}

func PutRedis(resource interface{}) bool {
	pool, ok := pools[keyRedis]
	if !ok {
		return false
	}
	pool.Put(resource)
	return true
}
