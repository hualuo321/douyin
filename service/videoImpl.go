package service

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
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

// 获取视频数据列表
func (videoi *VideoImpl) GetPublishList(userId int64, curId int64) ([]VideoData, error) {
	//依据用户id查询所有的视频，获取视频列表
	videoList, err := dao.QueryVideoListByUserId(userId)
	if err != nil {
		log.Printf("方法dao.QueryVideoListByUserId(%v)失败:%v", userId, err)
		return nil, err
	}
	log.Printf("方法dao.QueryVideoListByUserId(%v)成功:%v", userId, videoList)
	// 获取到的视频列表
	videoDataList, err := videoi.PrepareVideoData(videoList, userId)
	if err != nil {
		log.Printf("方法videoi.PrepareVideoData(videoList, userId)失败")
		return nil, err
	}
	// 如果数据没有问题，则直接返回
	return videoDataList, nil
}

func (videoi *VideoImpl) PrepareVideoData(videoList []dao.Video, userId int64) ([]VideoData, error) {
	videoDataList := make([]VideoData, 0, len(videoList))
	for _, video := range videoList {
		videoData, _ := videoi.CreatVideoData(video, userId)
		videoDataList = append(videoDataList, videoData)
	}
	return videoDataList, nil
}

// 创建 videoData 数据
func (videoi *VideoImpl) CreatVideoData(video dao.Video, userId int64) (VideoData, error) {
	//建立协程组，当这一组的携程全部完成后，才会结束本方法
	var wg sync.WaitGroup
	wg.Add(1) // 当前只有一个进程

	useri := UserImpl{}

	var err error
	var videoData VideoData
	videoData.Id = video.Id
	videoData.PlayUrl = video.PlayUrl
	videoData.CoverUrl = video.CoverUrl
	videoData.Title = video.Title
	//插入Author，这里需要将视频的发布者和当前登录的用户传入，才能正确获得isFollow，
	//如果出现错误，不能直接返回失败，将默认值返回，保证稳定
	go func() {
		videoData.Author, err = useri.GetUserByCurrentId(video.Id, userId)
		if err != nil {
			log.Printf("方法videoi.GetUserByCurrentId(video.Id, userId)失败：%v", err)
		} else {
			log.Printf("方法videoi.GetUserByCurrentId(video.Id, userId)成功")
		}
		wg.Done()
	}()

	//插入点赞数量, 待实现
	//获取该视屏的评论数字
	//获取当前用户是否点赞了该视频

	wg.Wait()
	return videoData, err
}
