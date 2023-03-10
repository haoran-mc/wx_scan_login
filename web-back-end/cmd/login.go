package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/route"
	"github.com/haoran-mc/wx_scan_login/web-back-end/pkg/config"
	_ "github.com/haoran-mc/wx_scan_login/web-back-end/pkg/config"
	_ "github.com/haoran-mc/wx_scan_login/web-back-end/pkg/db"
	"github.com/haoran-mc/wx_scan_login/web-back-end/pkg/logger"
	_ "github.com/haoran-mc/wx_scan_login/web-back-end/pkg/logger"
	_ "github.com/haoran-mc/wx_scan_login/web-back-end/pkg/sessions"
)

func main() {
	engine := gin.Default()
	route.Register(engine)

	if err := engine.Run(":" + config.Conf.Port); err != nil {
		logger.Logger.Fatal("server running error: ", zap.Error(err))
	}
}
