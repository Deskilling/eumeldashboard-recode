package api

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type globalStats struct {
	Kills    uint `json:"kills"`
	Deaths   uint `json:"deaths"`
	Playtime uint `json:"playtime"`
}

var (
	allStats        globalStats
	allStatsRWMLock = sync.RWMutex{}
)

func GetGlobalStats(context *gin.Context) {
	allStatsRWMLock.RLock()
	defer allStatsRWMLock.RUnlock()

	context.JSON(http.StatusOK, gin.H{"globalstats": allStats})
}

func PostGlobalStats(context *gin.Context) {
	var stat globalStats

	if err := context.ShouldBindJSON(&stat); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	allStatsRWMLock.Lock()

	allStats = stat

	allStatsRWMLock.Unlock()

	context.JSON(http.StatusCreated, gin.H{"globalstats": stat})
}
