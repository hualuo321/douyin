package redisDb

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()
var RdbMessageHelper *redis.Client
var RdbUserInfo *redis.Client

func InitRedis() (err error) {
	RdbMessageHelper = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0, // Message消息信息缓存存入 DB0：用于聊天记录更新.
	})
	_, err = RdbMessageHelper.Ping(Ctx).Result()

	// 初始化用不用清空缓存

	RdbUserInfo = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       1, // user用户信息缓存存入 DB1：用于缓存用户关注数粉丝数.
	})
	_, err = RdbUserInfo.Ping(Ctx).Result()
	return
}
