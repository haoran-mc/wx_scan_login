package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/service"
)

// visitor 访问者
func visitor(ctx *gin.Context) {
	s := service.Context(ctx)
	if !s.Check() {
		s.Ctx.Redirect(http.StatusFound, "/qrcode")
		ctx.Abort()
	} else {
		ctx.Next()
	}
}
