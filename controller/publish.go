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

// response类型
type FeedResponse struct {
	Response
	VideoList []service.VideoData `json:"video_list,omitempty"`
	NextTime  int64               `json:"next_time,omitempty"`
}

// Publish /publish/action/
func Publish(c *gin.Context) {
	title := c.PostForm("title")
	// fmt.Println("here is title:", title)

	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	fmt.Println("here is 获取到用户id:", userId)

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
	fmt.Println("filename:", filename) // rui_shi.mp4
	filename = fmt.Sprintf("%d_%s", userId, filename)
	fmt.Println("filename:", filename)                      // 2_rui_shi.mp4
	saveFile := filepath.Join("./public/videos/", filename) // ./public/videos/2_rui_shi.mp4
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 传一个视频名字，一个title，一个userID
	newVideo, err := videoi.PublishVideo(filename, title, userId) // 2_rui_shi.mp4
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
	fmt.Println("---当前位于 /controller/video/Feed()")
	inputTime := c.Query("latest_time")
	fmt.Println("-传入的时间戳", inputTime)
	var lastTime time.Time
	if inputTime != "0" {
		fmt.Println("这里是1")
		me, _ := strconv.ParseInt(inputTime, 10, 64)
		lastTime = time.UnixMilli(me)
	} else {
		fmt.Println("这里是2")
		lastTime = time.Now()
	}
	fmt.Println("-最新的时间", lastTime)
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64) // 当前 ID
	fmt.Println("-获取到当前用户id", userId)

	fmt.Println("-定义videoi接口, 根据 videoi.Feed 获取视频数据列表")
	videoi := service.VideoImpl{}
	videoDataList, nextTime, err := videoi.Feed(lastTime, userId)

	if err != nil {
		log.Printf("方法 videoi.Feed(lastTime, userId) 失败：%v\n", err)
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "获取视频流失败"},
		})
		return
	}
	log.Printf("方法videoService.Feed(lastTime, userId) 成功")
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videoDataList,
		NextTime:  nextTime.Unix(),
	})
}
