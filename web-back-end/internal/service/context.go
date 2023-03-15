package service

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/consts"
	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/models"
	"github.com/haoran-mc/wx_scan_login/web-back-end/pkg/db"
	login_session "github.com/haoran-mc/wx_scan_login/web-back-end/pkg/sessions"
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
func (c *BaseContext) SetAuth(users models.Users) {
	s, _ := json.Marshal(users)
	c.Session.Values[consts.SessionKeyUser] = string(s)
	_ = c.Session.Save(c.Ctx.Request, c.Ctx.Writer)
}

// Auth 获取授权
func (c *BaseContext) Auth() *models.Users {
	var user *models.Users
	str := c.Session.Values[consts.SessionKeyUser]
	if str == nil {
		return user
	}
	if v, ok := str.(string); ok {
		_ = json.Unmarshal([]byte(v), &user)
	}
	return user
}

// RefreshAuth 刷新授权
func (c *BaseContext) RefreshAuth() {
	var user models.Users
	db.DB().Model(&models.Users{}).Where("id", c.Auth().ID).Find(&user)
	c.SetAuth(user)
}

// CheckAuth 检查授权
func (c *BaseContext) CheckAuth() bool {
	user := c.Auth()
	if user == nil {
		return false
	} else {
		return true
	}
}

// ForgetAuth 清除授权
func (c *BaseContext) ForgetAuth() {
	delete(c.Session.Values, consts.SessionKeyUser)
	_ = c.Session.Save(c.Ctx.Request, c.Ctx.Writer)
}
