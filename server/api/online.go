package api

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type onlinePlayer struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Deaths uint   `json:"deaths"`
	Kills  uint   `json:"kills"`
	Online bool   `json:"online"`
}

var (
	onlinePlayersMap    = make(map[string]onlinePlayer)
	onlinePlayersRWLock = sync.RWMutex{}
)

func GetOnlinePlayers(context *gin.Context) {
	onlinePlayersRWLock.RLock()
	defer onlinePlayersRWLock.RUnlock()

	context.JSON(http.StatusOK, gin.H{"online": len(onlinePlayersMap), "players": onlinePlayersMap})
}

func PostOnlinePlayers(context *gin.Context) {
	var player onlinePlayer

	if err := context.ShouldBindJSON(&player); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	onlinePlayersRWLock.Lock()
	onlinePlayersMap[player.UUID] = player

	if onlinePlayersMap[player.UUID].Online {
		onlinePlayersMap[player.UUID] = onlinePlayer{
			UUID:   player.UUID,
			Name:   player.Name,
			Deaths: player.Deaths,
			Kills:  player.Kills,
			Online: true,
		}
	} else {
		delete(onlinePlayersMap, player.UUID)
	}
	onlinePlayersRWLock.Unlock()

	context.JSON(http.StatusCreated, gin.H{"online": len(onlinePlayersMap), "players": onlinePlayersMap})
}
