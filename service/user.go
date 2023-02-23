package service

import "github.com/hualuo321/douyin/dao"

// 封装返回的User结构体信息
type UserData struct {
	Id             int64  `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	FollowCount    int64  `json:"follow_count"`
	FollowerCount  int64  `json:"follower_count"`
	IsFollow       bool   `json:"is_follow"`
	Signature      string `json:"signature"`
	TotalFavorited int64  `json:"total_favorited,,omitempty"`
	FavoritedCount int64  `json:"favorited_count,omitempty"`
}

type UserService interface {
	//获得全部对象
	QueryAllUser() []dao.User
	//根据id查找用户
	QueryUserById(id int64) dao.User
	//根据用户名查找user
	QueryUserByUsername(name string) dao.User
	//新增用户
	InsertUser(tableUser *dao.User) bool
	//他人查看信息
	//通过id查询，未登录情况
	GetUserDataById(id int64) (UserData, error)
	//已经登录的情况下，根据user_id获得User对象
	GetUserByCurrentId(id int64, currentId int64) (UserData, error)
}
