package service

import (
	"fmt"
	"log"

	"github.com/hualuo321/douyin/dao"
)

type UserImpl struct {
	UserService
}

func (useri *UserImpl) QueryAllUser() []dao.User {
	users, err := dao.QueryAllUser()
	if err != nil {
		log.Println("error:" + err.Error())
		return users
	}
	return users
}

// 根据id查找用户
func (useri *UserImpl) QueryUserById(id int64) dao.User {
	user, err := dao.QueryUserById(id)
	if err != nil {
		log.Println("error:", err.Error())
		log.Println("通过id查询用户失败! 用户未找到!")
		return user
	}
	log.Println("查询用户成功！")
	return user
}

// 根据username查找用户
func (useri *UserImpl) QueryUserByUsername(name string) dao.User {
	fmt.Println("---当前位于service/userImpl/QueryUserByUsername()---")
	user, err := dao.QueryUserByUsername(name)
	if err != nil {
		log.Println("error:", err.Error())
		log.Println("通过用户名查询用户失败！用户未找到！")
		return user
	}
	log.Println("查询用户成功！")
	fmt.Println("-查询到的用户信息为", user)
	return user
}

// 新增用户
func (useri *UserImpl) InsertUser(user *dao.User) bool {
	flag := dao.InsertUser(user)
	if !flag {
		log.Println("新增用户失败!!")
		return false
	}
	return true
}

// 获取用户数据 UserData
// 未登录状态
func (useri *UserImpl) GetUserById(id int64) (UserData, error) {
	userData := UserData{
		Id:             0,
		Username:       "",
		FollowCount:    0,
		FollowerCount:  0,
		IsFollow:       false,
		TotalFavorited: 0,
		FavoritedCount: 0,
	}
	user, err := dao.QueryUserById(id)
	if err != nil {
		log.Println("err:", err.Error())
		log.Println("用户未找到！！")
		return userData, err
	}
	log.Println("查询用户成功!!")
	//这里查询关注数量
	//查询粉丝数量
	//查询总的点赞数
	//查询喜欢的视频数

	//返回封装的结构体
	userData = UserData{
		Id:             id,
		Username:       user.Username,
		FollowCount:    3,     //这里本应该是查询出来的变量,为了测试直接赋值
		FollowerCount:  5,     //这里本应该是查询出来的变量,为了测试直接赋值
		IsFollow:       false, //这里因为没有登录，查看不了是否关注，默认不关注
		TotalFavorited: 100,   //这里本应该是查询出来的变量,为了测试直接赋值
		FavoritedCount: 1,     //这里本应该是查询出来的变量,为了测试直接赋值
	}
	return userData, nil
}

// 已登录状态
func (useri *UserImpl) GetUserByCurrentId(id int64, currentId int64) (UserData, error) {
	userData := UserData{
		Id:             0,
		Username:       "",
		FollowCount:    0,
		FollowerCount:  0,
		IsFollow:       false,
		TotalFavorited: 0,
		FavoritedCount: 0,
	}
	user, err := dao.QueryUserById(id)
	if err != nil {
		log.Println("err:", err.Error())
		log.Println("用户未找到！！")
		return userData, err
	}
	log.Println("查询用户成功!!")
	//这里查询关注数量
	//查询粉丝数量
	//查询是否关注 利用当前id 和 目标id
	//查询总的点赞数
	//查询喜欢的视频数

	//返回封装的结构体
	userData = UserData{
		Id:             id,
		Username:       user.Username,
		FollowCount:    3,     //这里本应该是查询出来的变量,为了测试直接赋值
		FollowerCount:  5,     //这里本应该是查询出来的变量,为了测试直接赋值
		IsFollow:       false, //这里因为没有登录，查看不了是否关注，默认不关注
		TotalFavorited: 100,   //这里本应该是查询出来的变量,为了测试直接赋值
		FavoritedCount: 1,     //这里本应该是查询出来的变量,为了测试直接赋值
	}
	return userData, nil
}
