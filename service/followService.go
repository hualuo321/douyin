package service

// FollowService: 定义用户关系接口以及用户关系中的各种方法
type FollowService interface {
	/* 1. 其他模块需要调用的业务方法 */
	// IsFollowing: 根据当前 user_id 和目标 user_id 来判断当前用户是否关注了目标用户
	IsFollowing(userId int64, targetId int64) (bool, error)

	// GetFollowerCnt: 根据 user_id 来查询该用户被多少其他用户关注
	GetFollowerCnt(userId int64) (int64, error)

	// GetFollowingCnt: 根据 user_id 来查询该用户关注了多少其它用户
	GetFollowingCnt(userId int64) (int64, error)

	/* 2. 直接 request 需要的业务方法 */
	// AddFollowRelation: 当前用户关注目标用户
	AddFollowRelation(userId int64, targetId int64) (bool, error)

	// DeleteFollowRelation: 当前用户取消对目标用户的关注
	DeleteFollowRelation(userId int64, targetId int64) (bool, error)

	// GetFollowing: 根据 user_id 获取当前用户的关注列表
	GetFollowing(userId int64) ([]User, error)

	// GetFollowers: 根据 user_id 获取当前用户的粉丝列表
	GetFollowers(userId int64) ([]User, error)
}