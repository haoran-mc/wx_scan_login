package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func AesDecrypt(encryptedData, sessionKey, iv string) ([]byte, error) {
	// Base64 解码
	keyBytes, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, err
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}
	cryptData, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	origData := make([]byte, len(cryptData))
	block, err := aes.NewCipher(keyBytes) // AES
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, ivBytes) // CBC
	mode.CryptBlocks(origData, cryptData)          // 解密
	origData = PKCS7UnPadding(origData)            // 去除填充位
	return origData, nil
}

func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	if length > 0 {
		unPadding := int(plantText[length-1])
		return plantText[:(length - unPadding)]
	}
	return plantText
}
