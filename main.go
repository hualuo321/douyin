package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hualuo321/douyin/dao"
	"github.com/hualuo321/douyin/redisDb"
)

func main() {

	// 初始化数据库
	if err := dao.Init(); err != nil {
		// Exit 函数可以让当前程序以给出的状态码 code 退出。一般来说，状态码 0 表示成功，非 0 表示出错。
		os.Exit(-1)
	}
	// Redis初始化
	if err := redisDb.InitRedis(); err != nil {
		log.Printf("connect redis failed! err : %v\n", err)
		os.Exit(-1)
	}

	// gin 声明
	r := gin.Default()
	// gin 初始化
	initRouter(r)
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run()

}
