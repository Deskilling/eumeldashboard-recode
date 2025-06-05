package main

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/status", api.OnlineStatus)

	router.GET("/api/players/online", api.GetOnlinePlayers)
	router.POST("/api/players/online", api.PostOnlinePlayers)

	router.GET("/api/players/deaths", api.GetAllDeaths)
	router.POST("/api/players/deaths", api.PostAllDeaths)

	router.GET("/api/chat", api.GetChatMessages)
	router.POST("/api/chat", api.PostChatMessages)

	router.GET("/api/leaderboard", api.GetLeaderBoard)
	router.POST("/api/leaderboard", api.PostLeaderBoard)

	router.GET("/api/globalstats", api.GetGlobalStats)
	router.POST("/api/globalstats", api.PostGlobalStats)

	return router
}

func main() {

	router := setupRouter()

	_ = router.SetTrustedProxies([]string{"localhost"})

	_ = router.Run(":8080")
}
