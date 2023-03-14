package entity

// 封装前端传入的数据
type WxConfirmLogin struct {
	EncryptedData string `json:"encrypted_data"`
	Iv            string `json:"iv"`
	Code          string `json:"code"`
}

// 封装code2session接口返回数据
type WxSessionKey struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// 封装手机号信息数据
type WxPhone struct {
	PhoneNumber     string `json:"phone_number"`
	PurePhoneNumber string `json:"pure_phone_number"`
	CountryCode     string `json:"country_code"`
}
