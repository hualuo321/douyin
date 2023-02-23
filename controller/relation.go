package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/hualuo321/douyin/service"

	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []service.UserData `json:"user_list"`
}
type FriendListResponse struct {
	Response
	UserList []service.FriendUser `json:"friend_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	userId, err1 := strconv.ParseInt(c.GetString("userId"), 10, 64)
	toUserId, err2 := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	actionType := c.Query("action_type")

	if nil != err1 || nil != err2 {
		fmt.Printf("fail")
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户id格式错误"})
		return
	}
	ok := false

	switch actionType { // 实现关注&取消关注逻辑
	case "1": // 关注 u1关注u2
		fmt.Println("userId", userId, "toUserId", toUserId)
		ok = service.DoFollow(userId, toUserId)
	case "2": // 取消关注
		ok = service.CounselFollow(userId, toUserId)
	}
	// 响应
	if ok {
		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "关注成功"})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "关注失败"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	curId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)

	userList, ok := service.GetFollowList(userId, curId)
	// 响应
	if ok {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{StatusCode: 0},
			UserList: userList})
	} else {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "获取关注列表失败"},
			UserList: userList})
	}

}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	curId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	userList, ok := service.GetFollowerList(userId, curId) /////?????
	// 响应
	if ok {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{StatusCode: 0},
			UserList: userList,
		})
	} else {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "获取粉丝列表失败"},
			UserList: userList,
		})
	}

}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	curId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	_, userList, ok := service.GetFriendList(userId, curId)
	fmt.Println(userList)
	// 响应
	if ok {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{StatusCode: 0},
			UserList: userList,
		})
	} else {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "获取朋友列表失败"},
			UserList: userList,
		})
	}
}
