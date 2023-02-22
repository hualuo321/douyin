package controller

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hualuo321/douyin/service"
)

type VideoListResponse struct {
	Response
	VideoList []service.VideoData `json:"video_list"`
}

type FeedResponse struct {
	Response
	VideoList []service.VideoData `json:"video_list"`
	NextTime  int64               `json:"next_time"`
}

// Publish /publish/action/
func Publish(c *gin.Context) {
	title := c.PostForm("title")
	fmt.Println(title)

	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	fmt.Println("获取到用户id:", userId)

	videoi := service.VideoImpl{}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	print("filename:", filename)
	filename = fmt.Sprintf("%d_%s", userId, filename)
	print("filename:", filename)
	saveFile := filepath.Join("./public/videos/", filename)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 传一个视频名字，一个title，一个userID
	newVideo, err := videoi.PublishVideo(filename, title, userId)
	if newVideo == nil || err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 2,
			StatusMsg:  "save post video to db fail",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  filename + " uploaded successfully",
	})

}

// PublishList /publish/list/
func PublishList(c *gin.Context) {

	// 查询的别人的 uerId	为 1
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	fmt.Printf("获取到用户id:%v\n", userId)
	// 当前登陆人的 userId  为 1
	curId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	fmt.Printf("获取到当前用户id:%v\n", curId)

	videoi := service.VideoImpl{}
	videoDataList, err := videoi.GetPublishList(userId, curId)
	if err != nil {
		fmt.Printf("调用PublishList(%v)出现错误：%v\n", userId, err)
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "获取视频列表失败"},
		})
		return
	}

	fmt.Printf("调用GetPublishList(%v)成功", userId)
	c.JSON(http.StatusOK, VideoListResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videoDataList,
	})
}

// 拉取视频列表
func Feed(c *gin.Context) {
	// 拉去近期的部分视频
	queryTime := c.Query("lasted_time")
	fmt.Println("传入的时间:", queryTime)
	var lastTime int64
	if queryTime != "0" {
		lastTime, _ = strconv.ParseInt(queryTime, 10, 64)
	} else {
		lastTime = time.Now().Unix()
	}
	fmt.Printf("获取到时间戳%v", lastTime)

	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	fmt.Printf("获取到用户id:%v\n", userId)
	videoi := service.VideoImpl{}
	feed, nextTime, err := videoi.Feed(lastTime, userId)

	if err != nil {
		log.Printf("方法videoService.Feed(lastTime, userId) 失败：%v", err)
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "获取视频流失败"},
		})
		return
	}
	log.Printf("方法videoService.Feed(lastTime, userId) 成功")
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: feed,
		NextTime:  nextTime,
	})
}
