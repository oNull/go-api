package helpers

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"py-mxshop-api/user_web/global"
)

var rdb *redis.Client

func InitRedis(db int) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
		Password: global.ServerConfig.RedisInfo.Password, // no password set
		DB:       db,                                     // use default DB
	})
}

func RedisSet(ctx context.Context, key string, value interface{}) error {
	err := rdb.Set(ctx, key, value, 0).Err()
	return handleRedisError(err)
}

func RedisGet(ctx context.Context, key string) (string, error) {
	value, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return "", handleRedisError(err)
	}
	return value, nil
}

// 其他Redis操作函数（例如删除键、检查键是否存在等）可以根据需要类似地实现

func handleRedisError(err error) error {
	if err == nil {
		// 通知钉钉或QQ Redis报警
		return nil
	}
	if err == redis.Nil {
		return ErrKeyNotFound
	}
	return ErrRedisError
}

// ErrRedisError 自定义错误类型
var ErrRedisError = redisError("Redis operation error")
var ErrKeyNotFound = redisError("Key not found")

type redisError string

func (e redisError) Error() string {
	return string(e)
}
