package controller

import (
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
	User service.UserData `json:"user_data"`
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
	username := c.Query("username")
	password := c.Query("password")
	encoderPassword := service.EnCoder(password)
	useri := service.UserImpl{}
	user := useri.QueryUserByUsername(username)
	log.Println("正在登录的用户是：", user)

	// 登录成功，则返回 用户编号 与 token
	if encoderPassword == user.Password {
		token := service.GenerateToken(username)
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
	user_id, _ := strconv.ParseInt(id, 10, 64)
	useri := service.UserImpl{
		//RelationService: service.RelationServiceImpl{},
		//FavoriteService: service.FavoriteServiceImpl{},
	}
	//非登录情况下
	user, err := useri.GetUserById(user_id)
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
