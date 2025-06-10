package api

import (
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

type OnlinePlayers struct {
	Players map[string]onlinePlayer
	mu      sync.RWMutex
}

var onlinePlayers = &OnlinePlayers{
	Players: make(map[string]onlinePlayer),
}

func ReturnOnline() *OnlinePlayers {
	return onlinePlayers
}

func PostOnlinePlayers(context *gin.Context) {
	var player onlinePlayer

	if err := context.ShouldBindJSON(&player); err != nil {
		context.JSON(400, gin.H{"error": "invalid json format"})
		return
	}

	data := ReturnOnline()
	data.mu.Lock()
	defer data.mu.Unlock()

	if player.Online {
		data.Players[player.UUID] = player
	} else {
		delete(data.Players, player.UUID)
	}

	context.JSON(200, gin.H{
		"online":  len(data.Players),
		"players": data.Players,
	})
}

/*
func GetOnlinePlayers(context *gin.Context) {
	onlinePlayersRWLock.RLock()
	defer onlinePlayersRWLock.RUnlock()

	context.JSON(http.StatusOK, gin.H{"online": len(onlinePlayersMap), "players": onlinePlayersMap})
}
*/
