package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type home struct{}

var Home = home{}

func (c *home) Welcome(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "welcome!",
	})
}
