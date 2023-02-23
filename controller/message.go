package controller

import (
	"net/http"
	"strconv"

	"github.com/hualuo321/douyin/dao"
	"github.com/hualuo321/douyin/service"

	"github.com/gin-gonic/gin"
)

type ChatResponse struct {
	Response
	MessageList []dao.MessageData `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
	/* 给发送的消息存到tempChat里面
	   给发送的消息存到数据库里面，*/
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	toUserId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	content := c.Query("content")

	key := genChatKey(userId, toUserId)
	ok := service.ReceiveMessage(userId, toUserId, content, key)
	// 响应
	if ok {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "发送失败"})
	}
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	toUserId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	//messages := []dao.MessageData{}
	key := genChatKey(toUserId, userId)
	messages, ok := service.ResponseMessage(key)
	println(messages)
	if ok {
		//fmt.Println(user)
		c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0}, MessageList: messages})
	} else {
		c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 1, StatusMsg: "接收失败"}})
	}
}

//	func genChatKey(userIdA int64, userIdB int64) string {
//		if userIdA > userIdB {
//			return fmt.Sprintf("%d_%d", userIdB, userIdA)
//		}
//		return fmt.Sprintf("%d_%d", userIdA, userIdB)
//	}
func genChatKey(userId int64, toUserId int64) string {
	strSendId := strconv.FormatInt(userId, 10)
	strRecId := strconv.FormatInt(toUserId, 10)
	key := strSendId + " " + strRecId

	return key
}
