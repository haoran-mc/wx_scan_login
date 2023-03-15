package models

import (
	"time"
)

type Users struct {
	Model
	Name        string     // 用户名、昵称，首次注册使用微信昵称
	Phone       string     // 手机号
	Gender      uint8      // 性别
	Email       string     // 邮箱
	Avatar      string     // 头像
	IsAdmin     uint8      // 是否为管理员
	LastLoginAt *time.Time // 上一次登录的时间
}
