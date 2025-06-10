package api

import (
	"sync"

	"github.com/gin-gonic/gin"
)

const deathLimit = 10

type deathInfo struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Reason    string `json:"reason"`
	Timestamp string `json:"timestamp"`
}

type deaths struct {
	PlayerDeaths []deathInfo
	mu           sync.RWMutex
}

var Deaths = &deaths{
	PlayerDeaths: make([]deathInfo, 0),
}

func PostAllDeaths(context *gin.Context) {
	var death deathInfo

	if err := context.ShouldBindJSON(&death); err != nil {
		context.IndentedJSON(400, gin.H{"error": "invalid json format"})
		return
	}

	Deaths.mu.Lock()
	defer Deaths.mu.Unlock()

	Deaths.PlayerDeaths = append(Deaths.PlayerDeaths, death)
	if len(Deaths.PlayerDeaths) > deathLimit {
		Deaths.PlayerDeaths = Deaths.PlayerDeaths[len(Deaths.PlayerDeaths)-deathLimit:]
	}

	context.JSON(200, gin.H{"deaths": Deaths.PlayerDeaths})
}

/*
func GetAllDeaths(context *gin.Context) {
	allDeathRWMLock.RLock()
	defer allDeathRWMLock.RUnlock()

	context.JSON(http.StatusOK, gin.H{"deaths": allDeaths})
}
*/
