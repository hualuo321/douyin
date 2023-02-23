package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hualuo321/douyin/dao"
	"github.com/hualuo321/douyin/service"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User service.UserData `json:"user"`
}

// 注册 douyin/user/register/
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	useri := service.UserImpl{}
	user := useri.QueryUserByUsername(username)
	if username == user.Username {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户已存在！"},
		})
	} else {
		newUser := dao.User{
			Username: username,
			Password: service.EnCoder(password),
		}
		if !useri.InsertUser(&newUser) {
			log.Println("新增用户失败！")
		}
		user = useri.QueryUserByUsername(username)
		token := service.GenerateToken(username) //产生token
		log.Println("注册返回的id:", user.Id)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id, //user_id
			Token:    token,
		})
	}
}

// 登录 douyin/usr/login
func Login(c *gin.Context) {
	fmt.Println("---当前位于contorller/user/Login() 1 ---")
	username := c.Query("username")
	fmt.Println("-当前输入的用户是：", username)
	password := c.Query("password")
	fmt.Println("-当前输入的密码是：", password)
	encoderPassword := service.EnCoder(password)
	useri := service.UserImpl{}
	user := useri.QueryUserByUsername(username)
	fmt.Println("---当前位于contorller/user/Login() 2 ---")
	fmt.Println("-查询到的用户信息为", user)
	// 登录成功，则返回 用户编号 与 token
	if encoderPassword == user.Password {
		token := service.GenerateToken(username)
		fmt.Println("-登录的用户信息为", user)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户名或密码错误！"},
		})
	}
}

// userInfo douyin/user/
func UserInfo(c *gin.Context) {
	id := c.Query("user_id")
	print("------查询到的 userid 为", id)
	user_id, _ := strconv.ParseInt(id, 10, 64)
	print("------查询到的 userid 为", user_id)
	useri := service.UserImpl{}
	//非登录情况下
	user, err := useri.GetUserDataById(user_id)
	fmt.Println("--------查询到的用户信息为", user)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户不存在！"},
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	}

}
