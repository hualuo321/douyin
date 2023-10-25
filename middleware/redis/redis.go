package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

// 这段代码是一个Go语言程序中用于初始化连接到Redis数据库的代码片段。

var Ctx = context.Background()
// 用于连接到 Redis 数据库的客户端对象，每个客户端对象用于不同的数据库操作
var RdbFollowers *redis.Client
var RdbFollowing *redis.Client
var RdbFollowingPart *redis.Client

// 初始化 Redis 连接
func InitRedis() {
	RdbFollowers = redis.NewClient(&redis.Options {
		Addr:		"127.0.0.1:6379",
		Password:	"123456",
		DB:			0
	})
	RdbFollowing = redis.NewClient(&redis.Options{
		Addr:     	"127.0.0.1:6379",
		Password: 	"123456",
		DB:       	1, // 关注列表信息信息存入 DB1.
	})
	RdbFollowingPart = redis.NewClient(&redis.Options{
		Addr:     	"127.0.0.1:6379",
		Password: 	"123456",
		DB:       	3, // 当前用户是否关注了自己粉丝信息存入 DB3.
	})
}