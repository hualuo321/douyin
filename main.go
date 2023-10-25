package main

import (
	"github.com/hualuo321/douyin/dao"
	"github.com/hualuo321/douyin/middleware/ffmpeg"
	"github.com/hualuo321/douyin/middleware/ftp"
	"github.com/hualuo321/douyin/middleware/rabbitmq"
	"github.com/hualuo321/douyin/middleware/redis"
	"github.com/hualuo321/douyin/util"
	"github.com/gin-gonic/gin"
)

// 函数起始
func main() {
	// 初始化依赖
	initDeps()
	// gin 声明
	r := gin.Default()
	// 初始化路由
	initRouter(r)
	// listen and serve on 0.0.0.0:8080
	r.Run()
}

// 加载项目依赖
func initDeps() {
	// 初始化数据库
	dao.Init()
	// 初始化 FTP 服务器
	ftp.InitFTP()
	// 初始化 SSH
	ffmpeg.InitSSH()
	// 初始化 Redis
	redis.InitRedis()
	// 初始化消息队列
	rabbitmq.InitRabitMQ()
	// 初始化 Follow 的消息队列，并开启消费
	rabbitmq.InitFollowRabbitMQ()
}
