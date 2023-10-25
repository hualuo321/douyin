package config

import "time"

// Secret 密钥
var Secret = "douyin"

// ftp 服务器连接配置信息
const ConConfig = "192.168.31.117:21"
const FtpUser = "ftpuser"
const FtpPsw = "123456"
const HeartbeatTime = 2 * 60

// SSH 配置
const HostSSH = "192.168.31.117"
const UserSSH = "ftpuser"
const PasswordSSH = "123456"
const TypeSSH = "password"
const PortSSH = 22
const MaxMsgCount = 100
const SSHHeartbeatTime = 10 * 60

// PlayUrlPrefix 存储的图片和视频的链接
const PlayUrlPrefix = "http://43.138.25.60/"
const CoverUrlPrefix = "http://43.138.25.60/images/"

// VideoCount 每次获取视频流的数量
const VideoCount = 5

// OneDayOfHours 时间
var OneDayOfHours = 60 * 60 * 24
var OneMinute = 60 * 1
var OneMonth = 60 * 60 * 24 * 30
var OneYear = 365 * 60 * 60 * 24
var ExpireTime = time.Hour * 48 // 设置Redis数据热度消散时间。