package database

import (
	"cathub.me/go-web-examples/pkg/setting"
	"context"
	"github.com/go-redis/redis/v8"
	"sync"
)

var _redisClient *redis.Client
var _onceRedisClient sync.Once

func GetRedisClient() *redis.Client {
	_onceRedisClient.Do(func() {
		_redisClient = redis.NewClient(&redis.Options{
			Addr:     setting.Redis.Addr,
			Username: setting.Redis.Username,
			Password: setting.Redis.Password,
			DB:       setting.Redis.DB,
		})

		_redisClient.Conn(context.Background())
	})
	return _redisClient
}
