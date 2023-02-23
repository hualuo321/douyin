package service

import (
	"log"

	"github.com/hualuo321/douyin/dao"
)

type FollowImpl struct {
	UserService
}

// 添加订阅记录
func InsertFollow(userId int64, toUserId int64) error {
	err := dao.NewFollowDaoInstance().InsertFollow(userId, toUserId)
	if err != nil {
		log.Println("添加订阅记录失败!", err)
		return err
	}
	return nil
}

// 取消订阅记录
func DeleteFollow(userId int64, toUserId int64) error {
	err := dao.NewFollowDaoInstance().DeleteFollow(userId, toUserId)
	if err != nil {
		log.Println("取消订阅记录失败!", err)
		return err
	}
	return nil
}

func GetFollowDataList(userId int64) ([]UserData, error) {
	followIdList, err := dao.GetFollowIdList(userId)
	if err != nil {
		log.Println("获取订阅Id列表失败!", err)
		return nil, err
	}
	userDataList, err := new(UserImpl).GetUserDataList(followIdList)
	if err != nil {
		log.Println("获取订阅UserData列表失败!", err)
		return nil, err
	}
	return userDataList, nil
}

func GetFollowerDataList(userId int64) ([]UserData, error) {
	followerIdList, err := dao.GetFollowerIdList(userId)
	if err != nil {
		log.Println("获取粉丝Id列表失败!", err)
		return nil, err
	}
	userDataList, err := new(UserImpl).GetUserDataList(followerIdList)
	if err != nil {
		log.Println("获取粉丝UserData列表失败!", err)
		return nil, err
	}
	return userDataList, nil
}

func GetUserDataById(userId int64, isFollow bool) (UserData, error) {
	var userData = UserData{}
	userData, _ = new(UserImpl).GetUserDataById(userId)
	userData.FavoritedCount, _ = dao.CountFollow(userId)
	userData.FollowerCount, _ = dao.CountFollower(userId)
	userData.IsFollow = isFollow
	return userData, nil
}
