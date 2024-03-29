package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var Db *gorm.Db

func Init() {
	// 配置 Mysql 连接参数
	username := "root"  // 账号
	password := "0801"  // 密码
	host := "127.0.0.1" // ip
	port := 3306        // 端口
	Dbname := "douyin"  // 数据库
	// 连接参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("db connected error", err)
	} else {
		fmt.Println("db connected successful")
	}
}