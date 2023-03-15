package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/consts"
	"github.com/haoran-mc/wx_scan_login/web-back-end/pkg/sessions"
)

type store struct{}

var Store store

func (s *store) GetSessionInfo(ctx *gin.Context) {
	session := sessions.GetSession(ctx.Request)
	status := session.Values[consts.SessionKeyStatus]
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    status,
	})
}
