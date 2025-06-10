package api

import (
	"sync"

	"github.com/gin-gonic/gin"
)

const chatLimit = 10

type chatMessage struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type chat struct {
	Messages []chatMessage
	mu       sync.RWMutex
}

var Chat = &chat{
	Messages: make([]chatMessage, 0),
}

func PostChatMessages(context *gin.Context) {
	var message chatMessage

	if err := context.ShouldBindJSON(&message); err != nil {
		context.JSON(400, gin.H{"error": "invalid json format"})
		return
	}

	Chat.mu.Lock()
	defer Chat.mu.Unlock()

	Chat.Messages = append(Chat.Messages, message)
	if len(Chat.Messages) > chatLimit {
		Chat.Messages = Chat.Messages[len(Chat.Messages)-chatLimit:]
	}

	context.JSON(200, gin.H{"chat": Chat.Messages})
}

/*
func GetChatMessages(context *gin.Context) {
	allChatRWMLock.RLock()
	defer allChatRWMLock.RUnlock()

	context.JSON(http.StatusOK, gin.H{"chat": allChatMessages})
}
*/
