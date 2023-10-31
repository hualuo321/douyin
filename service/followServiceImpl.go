package service

import (
	"github.com/hualuo321/douyin/config"
	"github.com/hualuo321/douyin/dao"
	"github.com/hualuo321/douyin/middleware/rabbitmq"
	"github.com/hualuo321/douyin/middleware/redis"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

// FollowServiceImp 该结构体继承 FollowService 接口。
type FollowServiceImp struct {
	UserService
}