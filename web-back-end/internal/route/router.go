package route

import (
	"github.com/gin-gonic/gin"
	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/app"
)

func Register(engine *gin.Engine) {
	// web 获取二维码
	engine.GET("/qrcode", app.Auth.GetQRCode)

	// 已扫描，更改 web 状态，等待小程序确认
	engine.PUT("/change_status", app.Auth.ChangeStatus)

	// 从小程序方获取 code
	engine.POST("/login", app.Auth.Login)

	engine.Use(visitor)

	// 进入首页
	engine.GET("/welcome", app.Home.Welcome)

	// 登出清除缓存
	engine.GET("/logout", app.Auth.Logout)
}
