package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hualuo321/douyin/service"
)

// RelationActionResp 关注和取消关注需要返回结构。
type RelationActionResp struct {
	Response
}

// FollowingResp 获取关注列表需要返回的结构。
type FollowingResp struct {
	Response
	UserList []service.UserData `json:"user_list,omitempty"`
}

func RelationAction(c *gin.Context) {
	userId, err1 := strconv.ParseInt(c.GetString("userId"), 10, 64)
	toUserId, err2 := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	actionType, err3 := strconv.ParseInt(c.Query("action_type"), 10, 64)

	if nil != err1 || nil != err2 || nil != err3 || actionType < 1 || actionType > 2 {
		fmt.Printf("fail")
		c.JSON(http.StatusOK, RelationActionResp{
			Response{
				StatusCode: -1,
				StatusMsg:  "用户id格式错误",
			},
		})
		return
	}
	var err error
	if actionType == 1 {
		err = service.InsertFollow(userId, toUserId)
	} else {
		err = service.DeleteFollow(userId, toUserId)
	}
	if err == nil {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "update follow db fail"})
	}
}
