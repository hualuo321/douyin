package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hualuo321/douyin/controller"
	"github.com/hualuo321/douyin/middleware/jwt"
)

func initRouter(r *gin.Engine) {
	// gin 组
	apiRouter := r.Group("/douyin")
	// basic apis
	apiRouter.GET("/feed/", jwt.AuthWithoutLogin(), controller.Feed)
	apiRouter.GET("/user/", jwt.Auth(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", jwt.Auth(), controller.Publish)
	apiRouter.GET("/publish/list/", jwt.Auth(), controller.PublishList)

	// extra apis - II
	apiRouter.POST("/relation/action/", jwt.Auth(), controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", jwt.Auth(), controller.FollowList)
	apiRouter.GET("/relation/follower/list/", jwt.Auth(), controller.FollowerList)
	apiRouter.GET("/relation/friend/list/", jwt.Auth(), controller.FriendList)
	apiRouter.GET("/message/chat/", jwt.Auth(), controller.MessageChat)
	apiRouter.POST("/message/action/", jwt.Auth(), controller.MessageAction)
}
