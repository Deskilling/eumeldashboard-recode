package api

import (
	"github.com/gin-gonic/gin"
)

func OnlineStatus(context *gin.Context) {
	context.String(200, "ok")
}
