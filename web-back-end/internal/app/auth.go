package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/consts"
	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/entity"
	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/service"
)

type auth struct{}

var Auth = auth{}

func (c *auth) GetQRCode(ctx *gin.Context) {
	// 1. 生成二维码
	QRCodeUrl, err := service.AuthService(ctx).GenerateQRCode()
	if err != nil {
	}

	// 2. 修改 x-dl-status: 1  -->  做成一个函数
	if err := service.AuthService(ctx).
		ChangeStatus(consts.UnscannedStatus); err != nil {
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    QRCodeUrl,
	})
}

func (c *auth) ChangeStatus(ctx *gin.Context) {
	// 修改 x-dl-status: 2
	if err := service.AuthService(ctx).
		ChangeStatus(consts.ScannedStatus); err != nil {
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (c *auth) Login(ctx *gin.Context) {
	// 1. 获取 code, encryptedData, iv
	wxConfirmLogin := entity.WxConfirmLogin{}
	ctx.ShouldBindJSON(&wxConfirmLogin)

	// 2. code 换 sessionKey
	wxSessionKey, err := service.AuthService(ctx).
		Code2SessionKey(wxConfirmLogin.Code)
	if err != nil {
	}

	// 3. sessionKey, iv 将 encryptedData 解密为用户手机号
	phoneNumber, err := service.AuthService(ctx).
		DecryptPhoneData(wxConfirmLogin.EncryptedData, wxSessionKey.SessionKey, wxConfirmLogin.Iv)
	if err != nil {
	}

	// 4. 查询数据库，获取用户信息，Auth()，更改流量器状态 x-dl-status: 3
	user, err := service.AuthService(ctx).
		GetUserInfo(phoneNumber)
	if err != nil {
	}

	// 5. 授权
	service.Context(ctx).SetAuth(user)
}

func (c *auth) Logout(ctx *gin.Context) {
}
