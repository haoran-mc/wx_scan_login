package route

import (
	"github.com/gin-gonic/gin"
	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/app"
)

func Register(engine *gin.Engine) {
	engine.GET("/qrcode", app.Auth.GetQRCode)         // web 获取二维码
	engine.PUT("/scan", app.Auth.Scan)                // 微信扫码
	engine.GET("/check_status", app.Auth.CheckStatus) // web 查看后端状态
	engine.POST("/login", app.Auth.Login)             // 小程序确认登录

	engine.Use(visitor)

	engine.GET("/welcome", app.Home.Welcome) // 进入欢迎页面
	engine.GET("/logout", app.Auth.Logout)   // 登出清除缓存
}
