package utils

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/haoran-mc/wx_scan_login/web-back-end/pkg/logger"
	"go.uber.org/zap"
)

// Get HTTP GET request
func HttpGet(url string) (int, []byte) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		logger.Logger.Error("error sending GET request",
			zap.String("url", url),
			zap.Error(err),
		)
		return http.StatusInternalServerError, nil
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.Logger.Error("error decoding response from GET request",
				zap.String("url", url),
				zap.Error(err),
			)
		}
	}
	return resp.StatusCode, result.Bytes()
}
