package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/consts"
	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/entity"
	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/models"
	"github.com/haoran-mc/wx_scan_login/web-back-end/pkg/config"
	"github.com/haoran-mc/wx_scan_login/web-back-end/pkg/db"
	"github.com/haoran-mc/wx_scan_login/web-back-end/pkg/logger"
	"github.com/haoran-mc/wx_scan_login/web-back-end/pkg/utils"
	qrcode "github.com/skip2/go-qrcode"
	"go.uber.org/zap"
)

type sAuth struct {
	stx *BaseContext
}

func AuthService(ctx *gin.Context) *sAuth {
	return &sAuth{stx: Context(ctx)}
}

func (s *sAuth) GenerateQRCode() (string, error) {
	filename := "./assets/" + fmt.Sprintf("%d", time.Now().UnixMilli()) + ".png"
	err := qrcode.WriteFile(
		config.Conf.Applet.Url+":"+config.Conf.Applet.Port,
		qrcode.Medium, 256, filename,
	)
	if err != nil {
		return "", err
	}
	qraddr := config.Conf.Address + ":" + config.Conf.Port + "/assets/" + filename
	logger.Logger.Info("QR code", zap.String("address: ", qraddr))
	return qraddr, nil
}

func (s *sAuth) ChangeStatus(status string) error {
	s.stx.Session.Values[consts.SessionKeyStatus] = status
	s.stx.Session.Save(s.stx.Ctx.Request, s.stx.Ctx.Writer)
	return nil
}

func (s *sAuth) Code2SessionKey(code string) (entity.WxSessionKey, error) {
	var wxSessionKey entity.WxSessionKey
	httpState, bytes := utils.HttpGet(fmt.Sprintf(
		config.Conf.Applet.Code2sessionKeyUrl,
		config.Conf.Applet.Id,
		config.Conf.Applet.Secret,
		code,
	))
	if httpState != 200 {
		logger.Logger.Error("failed to get sessionKey",
			zap.Int("http code", httpState),
		)
		return wxSessionKey, errors.New("failed to get sessionKey")
	}
	err := json.Unmarshal(bytes, &wxSessionKey)
	if err != nil {
		logger.Logger.Error("failed to parse json: ", zap.Error(err))
		return wxSessionKey, errors.New("failed to parse json")
	}
	return wxSessionKey, nil
}

func (s *sAuth) DecryptPhoneData(
	encryptedData, sessionKey, iv string) (string, error) {
	decrypt, err := utils.AesDecrypt(encryptedData, sessionKey, iv)
	if err != nil {
		logger.Logger.Error("failed to decrypt", zap.Error(err))
		return "", err
	}
	var wxPhone = entity.WxPhone{}
	err = json.Unmarshal(decrypt, &wxPhone)
	if err != nil {
		logger.Logger.Error("failed to decrypt phone number", zap.Error(err))
		return "", err
	}
	var phone = wxPhone.PurePhoneNumber
	return phone, nil
}

func (s *sAuth) GetUserInfo(phoneNumber string) (models.Users, error) {
	user := models.Users{}
	err := db.DB().Model(models.Users{}).
		Where("phone = ?", phoneNumber).First(&user).Error
	return user, err
}
