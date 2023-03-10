package service

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/model"
	"github.com/haoran-mc/wx_scan_login/web-back-end/pkg/db"
	login_session "github.com/haoran-mc/wx_scan_login/web-back-end/pkg/sessions"
)

const (
	// session 中存的东西
	userKey = "user" // model.Users，用户信息
)

type BaseContext struct {
	Ctx     *gin.Context
	Session *sessions.Session
}

func Context(ctx *gin.Context) *BaseContext {
	stx := &BaseContext{
		Ctx:     ctx,
		Session: login_session.GetSession(ctx.Request),
	}
	return stx
}

// SetAuth 设置授权
func (c *BaseContext) SetAuth(users model.Users) {
	s, _ := json.Marshal(users)
	c.Session.Values[userKey] = string(s)
	_ = c.Session.Save(c.Ctx.Request, c.Ctx.Writer)
}

// Auth 获取授权
func (c *BaseContext) Auth() *model.Users {
	var user *model.Users
	str := c.Session.Values[userKey]
	if str == nil {
		return user
	}
	if v, ok := str.(string); ok {
		_ = json.Unmarshal([]byte(v), &user)
	}
	return user
}

// Refresh 刷新授权
func (c *BaseContext) Refresh() {
	var user model.Users
	db.DB().Model(&model.Users{}).Where("id", c.Auth().ID).Find(&user)
	c.SetAuth(user)
}

// Check 检查授权
func (c *BaseContext) Check() bool {
	user := c.Auth()
	if user == nil {
		return false
	} else {
		return true
	}
}

// Forget 清除授权
func (c *BaseContext) Forget() {
	delete(c.Session.Values, userKey)
	_ = c.Session.Save(c.Ctx.Request, c.Ctx.Writer)
}
