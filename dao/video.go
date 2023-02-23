package dao

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

var (
	VideoCount     = 5
	PlayUrlPrefix  = "http://192.168.249.81:8080/public/videos/"
	CoverUrlPrefix = "http://192.168.249.81:8080/public/covers/"
)

// 视频 表结构体
type Video struct {
	Id          int64     `gorm:"column:id"`
	UserId      int64     `gorm:"column:user_id"`
	PlayUrl     string    `gorm:"column:play_url"`
	CoverUrl    string    `gorm:"column:cover_url"`
	PublishTime time.Time `gorm:"column:publish_time"`
	Title       string    `gorm:"column:title"`
}

func (Video) TableName() string {
	return "video"
}

// 根据 UserId 查询 视频列表
func QueryVideoListByUserId(userId int64) ([]Video, error) {
	var videos []Video
	err := db.Model(&Video{}).Where("user_id=?", userId).Find(&videos).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		log.Println("query video list by user id error", err)
		return nil, err
	}
	return videos, nil
}

// 依据 视频id 查询 视频
func QueryVideoByVideoId(videoId int64) (Video, error) {
	var video Video
	video.Id = videoId
	//Init()
	err := db.First(&video).Error

	if err != nil {
		return video, err
	}
	return video, nil
}

// 依据一个时间，来获取这个时间之前的一些视频
func QueryVideosByLastTime(lastTime time.Time) ([]Video, error) {
	videoList := make([]Video, VideoCount)
	//format the time to compare with the time in db
	err := db.Where("publish_time<?", lastTime).Order("publish_time desc").Limit(VideoCount).Find(&videoList).Error
	// err := db.Model(&Video{}).Where("publish_time < ?", lastTime).Order("publish_time desc").Find(&videoList).Error
	fmt.Println("length of video list: ", len(videoList))
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return videoList, err
	}
	return videoList, nil
}

// 插入视频记录到db
func InsertVideo(videoName string, imageName string, userId int64, title string) error {
	var video Video
	video.UserId = userId
	video.PlayUrl = PlayUrlPrefix + videoName
	video.CoverUrl = CoverUrlPrefix + imageName
	video.PublishTime = time.Now()
	video.Title = title
	err := db.Save(&video).Error
	if err != nil {
		return err
	}
	return nil
}

// GetVideoIdsByAuthorId
// 通过作者id来查询发布的视频id切片集合
func QueryVideoIdsByAUserId(userId int64) ([]int64, error) {
	var ids []int64
	//通过pluck来获得单独的切片
	err := db.Model(&Video{}).Where("user_id", userId).Pluck("id", &ids).Error
	//如果出现问题，返回对应到空，并且返回error
	if err != nil {
		return nil, err
	}
	return ids, nil
}
