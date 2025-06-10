package api

import (
	"sync"

	"github.com/gin-gonic/gin"
)

type globalStats struct {
	Kills    uint `json:"kills"`
	Deaths   uint `json:"deaths"`
	Playtime uint `json:"playtime"`
}

type stats struct {
	Data globalStats
	mu   sync.RWMutex
}

var Stats = &stats{
	Data: globalStats{},
}

func PostGlobalStats(context *gin.Context) {
	var stat globalStats

	if err := context.ShouldBindJSON(&stat); err != nil {
		context.JSON(400, gin.H{"error": "invalid json format"})
		return
	}

	Stats.mu.Lock()
	defer Stats.mu.Unlock()

	Stats.Data = stat

	context.JSON(200, gin.H{"globalstats": Stats.Data})
}

/*
func GetGlobalStats(context *gin.Context) {
	allStatsRWMLock.RLock()
	defer allStatsRWMLock.RUnlock()

	context.JSON(http.StatusOK, gin.H{"globalstats": allStats})
}
*/
