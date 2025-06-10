package main

import (
	"server/api"
	"strings"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())

	router.GET("/status", api.OnlineStatus)

	/*
		router.GET("/api/players/online", api.GetOnlinePlayers)
		router.GET("/api/players/deaths", api.GetAllDeaths)
		router.GET("/api/chat", api.GetChatMessages)
		router.GET("/api/leaderboard", api.GetLeaderBoard)
		router.GET("/api/globalstats", api.GetGlobalStats)
	*/

	authGroup := router.Group("/api")
	authGroup.Use(authMiddleware)

	authGroup.POST("/players/online", api.PostOnlinePlayers)
	authGroup.POST("/players/deaths", api.PostAllDeaths)
	authGroup.POST("/chat", api.PostChatMessages)
	authGroup.POST("/leaderboard", api.PostLeaderboard)
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

func authMiddleware(context *gin.Context) {
	auth := context.GetHeader("Authorization")
	if auth == "" {
		context.AbortWithStatusJSON(401, gin.H{"error": "No authorization header"})
		return
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		context.AbortWithStatusJSON(401, gin.H{"error": "Invalid authorization format"})
		return
	}

	if parts[1] != Auth {
		context.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
		return
	}

	context.Next()
}
