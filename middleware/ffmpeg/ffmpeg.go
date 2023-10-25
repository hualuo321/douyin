package ffmpeg

import (
	"github.com/hualuo321/config"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"time"
)

type Ffmsg struct {
	VideoName string
	ImageName string
}

var ClientSSH *ssh.Client
var Ffchan chan Ffmsg

func InitSSH() {
	var err error
	SSHconfig := &ssh.ClientConfig(
		Timeout:			5 * time.Second,
		User:				config.UserSSH,
		HostKeyCallback:	ssh.InsecureIgnoreHostKey()
	)
	if config.TypeSSH == "password" {
		SSHconfig.Auth = []ssh.AuthMethod{ssh.Password(config.PasswordSSH)}
	}
	// 获取ssh client
	addr := fmt.Sprintf("%s:%d", config.HostSSH, config.PortSSH)
	ClientSSH, err = ssh.Dial("tcp", addr, SSHconfig)
	if err != nil {
		log.Fatal("创建ssh客户端失败", err)
	}
	log.Printf("获取到客户端：%v", ClientSSH)
	// 建立通道，作为队列使用,并且确立缓冲区大小
	Ffchan = make(chan Ffmsg, config.MaxMsgCount)
	// 建立携程用于派遣
	go dispatcher()
	go keepAlive()
}

//通过增加携程，将获取的信息进行派遣，当信息处理失败之后，还会将处理方式放入通道形成的队列中
func dispatcher() {
	for ffmsg := range Ffchan {
		go func(f Ffmsg) {
			err := Ffmpeg(f.VideoName, f.ImageName)
			if err != nil {
				Ffchan <- f
				log.Fatal("派遣失败：重新派遣")
			}
			log.Printf("视频%v截图处理成功", f.VideoName)
		}(ffmsg)
	}
}

// Ffmpeg 通过远程调用ffmpeg命令来创建视频截图
func Ffmpeg(videoName string, imageName string) error {
	session, err := ClientSSH.NewSession()
	if err != nil {
		log.Fatal("创建ssh session 失败", err)
	}
	defer session.Close()
	//执行远程命令 ffmpeg -ss 00:00:01 -i /home/ftpuser/video/1.mp4 -vframes 1 /home/ftpuser/images/4.jpg
	combo, err := session.CombinedOutput("ls;/usr/local/ffmpeg/bin/ffmpeg -ss 00:00:01 -i /home/ftpuser/video/" + videoName + ".mp4 -vframes 1 /home/ftpuser/images/" + imageName + ".jpg")
	if err != nil {
		log.Fatal("远程执行cmd失败", err)
		log.Fatal("命令输出:", string(combo))
	}
	return nil
}

// 维持长链接
func keepAlive() {
	time.Sleep(time.Duration(config.SSHHeartbeatTime) * time.Second)
	session, _ := ClientSSH.NewSession()
	session.Close()
}
