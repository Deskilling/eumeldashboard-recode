package api

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type deathInfo struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Reason    string `json:"reason"`
	Timestamp string `json:"timestamp"`
}

var (
	allDeaths       []deathInfo
	allDeathRWMLock = sync.RWMutex{}
)

func GetAllDeaths(context *gin.Context) {
	allDeathRWMLock.RLock()
	defer allDeathRWMLock.RUnlock()

	context.JSON(http.StatusOK, gin.H{"deaths": allDeaths})
}

func PostAllDeaths(context *gin.Context) {
	var death deathInfo

	if err := context.ShouldBindJSON(&death); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	allDeathRWMLock.Lock()

	allDeaths = append(allDeaths, death)
	if len(allDeaths) > 10 {
		allDeaths = allDeaths[len(allDeaths)-10:]
	}
	allDeathRWMLock.Unlock()

	context.JSON(http.StatusCreated, gin.H{"deaths": allDeaths})
}
