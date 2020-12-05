package utils

import (
	"sync"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var redis_once sync.Once

// GetDB : return databast instance
func GetRedis() *redis.Client {
	redis_once.Do(func() {
		rdb = redis.NewClient(getRedisOption())
	})
	return rdb
}
