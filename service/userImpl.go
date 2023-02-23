package service

import (
	"fmt"
	"log"

	"github.com/hualuo321/douyin/dao"
)

type UserImpl struct {
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
	user, err := dao.GetUserById(id)
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
func (useri *UserImpl) GetUserDataById(id int64) (UserData, error) {
	userData := UserData{
		Id:             0,
		Name:           "",
		FollowCount:    0,
		FollowerCount:  0,
		IsFollow:       false,
		Signature:      "",
		TotalFavorited: 0,
		FavoritedCount: 0,
	}
	user, err := dao.GetUserById(id)
	if err != nil {
		log.Println("err:", err.Error())
		log.Println("用户未找到！！")
		return userData, err
	}
	// followCount := GetUserFollowInfo(id)
	// fmt.Println("----------------followcount, ", followCount)
	// followerCount := GetUserFollowerInfo(id)
	signature := "世界这么大，谢谢你的关注！"
	dao.NewRelationDaoInstance().FindRelationBetween(id, id)
	userData = UserData{
		Id:             id,
		Name:           user.Username,
		FollowCount:    333,
		FollowerCount:  444,
		IsFollow:       false,
		Signature:      signature,
		TotalFavorited: 100,
		FavoritedCount: 10,
	}

	log.Println("查询用户成功!!")
	//返回封装的结构体
	return userData, nil
}

// 根据当前Id获取用户数据
func (useri *UserImpl) GetUserDataByCurId(id int64, currentId int64) (UserData, error) {
	userData := UserData{
		Id:             0,
		Name:           "",
		FollowCount:    0,
		FollowerCount:  0,
		IsFollow:       false,
		Signature:      "",
		TotalFavorited: 0,
		FavoritedCount: 0,
	}
	user, err := dao.GetUserById(id) // 获取视频作者
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
	followCount := GetUserFollowInfo(id)
	fmt.Println("----------------followcount, ", followCount)
	followerCount := GetUserFollowerInfo(id)
	signature := "世界这么大，谢谢你的关注！"
	isFollow, _ := dao.NewRelationDaoInstance().FindRelationBetween(currentId, id)
	fmt.Println(currentId, id, isFollow)
	//返回封装的结构体
	userData = UserData{
		Id:             id,
		Name:           user.Username,
		FollowCount:    followCount,   //这里本应该是查询出来的变量,为了测试直接赋值
		FollowerCount:  followerCount, //这里本应该是查询出来的变量,为了测试直接赋值
		IsFollow:       isFollow,      //这里因为没有登录，查看不了是否关注，默认不关注
		Signature:      signature,
		TotalFavorited: 100, //这里本应该是查询出来的变量,为了测试直接赋值
		FavoritedCount: 10,  //这里本应该是查询出来的变量,为了测试直接赋值
	}
	return userData, nil
}

// / 获取用户数据列表
// func (useri *UserImpl) GetUserDataList(userIdList []int64) ([]UserData, error) {
// 	var userDataList []UserData
// 	for _, userId := range userIdList {
// 		userData, err := useri.GetUserDataById(userId)
// 		if err != nil {
// 			log.Println("通过Id获取用户失败", err)
// 			return nil, err
// 		}
// 		userDataList = append(userDataList, userData)
// 	}
// 	return userDataList, nil
// }
