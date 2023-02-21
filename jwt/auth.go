package jwt

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int32  `json:"status_code"` //1为错误 0为正确
	StatusMsg  string `json:"status_msg,omitempty"`
}

// 放在路由访问
// 如果用户的token正确，解析token，将userId放入context中，否则报错
func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		auth := context.Query("token")
		if len(auth) == 0 {
			context.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "未授权!",
			})
		}
		//解析token
		token, err := parseToken(auth)
		if err != nil {
			context.Abort()
			context.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Token 错误！",
			})
		} else {
			println("token 正确！")
		}
		context.Set("userId", token.Id)
		context.Next()
	}
}

// 未登录情况下
// 若携带token则解析user_id放入context，未携带则放入用户id 默认值0
func AuthWithoutLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Query("token")
		var userId string
		if len(auth) == 0 {
			userId = "0"
		} else {
			token, err := parseToken(auth)
			if err != nil {
				c.Abort() //阻止挂起的函数
				c.JSON(http.StatusUnauthorized, Response{
					StatusCode: -1,
					StatusMsg:  "Token Error",
				})
			} else {
				userId = token.Id
				println("token 正确")
			}
			c.Set("userId", userId)
			c.Next() //执行挂起的函数
		}
	}
}

func parseToken(token string) (*jwt.StandardClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte("key"), nil
	})

	if err == nil && jwtToken != nil {
		if claim, ok := jwtToken.Claims.(*jwt.StandardClaims); ok && jwtToken.Valid {
			return claim, nil
		}
	}
	return nil, err
}
