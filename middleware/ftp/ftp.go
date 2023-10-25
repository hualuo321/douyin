package ftp

import (
	"github.com/hualuo321/config"
	"github.com/dutchcoders/goftp"
	"log"
	"time"
)

var MyFTP *goftp.FTP

func InitFTP() {
	var err error
	MyFTP, err = goftp.Connect(config.ConConfig)
	if err != nil {
		log.Printf("获取TCP链接失败")
	}
	log.Printf("获取TCP链接成功")
	// 登录
	err = MyFTP.Login(config.FtpUser, config.FtpPsw)
	if err != nil {
		log.Printf("FTP登录失败")
	}
	log.Printf("FTP登录成功")
	// 维持长链接
	go keepAlive()
}

func keepAlive() {
	time.Sleep(time.Duration(config.HeartbeatTime) * time.Second)
	// 定期发送 FTP 的 NOOP（No Operation）命令以保持 FTP 连接的活跃状态。
	MyFTP.Noop()
}