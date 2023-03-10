package app

import (
	"github.com/gin-gonic/gin"
)

type auth struct{}

var Auth = auth{}

func (c *auth) GetQRCode(ctx *gin.Context) {
	// set session
}

func (c *auth) ChangeStatus(ctx *gin.Context) {
}

// get wx code
func (c *auth) Login(ctx *gin.Context) {
}

func (c *auth) Logout(ctx *gin.Context) {
}
