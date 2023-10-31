package service

import "github.com/hualuo321/douyin/dao"

// 接口里列出了很多方法
type UserService interface {
	// GetTableUserList: 获得全部 TableUser 对象
	GetTableUserList() []dao.TableUser

	// GetTableUserByUsername: 根据 username 获得 TableUser 对象
	GetTableUserByUsername(name string) dao.TableUser

	// GetTableUserById: 根据 user_id 获得 TableUser 对象
	GetTableUserById(id int64) dao.TableUser

	// InsertTableUser: 将 tableUser 插入表内
	InsertTableUser(tableUser *dao.TableUser) bool

	// GetUserById: 未登录情况下根据 user_id 获得 User 对象
	GetUserById(id int64) (User, error)

	// GetUserByIdWithCurId: 已登录 (curID) 情况下,根据 user_id 获得 User 对象
	GetUserByIdWithCurId(id int64, curId int64) (User, error)

	// 根据 token 返回 id
	// 接口: auth 中间件, 解析完 token, 将 userid 放入 context
	// 调用方法: 直接在 context 内拿参数 "userId" 的值
}

// User: 最终封装后, controller 返回的 User 结构体
type User struct {
	Id             int64  `json:"id,omitempty"`					// ID
	Name           string `json:"name,omitempty"`				// 姓名
	FollowCount    int64  `json:"follow_count"`					// 关注了多少用户
	FollowerCount  int64  `json:"follower_count"`				// 被多少用户关注	
	IsFollow       bool   `json:"is_follow"`					// 自己是否关注该用户
	TotalFavorited int64  `json:"total_favorited,omitempty"`	// 自己作品被多少人点赞
	FavoriteCount  int64  `json:"favorite_count,omitempty"`		// 自己点赞了多少其他作品
}