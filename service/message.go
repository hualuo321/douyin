package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hualuo321/douyin/dao"
	"github.com/hualuo321/douyin/redisDb"

	"github.com/go-redis/redis/v8"
)

var (
	RECEIVETYPE int8 = 0
	SENDTYPE    int8 = 1
	ERRORTYPE   int8 = 2
)

// 服务器接收到发送端发送消息，存入数据表
func ReceiveMessage(userId int64, toUserId int64, content string, key string) bool {
	/*
		写入message表需要返回id就完事了，给用户响应
		写入recMessage
	*/
	// 写入mysql
	messageDao := dao.NewMessageDaoInstance()
	messageId, err := messageDao.CreateMessage(userId, toUserId, content)
	if err != nil {
		log.Println("接收消息逻辑时创建消息失败: ", err)
		return false
	}
	// 写入redis
	strMessageId := strconv.FormatInt(messageId, 10) + " "
	redisDb.RdbMessageHelper.Append(redisDb.Ctx, key, strMessageId)
	return true
}

// 服务器接到接收端请求接收消息
func ResponseMessage(key string) ([]dao.MessageData, bool) {
	/*
		说明：接收端会5s一次轮询，只能给它指定它接收的消息
		1. 从redis中读取key("touserid-user")
		2. 读完后删除该键
		3. 若值为空，则不存在消息
		4. 若值不为空，将值解析为一个个message_id，int64
			解析："id1 id2 id3 "--空格切割-->["id1" "id2" "id3"]--遍历类型转换-->id1, id2, id3
		5. 查找message表找到message_id对应的消息存入list
	*/
	messageList := []dao.MessageData{}
	// 1.
	value, err := redisDb.RdbMessageHelper.Get(redisDb.Ctx, key).Result() // redis提供方法将ridis.stringcmd格式转化string
	if err != nil && (!errors.Is(err, redis.Nil)) {
		fmt.Println("获取失败")
		return messageList, false
	}
	// 2.
	redisDb.RdbMessageHelper.Del(redisDb.Ctx, key)
	// 3.
	if len(value) == 0 {
		return messageList, true
	}
	// 4.
	ids := strings.Fields(value)
	messageDao := dao.NewMessageDaoInstance()
	for _, eachId := range ids {
		id, _ := strconv.ParseInt(eachId, 10, 64)
		// 5.
		message, err := messageDao.FindMessageById(id)
		message.CreatTime = ""
		if err != nil {
			log.Println("接收消息逻辑中查找消息失败: ", err)
			return messageList, false
		}
		messageList = append(messageList, message)
		fmt.Println(messageList)
	}
	return messageList, true
}

func getLastMessageInfo(userId1 int64, userId2 int64) (string, int8) {
	messageDao := dao.NewMessageDaoInstance()
	message1, err1 := messageDao.FindLastMessage(userId1, userId2)
	if err1 != nil {
		return "", ERRORTYPE
	}
	message2, err2 := messageDao.FindLastMessage(userId2, userId1)
	if err2 != nil {
		return "", ERRORTYPE
	}
	if message1.Id > message2.Id {
		return message1.Content, SENDTYPE
	}
	if message2.Id > message1.Id {
		return message2.Content, RECEIVETYPE
	}
	return "", ERRORTYPE
}
