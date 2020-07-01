package pool

import (
	"context"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
)

var pool sync.Pool

//DefaultCtx the default context for redis session.
var DefaultCtx = context.Background()

//Init initionlize the resouce pool
func Init(opt *redis.Options) {
	pool = sync.Pool{

		New: func() interface{} {
			return redis.NewClient(opt)
		},
	}

}

//Get get a resource from the pool
func Get() *redis.Client {
	cli := pool.Get().(*redis.Client)
	_, err := cli.Ping(DefaultCtx).Result()
	if err != nil {
		log.Fatal("redis cant connect ", err)
	}
	return cli
}

//Put return back the resource to the pool
func Put(resource interface{}) {
	pool.Put(resource)
}
