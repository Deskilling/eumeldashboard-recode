package api

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type chatMessage struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

var (
	allChatMessages []chatMessage
	allChatRWMLock  = sync.RWMutex{}
)

func GetChatMessages(context *gin.Context) {
	allChatRWMLock.RLock()
	defer allChatRWMLock.RUnlock()

	context.JSON(http.StatusOK, gin.H{"chat": allChatMessages})
}

func PostChatMessages(context *gin.Context) {
	var chat chatMessage

	if err := context.ShouldBindJSON(&chat); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	allChatRWMLock.Lock()

	allChatMessages = append(allChatMessages, chat)
	if len(allChatMessages) > 10 {
		allChatMessages = allChatMessages[len(allChatMessages)-10:]
	}
	allChatRWMLock.Unlock()

	context.JSON(http.StatusCreated, gin.H{"chat": allChatMessages})
}
