package service

import "github.com/hualuo321/douyin/dao"

type VideoData struct {
	Id       int64    `json:"id,omitempty"`
	Author   UserData `json:"author"`
	PlayUrl  string   `json:"play_url,omitempty"`
	CoverUrl string   `json:"cover_url,omitempty"`
	Title    string   `json:"title,omitempty"`
}

type VideoService interface {
	// 上传视频
	PublishVideo(filename string, title string, userId int64) (*dao.Video, error)
	// 获取某用户视频列表
	GetPublishList(userID int64, curId int64) ([]VideoData, error)
	// 获取所有视频列表, 及最近视频发布时间
	Feed(lastTime int64, userId int64) ([]VideoData, int64, error)
}
