package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hualuo321/douyin/controller"
	"github.com/hualuo321/douyin/jwt"
)

func initRouter(r *gin.Engine) {
	// gin 组
	apiRouter := r.Group("/douyin")
	// public directory is used to serve static resources
	r.StaticFS("/public", http.Dir("./public"))
	//basic apis
	apiRouter.GET("/user/", jwt.Auth(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", jwt.Auth(), controller.Publish)
	apiRouter.GET("/publish/list/", jwt.Auth(), controller.PublishList)
	apiRouter.GET("/feed/", jwt.AuthWithoutLogin(), controller.Feed)

	// extra one
	// apiRouter.POST("/relation/action/", jwt.Auth(), controller.RelationAction)
	// extra apis - II
	apiRouter.POST("/relation/action/", jwt.Auth(), controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", jwt.Auth(), controller.FollowList)
	apiRouter.GET("/relation/follower/list/", jwt.Auth(), controller.FollowerList)
	apiRouter.GET("/relation/friend/list/", jwt.Auth(), controller.FriendList)
	apiRouter.GET("/message/chat/", jwt.Auth(), controller.MessageChat)
	apiRouter.POST("/message/action/", jwt.Auth(), controller.MessageAction)
}
