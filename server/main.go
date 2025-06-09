package main

import (
	"net/http"
	"server/api"
	"strings"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())

	router.GET("/status", api.OnlineStatus)
	router.GET("/api/players/online", api.GetOnlinePlayers)
	router.GET("/api/players/deaths", api.GetAllDeaths)
	router.GET("/api/chat", api.GetChatMessages)
	router.GET("/api/leaderboard", api.GetLeaderBoard)
	router.GET("/api/globalstats", api.GetGlobalStats)

	authGroup := router.Group("/api")
	authGroup.Use(authMiddleware)

	authGroup.POST("/players/online", api.PostOnlinePlayers)
	authGroup.POST("/players/deaths", api.PostAllDeaths)
	authGroup.POST("/chat", api.PostChatMessages)
	authGroup.POST("/leaderboard", api.PostLeaderBoard)
	authGroup.POST("/globalstats", api.PostGlobalStats)

	return router
}

// Auth is the authorization token used for the API, for now
const Auth = "eumeldashboard"

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := setupRouter()

	_ = router.SetTrustedProxies([]string{"localhost"})

	_ = router.Run(":8080")
}

func authMiddleware(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No authorization header"})
		return
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
		return
	}

	if parts[1] != Auth {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.Next()
}
