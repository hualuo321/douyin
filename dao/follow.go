package dao

import (
	"log"
	"sync"

	"gorm.io/gorm"
)

type Follow struct {
	UserId   int64 `gorm:"column:user_id"`
	ToUserId int64 `gorm:"column:to_user_id"`
}

// Follow 用户关系结构，对应用户关系表。
func (Follow) TableName() string {
	return "follow"
}

// FollowDao 把dao层看成整体，把dao的curd封装在一个结构体中。
type FollowDao struct {
}

var (
	followDao  *FollowDao //操作该dao层crud的结构体变量。
	followOnce sync.Once  //单例限定，去限定申请一个followDao结构体变量。
)

// NewFollowDaoInstance 生成并返回RelationDao的单例对象。
func NewFollowDaoInstance() *FollowDao {
	followOnce.Do(
		func() {
			followDao = &FollowDao{}
		})
	return followDao
}

// 添加一条 follow 记录
func (*FollowDao) InsertFollow(userId int64, toUserId int64) error {
	// 新建关系
	follow := Follow{
		UserId:   userId,
		ToUserId: toUserId,
	}
	if err := db.Create(&follow).Error; err != nil {
		log.Println("创建 follow 关系失败", err.Error())
		return err
	}
	return nil
}

// 删除一条 follow 记录
func (*FollowDao) DeleteFollow(userId int64, toUserId int64) error {
	// delete follow record in follows
	var follow Follow
	err := db.Where("user_id = ? && to_user_id = ?", userId, toUserId).Delete(&follow).Error
	return err
}

// 查询两者之间的关系
func QueryFollowInfo(userId int64, toUserId int64) (bool, error) {
	var count int64 = 0
	err := db.Model(&Follow{}).Where(Follow{UserId: userId, ToUserId: toUserId}).Count(&count).Error
	if err != nil {
		log.Println("查询关系失败:" + err.Error())
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

// 统计订阅数
func CountFollow(userId int64) (int64, error) {
	var count int64
	err := db.Model(&Follow{}).Where("user_id = ?", userId).Count(&count).Error
	if err != nil {
		log.Println("查询订阅数失败:" + err.Error())
		return 0, err
	}
	return count, nil
}

// 统计粉丝数
func CountFollower(toUserId int64) (int64, error) {
	var count int64
	err := db.Model(&Follow{}).Where("to_user_id = ?", toUserId).Count(&count).Error
	if err != nil {
		log.Println("查询粉丝数失败:" + err.Error())
		return 0, err
	}
	return count, nil
}

// 获取订阅列表
func GetFollowIdList(userId int64) ([]int64, error) {
	var followIdList []int64
	err := db.Model(&Follow{}).Where("user_id = ?", userId).Select("to_user_id").Find(&followIdList).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		log.Println("查询订阅列表失败:" + err.Error())
		return nil, err
	}
	return followIdList, nil
}

// 获取粉丝列表
func GetFollowerIdList(toUserId int64) ([]int64, error) {
	var followerIdList []int64
	err := db.Model(&Follow{}).Where("to_user_id = ?", toUserId).Select("user_id").Find(&followerIdList).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		log.Println("查询粉丝列表失败:" + err.Error())
		return nil, err
	}
	return followerIdList, nil
}
