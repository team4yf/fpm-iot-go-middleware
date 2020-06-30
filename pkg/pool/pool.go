package pool

import (
	"context"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
)

var pool sync.Pool

var TIMEOUT_CTX = context.Background()

func Init(opt *redis.Options) {
	pool = sync.Pool{

		New: func() interface{} {
			return redis.NewClient(opt)
		},
	}

}

func Get() *redis.Client {
	cli := pool.Get().(*redis.Client)
	_, err := cli.Ping(TIMEOUT_CTX).Result()
	if err != nil {
		log.Fatal("redis cant connect ", err)
	}
	return cli
}

func Put(resource interface{}) {
	pool.Put(resource)
}
