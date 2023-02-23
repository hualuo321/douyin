package dao

import (
	"log"
	"sync"

	"gorm.io/gorm"
)

type RelationData struct {
	UserId1  int64
	UserId2  int64
	Relation int8
}

// type UserData struct {
// 	Id            int64  `json:"id,omitempty"`
// 	Name          string `json:"name,omitempty"`
// 	FollowCount   int64  `json:"follow_count,omitempty"`
// 	FollowerCount int64  `json:"follower_count,omitempty"`
// 	IsFollow      bool   `json:"is_follow,omitempty"`
// }

// Follow 用户关系结构，对应用户关系表。

//	func (p UserData) TableName() string {
//		return "user"
//	}
func (p RelationData) TableName() string {
	return "follow_detail"
}

// FollowDao 把dao层看成整体，把dao的curd封装在一个结构体中。
type RelationDao struct {
}

var (
	relationDao  *RelationDao //操作该dao层crud的结构体变量。
	relationOnce sync.Once    //单例限定，去限定申请一个followDao结构体变量。
)

// NewFollowDaoInstance 生成并返回RelationDao的单例对象。
func NewRelationDaoInstance() *RelationDao {
	relationOnce.Do(
		func() {
			relationDao = &RelationDao{}
		})
	return relationDao
}

/* CRUD实现 */
// // (id用户关注数+1) 给user表中的对应id的FollowCount +1
// func (*RelationDao) PlusFollowCount(id int64) error {
// 	err := Db.Model(&UserData{Id: id}).Update("follow_count", gorm.Expr("follow_count + 1")).Error
// 	return err
// }

// // (id用户关注数-1) 给user表中的对应id的FollowCount-1
// func (*RelationDao) SubFollowCount(id int64) error {
// 	err := Db.Model(&UserData{Id: id}).Update("follow_count", gorm.Expr("follow_count - 1")).Error
// 	return err
// }

// // (id用户粉丝数+1) 给给user表中的对应id的FollowerCount+1
// func (*RelationDao) PlusFollowerCount(id int64) error {
// 	err := Db.Model(&UserData{Id: id}).Update("follower_count", gorm.Expr("follower_count + 1")).Error
// 	return err
// }

// // (id用户粉丝数-1) 给user表中的对应id的FollowerCount-1
// func (*RelationDao) SubFollowerCount(id int64) error {
// 	err := Db.Model(&UserData{Id: id}).Update("follower_count", gorm.Expr("follower_count - 1")).Error
// 	return err
// }

func (*RelationDao) FindFollowId(userId int64) ([]int64, error) {
	var ids []int64
	err := db.Model(RelationData{}).Where("user_id2=?", userId).Pluck("user_id1", &ids).Error
	log.Println(err)
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	// 查询成功。
	return ids, nil
}

// (找到id用户的所有粉丝id ) 搜索follow_detail表中给定user_id2对应的所有行数据
func (*RelationDao) FindFollowerId(userId int64) ([]int64, error) {
	var ids []int64
	err := db.Model(RelationData{}).Where("user_id1=?", userId).Pluck("user_id2", &ids).Error
	log.Println(err)
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	// 查询成功。
	return ids, nil
}

// (找到id用户的所有关注者id ) 搜索follow_detail表中给定user_id1对应的所有行数据
func (*RelationDao) FindFollows(userId int64) ([]RelationData, error) {
	var rels []RelationData
	err := db.Find(&rels, "user_id1 = ?", userId).Error
	return rels, err
}

// (计数id用户的所有关注者 )
func (*RelationDao) CountFollows(userId int64) (int64, error) {
	var count int64
	err := db.Model(&RelationData{}).Where("user_id1 = ?", userId).Count(&count).Error
	return count, err
}

// (找到id用户的所有粉丝id ) 搜索follow_detail表中给定user_id2对应的所有行数据
func (*RelationDao) FindFollowers(userId int64) ([]RelationData, error) {
	var rels []RelationData
	err := db.Find(&rels, "user_id2 = ?", userId).Error
	return rels, err
}

// (计数id用户的所有粉丝 )
func (*RelationDao) CountFollowers(userId int64) (int64, error) {
	var count int64
	err := db.Model(&RelationData{}).Where("user_id2 = ?", userId).Count(&count).Error
	return count, err
}

// (找到id用户的所有关系id ) 搜索follow_detail表中userid1=id or userid2=id对应的所有行数据
func (*RelationDao) FindRelations(userId int64) ([]RelationData, error) {
	var rels []RelationData
	err := db.Find(&rels, "user_id1 = ? or user_id2=?", userId, userId).Error
	return rels, err
}

// 判断userId有没有关注toUserId
func (*RelationDao) FindRelationBetween(userId int64, toUserId int64) (bool, error) {
	var rel RelationData
	var ok bool
	err := db.Find(&rel, "user_id1 = ? and user_id2=?", userId, toUserId).Error
	if rel == (RelationData{}) { // 结构体变量判空
		ok = false
	} else {
		ok = true
	}
	return ok, err
}

// // 搜索user表中给定id用户的全部信息
// func (*RelationDao) FindUserById(id int64) (UserData, error) {
// 	var u UserData
// 	err := Db.Find(&u, "id = ?", id).Error
// 	return u, err
// }

// (新增关系 userId用户关注了toUserId用户) 给follow_detail表中新增一条关系数据(行)
func (*RelationDao) CreateRelationInfo(userId int64, toUserId int64) error {
	err := db.Create(&RelationData{UserId1: userId, UserId2: toUserId, Relation: 1}).Error
	return err
}

// (删除关系 userId用户取关了toUserId用户) 给follow_detail表中删除一条关系数据(行)
func (*RelationDao) DeleteRelationInfo(userId int64, toUserId int64) error {
	err := db.Where("user_id1 = ? and user_id2=?", userId, toUserId).Delete(RelationData{}).Error
	return err
}
