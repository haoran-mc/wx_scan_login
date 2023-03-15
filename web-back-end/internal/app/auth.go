package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/consts"
	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/entity"
	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/service"
	"github.com/haoran-mc/wx_scan_login/web-back-end/pkg/config"
	"github.com/haoran-mc/wx_scan_login/web-back-end/pkg/logger"
	"github.com/haoran-mc/wx_scan_login/web-back-end/pkg/sessions"
	"go.uber.org/zap"
)

type auth struct{}

var Auth = auth{}

func (c *auth) GetQRCode(ctx *gin.Context) {
	// 1. 生成二维码
	QRCodeUrl, err := service.AuthService(ctx).GenerateQRCode()
	if err != nil {
		logger.Logger.Error("failed to generate qrcode", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to generate qrcode",
		})
		return
	}

	// 2. 修改 session 中的 x-dl-status: 1，同时随响应返回 session_id
	if err := service.AuthService(ctx).
		ChangeStatus(consts.StatusUnscanned); err != nil {
		logger.Logger.Error("failed to modify x-dl-status", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to change status",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    QRCodeUrl,
	})
}

// 微信扫描二维码，小程序发送请求修改后端 session 中的 x-dl-status
func (c *auth) Scan(ctx *gin.Context) {
	if err := service.AuthService(ctx).
		ChangeStatus(consts.StatusScanned); err != nil {
		logger.Logger.Error("failed to modify x-dl-status", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to change status",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// 浏览器不断发送请求，查看后端 session 中的状态
func (c *auth) CheckStatus(ctx *gin.Context) {
	session := sessions.GetSession(ctx.Request)
	backStatus := session.Values[consts.SessionKeyStatus] // session 中的状态
	frontStatus, _ := ctx.Cookie(consts.SessionKeyStatus) // 浏览器的状态
	if frontStatus == consts.StatusUnscanned &&
		backStatus == consts.StatusScanned {
		err := service.AuthService(ctx).ChangeStatus(consts.StatusScanned)
		if err != nil {
			ctx.JSON(http.StatusNotModified, gin.H{
				"message": "failed to change status",
			})
		}
		ctx.Redirect(http.StatusTemporaryRedirect, config.Conf.Web.Router.Scanned)
	} else if frontStatus == consts.StatusScanned &&
		backStatus == consts.StatusLogined {
		err := service.AuthService(ctx).ChangeStatus(consts.StatusLogined)
		if err != nil {
			ctx.JSON(http.StatusNotModified, gin.H{
				"message": "failed to change status",
			})
		}
		ctx.Redirect(http.StatusTemporaryRedirect, config.Conf.Web.Router.Welcome)
	} else { // 否则不应该修改状态
		ctx.JSON(http.StatusNotModified, gin.H{
			"message": "status not modified",
		})
	}
}

func (c *auth) Login(ctx *gin.Context) {
	// 1. 判断后台 session 的状态
	session := sessions.GetSession(ctx.Request)
	backStatus := session.Values[consts.SessionKeyStatus]
	if backStatus != consts.StatusScanned {
		logger.Logger.Info("an ill-timed request")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "an ill-timed request",
		})
	}

	// 2. 获取 code, encryptedData, iv
	wxConfirmLogin := entity.WxConfirmLogin{}
	ctx.ShouldBindJSON(&wxConfirmLogin)

	authService := service.AuthService(ctx)

	// 3. code 换 sessionKey
	wxSessionKey, err := authService.Code2SessionKey(wxConfirmLogin.Code)
	if err != nil {
		logger.Logger.Error("code exchange sessionKey failed", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "exchange sessionKey failed",
		})
	}

	// 4. sessionKey, iv 将 encryptedData 解密为用户手机号
	phoneNumber, err := authService.DecryptPhoneData(
		wxConfirmLogin.EncryptedData, wxSessionKey.SessionKey, wxConfirmLogin.Iv)
	if err != nil {
		logger.Logger.Error("failed to decrypt encryptedData", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to decrypt encryptedData",
		})
	}

	// 5. 查询数据库，获取用户信息
	user, err := authService.GetUserInfo(phoneNumber)
	if err != nil {
		logger.Logger.Error("failed to get user info", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to query user",
		})
	}

	// 6. 授权、status: 3
	service.Context(ctx).SetAuth(user)
	authService.ChangeStatus(consts.StatusLogined)
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (c *auth) Logout(ctx *gin.Context) {
	service.Context(ctx).ForgetAuth()
	ctx.Redirect(http.StatusFound, "/")
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}
