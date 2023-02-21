package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hualuo321/douyin/dao"
)

// 密码加密 md5算法
func EnCoder(password string) string {
	p := []byte(password)
	m := md5.New()
	m.Write(p)
	return hex.EncodeToString(m.Sum(nil))
	//sum()对hash对象内部存储的内容进行校验，
	//追加到data的后面形成一个新的byte切片，
	//encoding/hex是一个将byte切片转换为字符串的编码工具库
}

// 根据username生成一个token
func GenerateToken(username string) string {
	user := new(UserImpl).QueryUserByUsername(username)
	fmt.Printf("generate token for: %v\n", user)
	token := NewToken(user)
	println("签名令牌：", token)
	return token
}

// NewToken 实现方法
func NewToken(user dao.User) string {
	expiresTime := time.Now().Unix() + int64(60*60*24) //过期时间 一天
	fmt.Printf("expiresTime: %v\n", expiresTime)
	id := user.Id
	fmt.Printf("id: %v\n", strconv.FormatInt(id, 10))
	claims := jwt.StandardClaims{ //包含的字段 定义需求，也就是通过jwt传输的数据
		Audience:  user.Username, //签发人
		ExpiresAt: expiresTime,   //过期时间
		Id:        strconv.FormatInt(id, 10),
		IssuedAt:  time.Now().Unix(), //签发时间
		Issuer:    "hualuo321",       //签发者
		NotBefore: time.Now().Unix(), //生效时间 时间戳 单位s
		Subject:   "token",           //主题
	}
	if token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("key")); err == nil {
		println("generate token success!\n")
		return token
	} else {
		println("generate token fail\n")
		return "fail"
	}
}
