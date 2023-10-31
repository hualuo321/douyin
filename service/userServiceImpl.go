package service

import (
	"github.com/hualuo321/douyin/config"
	"github.com/hualuo321/douyin/dao"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"strconv"
	"time"
)

// UserServiceImpl: 接口组合，意味着结构体需要完全实现这两个接口的方法
type UserServiceImpl struct {
	FollowService
	LikeService
}

// GetTableUserList: 获得全部 TableUser 对象
func (usi *UserServiceImpl) GetTableUserList() []dao.TableUser {
	// 调用 dao 包中的 GetTableUserList 函数，尝试获取所有 TableUser 对象
	tableUsers, err := dao.GetTableUserList()
	if err != nil {
		log.Println("Err:", err.Error())
		return tableUsers
	}
	// 如果成功，返回获取到的 TableUser 对象列表
	log.Println("Get TableUser Success")
	return tableUsers
}

// GetTableUserByUsername: 根据 username 获得 TableUser 对象
func (usi *UserServiceImpl) GetTableUserByUsername(userName string) dao.TableUser {
	// 调用 dao 包中的 GetTableUserByUsername 函数，尝试获取 TableUser 对象
	tableUser, err := dao.GetTableUserByUsername(userName)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return tableUser
	}
	// 如果成功，返回获取到的 TableUser 对象
	log.Println("Query User Success")
	return tableUser
}

// GetTableUserById: 根据 user_id 获得 TableUser 对象
func (usi *UserServiceImpl) GetTableUserById(userId int64) dao.TableUser {
	// 调用 dao 包中的 GetTableUserById 函数，尝试获取 TableUser 对象
	tableUser, err := dao.GetTableUserById(userId)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return tableUser
	}
	// 如果成功，返回获取到的 TableUser 对象
	log.Println("Query User Success")
	return tableUser
}

// InsertTableUser: 将 tableUser 插入表内
func (usi *UserServiceImpl) InsertTableUser(tableUser *dao.TableUser) bool {
	// 调用 dao 包中的 InsertTableUser 函数，尝试插入数据到表中
	flag := dao.InsertTableUser(tableUser)
	if flag == false {
		log.Println("插入失败")
		return false
	}
	// 如果成功，返回 true
	return true
}

// GetUserById 未登录情况下,根据user_id获得User对象
func (usi *UserServiceImpl) GetUserById(userId int64) (User, error) {
	// 创建一个初始的 User 结构体，该结构体包含默认值
	user := User{
		Id:             0,
		Name:           "",
		FollowCount:    0,
		FollowerCount:  0,
		IsFollow:       false,
	}
	// 调用 dao 包中的 GetTableUserById 函数，尝试根据用户ID获取 TableUser 对象
	tableUser, err := dao.GetTableUserById(userId)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return user, err
	}
	log.Println("Query User Success")
	// 调用 usi 的 GetFollowingCnt 方法，获取用户的关注数
	followCount, _ := usi.GetFollowingCnt(id)
	if err != nil {
		log.Println("Err:", err.Error())
	}
	// 调用 usi 的 GetFollowerCnt 方法，获取用户的粉丝数
	followerCount, _ := usi.GetFollowerCnt(id)
	if err != nil {
		log.Println("Err:", err.Error())
	}
	// 将获取到的数据填充到结构体中
	user = User{
		Id:             userId,
		Name:           tableUser.Name,
		FollowCount:    followCount,
		FollowerCount:  followerCount,
		IsFollow:       false,
	}
	// 返回填充后的 User 结构体
	return user, nil
}

// GetUserByIdWithCurId: 已登录 (curID) 情况下, 根据 user_id 获得 User 对象
func (usi *UserServiceImpl) GetUserByIdWithCurId(userId int64, curId int64) (User, error) {
	// 创建一个初始的 User 结构体，该结构体包含默认值
	user := User{
		Id:             0,
		Name:           "",
		FollowCount:    0,
		FollowerCount:  0,
		IsFollow:       false,
		TotalFavorited: 0,
		FavoriteCount:  0,
	}
	// 调用 dao 包中的 GetTableUserById 函数，尝试根据用户 ID 获取 TableUser 对象
	tableUser, err := dao.GetTableUserById(userId)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return user, err
	}
	log.Println("Query User Success")
	// 调用 usi 的 GetFollowingCnt 方法，获取用户的关注数
	followCount, err := usi.GetFollowingCnt(userId)
	if err != nil {
		log.Println("Err:", err.Error())
	}
	// 调用 usi 的 GetFollowerCnt 方法，获取用户的粉丝数
	followerCount, err := usi.GetFollowerCnt(userId)
	if err != nil {
		log.Println("Err:", err.Error())
	}
	// 调用 usi 的 IsFollowing 方法，检查当前用户是否关注了目标用户
	isfollow, err := usi.IsFollowing(curId, userId)
	if err != nil {
		log.Println("Err:", err.Error())
	}
	// 获取 LikeService 实例，解决循环依赖
	u := GetLikeService()
	// 调用 u 的 TotalFavourite 方法，获取用户的总点赞数
	totalFavorited, _ := u.TotalFavourite(userId)
	// 调用 u 的 FavouriteVideoCount 方法，获取用户点赞的视频数量
	favoritedCount, _ := u.FavouriteVideoCount(userId)
	// 将获取到的数据填充到结构体中，并返回该结构体
	user = User{
		Id:             id,
		Name:           tableUser.Name,
		FollowCount:    followCount,
		FollowerCount:  followerCount,
		IsFollow:       isfollow,
		TotalFavorited: totalFavorited,
		FavoriteCount:  favoritedCount,
	}
	// 返回填充后的 User 结构体
	return user, nil
}

// GenerateToken: 根据 username 生成一个 token
func GenerateToken(username string) string {
	// 使用 UserService 的 GetTableUserByUsername 方法，通过用户名获取用户信息
	u := UserService.GetTableUserByUsername(new(UserServiceImpl), username)
	fmt.Printf("generatetoken: %v\n", u)
	// 调用 NewToken 函数，根据用户信息创建一个 token
	token := NewToken(u)
	println(token)
	// 返回生成的 token
	return token
}

// NewToken 根据用户信息创建 token 的方法
func NewToken(u dao.TableUser) string {
	// 计算 token 过期时间，为当前时间加上一天的秒数
	expiresTime := time.Now().Unix() + int64(config.OneDayOfHours)
	fmt.Printf("expiresTime: %v\n", expiresTime)
	id64 := u.Id
	fmt.Printf("id: %v\n", strconv.FormatInt(id64, 10))
	// 创建 JWT 标准声明（claims）
	claims := jwt.StandardClaims{
		Audience:  u.Name,
		ExpiresAt: expiresTime,
		Id:        strconv.FormatInt(id64, 10),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "tiktok",
		NotBefore: time.Now().Unix(),
		Subject:   "token",
	}
	// 定义 JWT 密钥
	var jwtSecret = []byte(config.Secret)
	// 创建基于声明的 token
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 签署 token
	if token, err := tokenClaims.SignedString(jwtSecret); err == nil {
		token = "Bearer " + token
		println("generate token success!\n")
		return token
	} else {
		println("generate token fail\n")
		return "fail"
	}
}

// EnCoder 密码加密
func EnCoder(password string) string {
	h := hmac.New(sha256.New, []byte(password))
	// 计算哈希值并将其转为十六进制格式
	sha := hex.EncodeToString(h.Sum(nil))
	fmt.Println("Result: " + sha)
	// 返回加密后的密码
	return sha
}