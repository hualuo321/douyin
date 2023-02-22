package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hualuo321/douyin/controller"
	"github.com/hualuo321/douyin/jwt"
)

func initRouter(r *gin.Engine) {
	// gin 组
	apiRouter := r.Group("/douyin")
	// public directory is used to serve static resources
	r.Static("/static", "./public")
	//basic apis
	apiRouter.GET("/user/", jwt.Auth(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", jwt.Auth(), controller.Publish)
	apiRouter.GET("/publish/list/", jwt.Auth(), controller.PublishList)
	apiRouter.GET("/feed/", jwt.AuthWithoutLogin(), controller.Feed)
}
