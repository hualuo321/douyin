package main

import "github.com/gin-gonic/gin"

func initRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")
	//basic apis
	apiRouter.POST("/user/")
}
