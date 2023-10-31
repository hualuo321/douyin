package service

// LikeService: 定义点赞状态和点赞数量
type LikeService interface {
	/* 1. 其他模块 (video) 需要使用的业务方法 */
	// IsFavorite: 根据当前 videoId 判断是否点赞了该视频
	IsFavourite(videoId int64, userId int64) (bool, error)

	// FavouriteCount 根据当前 videoId 获取当前视频点赞数量
	FavouriteCount(videoId int64) (int64, error)

	// TotalFavourite 根据 userId 获取这个用户总共被点赞数量
	TotalFavourite(userId int64) (int64, error)

	// FavouriteVideoCount 根据 userId 获取这个用户点赞视频数量
	FavouriteVideoCount(userId int64) (int64, error)

	/* 2. 直接 request 需要实现的功能 */
	// FavouriteAction: 当前操作行为，点赞, 或取消点赞。
	FavouriteAction(userId int64, videoId int64, actionType int32) error

	// GetFavouriteList: 获取当前用户的所有点赞视频，调用 videoService 的方法
	GetFavouriteList(userId int64, curId int64) ([]Video, error)
}