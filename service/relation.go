package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hualuo321/douyin/dao"
	"github.com/hualuo321/douyin/redisDb"
)

/* 方法实现 */

type FriendUser struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`

	Message string `json:"message,omitempty"`
	MsgType int8   `json:"msgType,omitempty"`
}

func GetUserFollowInfo(id int64) int64 {
	/*
		写入：1. 更新数据库，删掉缓存 --> DoFollow
		读取：1. 缓存命中？是：读取，否：从数据库读取写入缓存
	*/
	idStr := strconv.FormatInt(id, 10)
	follow, err := redisDb.RdbUserInfo.HGet(redisDb.Ctx, "follow", idStr).Result()
	if err != nil && (!errors.Is(err, redis.Nil)) {
		fmt.Println(err)
	}
	if follow != "" { // 缓存命中
		followCount, _ := strconv.ParseInt(follow, 10, 64)
		return followCount
	}
	// 缓存未命中
	followCount, err := dao.NewRelationDaoInstance().CountFollows(id)
	err = redisDb.RdbUserInfo.HSet(redisDb.Ctx, "follow", idStr, followCount).Err()
	if err != nil {
		panic(err)
	}
	fmt.Println(followCount)
	return followCount
}
func GetUserFollowerInfo(id int64) int64 {
	/*
		写入：1. 更新数据库，删掉缓存 --> CounselFollow
		读取：1. 缓存命中？是：读取，否：从数据库读取写入缓存
	*/
	idStr := strconv.FormatInt(id, 10)
	follower, err := redisDb.RdbUserInfo.HGet(redisDb.Ctx, "follower", idStr).Result()
	if err != nil && (!errors.Is(err, redis.Nil)) {
		fmt.Println(err)
	}
	if follower != "" { // 缓存命中
		fmt.Println("缓存命中")
		followerCount, _ := strconv.ParseInt(follower, 10, 64)
		return followerCount
	}
	// 缓存未命中
	log.Println("缓存未命中")
	followerCount, err := dao.NewRelationDaoInstance().CountFollowers(id)
	err = redisDb.RdbUserInfo.HSet(redisDb.Ctx, "follower", idStr, followerCount).Err()
	if err != nil {
		panic(err)
	}
	fmt.Println(followerCount)
	return followerCount
}

// id1用户关注id2用户逻辑实现
func DoFollow(id1 int64, id2 int64) bool {
	/* 实现：
	更新数据库
		0. 判断是否已关注，已关注就什么都不做
		// 1. id1用户关注数量增加——user表更新
		// 2. id2用户粉丝数量增加——user表更新
		3. 新建一条id1和id2的关系——relation表新建
	更新缓存
	*/
	relationDao := dao.NewRelationDaoInstance()
	if ok, _ := relationDao.FindRelationBetween(id1, id2); ok {
		return true
	}
	// err := relationDao.PlusFollowCount(id1)
	// if err != nil {
	// 	log.Println("关注逻辑时增加关注数发生了错误：", err)
	// 	return false
	// }
	// err1 := relationDao.PlusFollowerCount(id2)
	// if err1 != nil {
	// 	log.Println("关注逻辑时增加粉丝数发生了错误：", err1)
	// 	return false
	// }
	// 更新数据库
	err2 := relationDao.CreateRelationInfo(id1, id2)
	if err2 != nil {
		log.Println("关注逻辑时新增关系记录发生了错误：", err2)
		return false
	}
	// 删除关注缓存——id1关注数
	id1Str := strconv.FormatInt(id1, 10)
	redisDb.RdbUserInfo.HDel(redisDb.Ctx, "follow", id1Str)

	// 删除粉丝缓存——id2粉丝数
	id2Str := strconv.FormatInt(id2, 10)
	redisDb.RdbUserInfo.HDel(redisDb.Ctx, "follower", id2Str)

	addRelationToRedis(int(id1), int(id2))

	return true
}

// id1用户取关id2用户逻辑实现
func CounselFollow(id1 int64, id2 int64) bool {
	/* 实现：
	更新数据库
		0. 判断id1是否关注id2，未关注直接返回取关成功无需操作
		// 1. id1用户关注数量减少——user表更新
		// 2. id2用户粉丝数量减少——user表更新
		3. 删除一条id1和id2的关系——relation表删除
	更新缓存
	*/
	relationDao := dao.NewRelationDaoInstance()
	if ok, _ := relationDao.FindRelationBetween(id1, id2); !ok {
		return true
	}
	// err := relationDao.SubFollowCount(id1)
	// if err != nil {
	// 	log.Println("关注逻辑时减少关注数发生了错误：", err)
	// 	return false
	// }
	// err1 := relationDao.SubFollowerCount(id2)
	// if err1 != nil {
	// 	log.Println("关注逻辑时减少粉丝数发生了错误：", err1)
	// 	return false
	// }
	err2 := relationDao.DeleteRelationInfo(id1, id2)
	if err2 != nil {
		log.Println("关注逻辑时删除关系记录发生了错误：", err2)
		return false
	}
	// 删除关注缓存——id1关注数
	id1Str := strconv.FormatInt(id1, 10)
	redisDb.RdbUserInfo.HDel(redisDb.Ctx, "follow", id1Str)

	// 删除粉丝缓存——id2粉丝数
	id2Str := strconv.FormatInt(id2, 10)
	redisDb.RdbUserInfo.HDel(redisDb.Ctx, "follower", id2Str)
	return true
}

// 获取id用户的关注列表
func GetFollowList(id int64, curId int64) ([]UserData, bool) {
	/* 实现：
	1. 查询id用户的关注记录得到关注用户id——follow_list表
	2. 查询关注用户id对应的用户信息      ——user表
	3. 将isFollow改为True后存储返回
	*/
	relationDao := dao.NewRelationDaoInstance()
	userList := []UserData{}
	relations, err := relationDao.FindFollows(id)
	if err != nil {
		log.Println("获取关注列表逻辑时步骤1发生了错误: ", err)
		return userList, false
	}
	for _, rel := range relations {
		var user UserData

		user, err := new(UserImpl).GetUserDataByCurId(rel.UserId2, curId)
		if err != nil {
			log.Println("获取关注列表逻辑时步骤2发生了错误: ", err)
			return userList, false
		}
		user.IsFollow = true
		userList = append(userList, user)
	}
	return userList, true
}

func GetUserData(i int64) {
	panic("unimplemented")
}

// 获取id用户的粉丝列表
func GetFollowerList(id int64, curId int64) ([]UserData, bool) {
	/* 实现：
	1. 查询id用户的粉丝记录得到粉丝用户id——follow_list表
	2. 查询粉丝id对应的用户信息         ——user表
	3. 判断该用户是否为关注用户修改isFollow信息
	*/
	relationDao := dao.NewRelationDaoInstance()
	userList := []UserData{}
	relations, err := relationDao.FindFollowers(id)
	if err != nil {
		log.Println("获取粉丝列表逻辑时步骤1发生了错误: ", err)
		return userList, false
	}
	for _, rel := range relations {
		user, err := new(UserImpl).GetUserDataByCurId(rel.UserId1, curId)
		if err != nil {
			log.Println("获取粉丝列表逻辑时步骤2发生了错误: ", err)
			return userList, false
		}
		// 判断该用户是否有关注当前这位user粉丝
		isfollow, err := relationDao.FindRelationBetween(id, user.Id)
		if err != nil {
			log.Println("获取粉丝列表逻辑时步骤3发生了错误: ", err)
			return userList, false
		}
		user.IsFollow = isfollow
		userList = append(userList, user)
	}
	return userList, true
}

// 获取id用户的朋友列表
func GetFriendList(id int64, curId int64) ([]FriendUser, []UserData, bool) {
	/* 朋友定义：双向关注，既是id的粉丝，又是id的关注者
	实现：
	1. 查询id用户的所有关系记录得到所有与用户有关系的id (follow_list表)
	2. 如果某个关系id出现两次(set冲突判断)，则该id一定是用户的朋友，
	*/
	relationDao := dao.NewRelationDaoInstance()
	friendList := []FriendUser{}
	userList := []UserData{}
	relations, err := relationDao.FindRelations(id)
	if err != nil {
		log.Println("获取朋友列表逻辑时步骤1发生了错误: ", err)
		return friendList, userList, false
	}
	set := make(map[int64]bool)
	for _, rel := range relations {
		// 查询是否重复
		friendId := int64(0)
		_, ok1 := set[rel.UserId1]
		_, ok2 := set[rel.UserId2]
		switch {
		case !ok1 && !ok2:
			set[id^rel.UserId1^rel.UserId2] = true
		case !ok1 && ok2:
			friendId = rel.UserId2
		case ok1 && !ok2:
			friendId = rel.UserId1
		default:
			// 不可能同时为true: userid1和userid2其中一个不会存set
			log.Println("follow_detail数据库异常")
		}

		if friendId != 0 { // 找到朋友
			fmt.Println(friendId)
			delete(set, friendId)
			user, err := new(UserImpl).GetUserDataByCurId(friendId, curId)
			if err != nil {
				log.Println("获取朋友列表时根据id查找用户失败: ", err)
				return friendList, userList, false
			}
			userList = append(userList, user)
			// curMsg, curMsgType := getLastMessageInfo(id, friendId)
			// if curMsgType == ERRORTYPE {
			// 	log.Println("获取朋友列表时最新消息类型异常: ", err)
			// 	return friendList, userList, false
			// }
			var friend FriendUser
			// friend.Id = user.Id
			// friend.Name = user.Name
			// friend.FollowCount = user.FollowCount
			// friend.FollowerCount = user.FollowerCount
			// friend.IsFollow = true
			// friend.Message = curMsg
			// friend.MsgType = curMsgType
			if err != nil {
				log.Println("获取朋友列表逻辑中从user查询朋友id异常:", err)
				return friendList, userList, false
			}
			friendList = append(friendList, friend)
		}
	}
	return friendList, userList, true
}

// func GetUserDataById(userId int64, isFollow bool) (UserData, error) {
// 	var userData = UserData{}
// 	userData, _ = new(UserImpl).GetUserDataById(userId)
// 	userData.FollowCount, _ = dao.CountFollow(userId)
// 	userData.FollowerCount, _ = dao.CountFollower(userId)
// 	userData.IsFollow = isFollow
// 	return userData, nil
// }

func addRelationToRedis(userId int, targetId int) {
	// 第一次存入时，给该key添加一个-1为key，防止脏数据的写入。当然set可以去重，直接加，便于CPU。
	redisDb.RdbMessageHelper.SAdd(redisDb.Ctx, strconv.Itoa(int(userId)), -1)
	// 将查询到的关注关系注入Redis.
	redisDb.RdbMessageHelper.SAdd(redisDb.Ctx, strconv.Itoa(int(userId)), targetId)
	// 更新过期时间。
	redisDb.RdbMessageHelper.Expire(redisDb.Ctx, strconv.Itoa(int(userId)), time.Hour*48)

}
