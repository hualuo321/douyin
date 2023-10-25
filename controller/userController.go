package controller

import (
	"github.com/hualuo321/douyin/dao"
	"github.com/hualuo321/douyin/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
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
	User service.User `json:"user"`
}

// 注册 douyin/user/register/
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	usi := service.UserServiceImpl{}
	u := usi.GetTableUserByUsername(username)

	user := useri.QueryUserByUsername(username)
	if username == u.Name {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		newUser := dao.TableUser{
			Name:     username,
			Password: service.EnCoder(password),
		}
		if usi.InsertTableUser(&newUser) != true {
			println("Insert Data Fail")
		}
		u := usi.GetTableUserByUsername(username)
		token := service.GenerateToken(username)
		log.Println("注册返回的id: ", u.Id)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   u.Id,
			Token:    token,
		})
	}
}

// 登录 douyin/usr/login
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	encoderPassword := service.EnCoder(password)
	println(encoderPassword)
	encoderPassword := service.EnCoder(password)
	usi := service.UserServiceImpl{}
	u := usi.GetTableUserByUsername(username)
	if encoderPassword == u.Password {
		token := service.GenerateToken(username)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   u.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Username or Password Error"},
		})
	}
}

// UserInfo GET douyin/user/ 用户信息
func UserInfo(c *gin.Context) {
	user_id := c.Query("user_id")
	id, _ := strconv.ParseInt(user_id, 10, 64)
	usi := service.UserServiceImpl{
		FollowService: &service.FollowServiceImp{},
		LikeService:   &service.LikeServiceImpl{},
	}
	if u, err := usi.GetUserById(id); err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User Doesn't Exist"},
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     u,
		})
	}
}
