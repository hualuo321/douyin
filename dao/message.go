package dao

import (
	"sync"
	"time"
)

type MessageData struct {
	Id         int64  `json:"id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	FromUserId int64  `json:"from_user_id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreatTime  string `json:"create_time,omitempty"`
}

func (p MessageData) TableName() string {
	return "message"
}

// FollowDao 把dao层看成整体，把dao的curd封装在一个结构体中。
type MessageDao struct {
}

var (
	messageDao  *MessageDao //操作该dao层crud的结构体变量。
	messageOnce sync.Once   //单例限定，去限定申请一个followDao结构体变量。
)

// NewFollowDaoInstance 生成并返回RelationDao的单例对象。
func NewMessageDaoInstance() *MessageDao {
	relationOnce.Do(
		func() {
			messageDao = &MessageDao{}
		})
	return messageDao
}

// var messageIdSequence = int64(1)

/* CRUD实现 */
func (*MessageDao) CreateMessage(fromUserId int64, toUserId int64, content string) (int64, error) {
	curMessage := MessageData{
		//		Id:         messageIdSequence,
		ToUserId:   toUserId,
		FromUserId: fromUserId,
		Content:    content,
		CreatTime:  time.Now().Format("2006-01-02 15:04:05"),
	}
	err := db.Create(&curMessage).Error
	return curMessage.Id, err
}

func (*MessageDao) FindMessageById(id int64) (MessageData, error) {
	var message MessageData
	err := db.Find(&message, id).Error
	return message, err
}

func (*MessageDao) FindLastMessage(userId int64, toUserId int64) (MessageData, error) {
	var message MessageData
	err := db.Where("from_user_id=? AND to_user_id = ?", userId, toUserId).Last(&message).Error
	return message, err
}

func (*MessageDao) FindLast10Message() ([]MessageData, error) {
	var messages []MessageData
	limit := 10

	err := db.Limit(limit).Find(&messages).Order("id desc").Error

	return messages, err
}
