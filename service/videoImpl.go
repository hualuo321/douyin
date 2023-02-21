package service

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/hualuo321/douyin/dao"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type VideoImpl struct {
}

// 上传视频
func (videoi *VideoImpl) PublishVideo(filename string, title string, userId int64) (*dao.Video, error) {
	playUrl := filename
	//添加生成视频关键帧并上传到public目录的函数
	_, err := GetSnapshot("./public/videos/"+filename, "./public/covers/"+filename, 1)
	if err != nil {
		log.Println("generate cover err:" + err.Error())
		return nil, err
	}
	coverUrl := strings.TrimSuffix(filename, ".mp4") + ".jpg"
	newVideo := dao.Video{PlayUrl: playUrl, CoverUrl: coverUrl, Title: title, UserId: userId, PublishTime: time.Now()}
	// func InsertVideo(videoName string, imageName string, userId int64, title string) error
	if err := dao.InsertVideo(playUrl, coverUrl, userId, title); err != nil {
		log.Println("post video to db err:" + err.Error())
		return nil, err
	}
	return &newVideo, nil
}

// 获取视频中的关键帧当作封面
func GetSnapshot(videoPath, snapshotPath string, frameNum int) (snapshotName string, err error) {
	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		log.Println("generate cover fail:" + err.Error())
		return "", err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		log.Println("decoding cover fail:" + err.Error())
		return "", err
	}

	err = imaging.Save(img, snapshotPath+".jpg")
	if err != nil {
		log.Println("saving cover fail:" + err.Error())
		return "", err
	}
	names := strings.Split(snapshotPath, "\\")
	snapshotName = names[len(names)-1] + ".jpg"
	return snapshotName, nil
}
