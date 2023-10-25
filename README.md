<a name="readme-top"></a>

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/hualuo321/douyin">
    <img src="images/error404.png" alt="Logo" width="180" height="180">
  </a>
</div>

# Easy Douyin
一个极简版的抖音服务器，包含登录，注册，发布视频，首页刷新视频，关注取关用户等功能。

# Dependency
**开发前的配置要求**：
- go 1.18.1
- MySQL
- 搭建Redis、RabbitMQ环境
- 配置静态资源服务器：安装Nginx、vsftpd、ffmpeg
- [最新版抖音客户端软件](https://bytedance.feishu.cn/docs/doccnM9KkBAdyDhg8qaeGlIz7S7)

`go build && ./easy-douyin`
# Directory
```
├── controller                          # API functions and response structure
│   ├── comment.go
│   ├── common.go
│   ├── favorite.go
│   ├── feed.go
│   ├── publish.go
│   ├── relation.go
│   └── user.go
├── public                              # public static resources on server
│   ├── covers                          # store static pictures of video cover
│   │   └──*.png
│   └── videos                          # store static videos
│       └──*.mp4
├── repository                          # init, models and CRUD of database
│   ├── comment.go
│   ├── db_init.go
│   ├── follow.go
│   ├── like.go
│   ├── user.go
│   └── video.go
├── service                             # realisation of functions in controller
│   ├── comment.go
│   ├── follow.go
│   ├── like.go
│   ├── user.go
│   └── video.go
├── test                                # test files
│   └── ...
├── util
│   └── MD5.go                          # encryption function
├── middleware
│   ├── logger.go                       # provide error log
│   ├── auth.go                         # authority function handler
│   └── jwt.go                          # generate and parse jwt
├── .gitattributes
├── .gitignore
├── go.mod
├── go.sum
├── main.go                             # start of execution
└── router.go                           # path configuration
```

