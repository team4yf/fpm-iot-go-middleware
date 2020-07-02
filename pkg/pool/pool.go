package pool

import (
	"context"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/team4yf/fpm-iot-go-middleware/external/device/light/lintai"
	"github.com/team4yf/fpm-iot-go-middleware/external/rest"
)

var pools map[string]sync.Pool

//DefaultCtx the default context for redis session.
var DefaultCtx = context.Background()

func init() {
	pools = make(map[string]sync.Pool)
}

//InitRedis initionlize the resouce pool
func InitRedis(opt *redis.Options) {
	pools["redis"] = sync.Pool{

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

func InitRestClient(opt *rest.LinTaiOptions) {
	pools["rest:"+opt.AppID] = sync.Pool{
		New: func() interface{} {
			return lintai.NewClient(opt)
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

//Put return back the resource to the pool
func Put(key string, resource interface{}) bool {
	pool, ok := pools[key]
	if !ok {
		return false
	}
	pool.Put(resource)
	return true
}
