package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func OnlineStatus(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
