package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/consts"
	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/entity"
	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/service"
	"github.com/haoran-mc/wx_scan_login/web-back-end/pkg/logger"
	"go.uber.org/zap"
)

type auth struct{}

var Auth = auth{}

func (c *auth) GetQRCode(ctx *gin.Context) {
	// 1. 生成二维码
	QRCodeUrl, err := service.AuthService(ctx).GenerateQRCode()
	if err != nil {
		logger.Logger.Error("failed to generate qrcode", zap.Error(err))
	}

	// 2. 修改 x-dl-status: 1  -->  做成一个函数
	if err := service.AuthService(ctx).
		ChangeStatus(ctx, consts.UnscannedStatus); err != nil {
		logger.Logger.Error("failed to modify x-dl-status", zap.Error(err))
	}

	// 3. 使用 session，将当前状态存入 session
	// 此步为什么不在「更改状态」中完成？
	// 因为后面在更改状态之前还要判断当前状态
	stx := service.Context(ctx)
	stx.Session.Values[consts.StatusSessionKey] = consts.UnscannedStatus
	stx.Session.Save(ctx.Request, ctx.Writer)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    QRCodeUrl,
	})
}

// 小程序发送请求修改后端 session 中的 x-dl-status
func (c *auth) ChangeStatus(ctx *gin.Context) {
	if err := service.AuthService(ctx).
		ChangeStatus(ctx, consts.ScannedStatus); err != nil {
		logger.Logger.Error("failed to modify x-dl-status", zap.Error(err))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

// 浏览器不断发送请求，查看后端 session 中的状态
func (c *auth) CheckStatus(ctx *gin.Context) {
	stx := service.Context(ctx)
	back_status := stx.Session.Values[consts.StatusSessionKey] // 后台的状态
	front_status, _ := ctx.Cookie(consts.StatusSessionKey)
	var err error
	if front_status == consts.UnscannedStatus && back_status == consts.ScannedStatus {
		err = service.AuthService(ctx).ChangeStatus(ctx, consts.ScannedStatus)
	} else if front_status == consts.ScannedStatus && back_status == consts.LoginedStatus {
		err = service.AuthService(ctx).ChangeStatus(ctx, consts.LoginedStatus)
	}
	if err != nil {
		logger.Logger.Error("failed to change status", zap.Error(err))
		ctx.JSON(http.StatusNotModified, gin.H{"message": "failed to change status"})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (c *auth) Login(ctx *gin.Context) {
	// 1. 判断后台 session 的状态
	stx := service.Context(ctx)
	back_status := stx.Session.Values[consts.StatusSessionKey]
	if back_status != consts.ScannedStatus {
		logger.Logger.Info("an ill-timed request")
		ctx.JSON(http.StatusForbidden, gin.H{"message": "an ill-timed request"})
	}

	// 2. 获取 code, encryptedData, iv
	wxConfirmLogin := entity.WxConfirmLogin{}
	ctx.ShouldBindJSON(&wxConfirmLogin)

	// 3. code 换 sessionKey
	wxSessionKey, err := service.AuthService(ctx).
		Code2SessionKey(wxConfirmLogin.Code)
	if err != nil {
		logger.Logger.Error("code exchange sessionKey failed", zap.Error(err))
	}

	// 4. sessionKey, iv 将 encryptedData 解密为用户手机号
	phoneNumber, err := service.AuthService(ctx).
		DecryptPhoneData(wxConfirmLogin.EncryptedData, wxSessionKey.SessionKey, wxConfirmLogin.Iv)
	if err != nil {
		logger.Logger.Error("failed to decrypt encryptedData", zap.Error(err))
	}

	// 5. 查询数据库，获取用户信息，Auth()，更改流量器状态 x-dl-status: 3
	user, err := service.AuthService(ctx).
		GetUserInfo(phoneNumber)
	if err != nil {
		logger.Logger.Error("failed to get user info", zap.Error(err))
	}

	// 6. 授权
	service.Context(ctx).SetAuth(user)
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (c *auth) Logout(ctx *gin.Context) {
	service.Context(ctx).Forget()
	ctx.Redirect(http.StatusFound, "/")
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}
