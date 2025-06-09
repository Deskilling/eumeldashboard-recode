package api

import (
	"net/http"
	"sort"
	"sync"

	"github.com/gin-gonic/gin"
)

// Hoffentlich hat der Aal gute Performance

type leaderBoardEntry struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Kills    uint   `json:"kills"`
	Deaths   uint   `json:"deaths"`
	Playtime uint   `json:"playtime"`
}

var (
	leaderboardRWLock = sync.RWMutex{}
	entries           = map[string]leaderBoardEntry{}
	sortedByKills     []leaderBoardEntry
	sortedByDeaths    []leaderBoardEntry
	sortedByPlaytime  []leaderBoardEntry
)

func remove(leaderboard []leaderBoardEntry, uuid string) []leaderBoardEntry {
	for i, entry := range leaderboard {
		if entry.UUID == uuid {
			return append(leaderboard[:i], leaderboard[i+1:]...)
		}
	}
	return leaderboard
}

// -- sort.Search -- Kuss für die gute Documentation
func insert(sortedSlice []leaderBoardEntry, newEntry leaderBoardEntry, compare func(currenEntry leaderBoardEntry, oldEntry leaderBoardEntry) bool) []leaderBoardEntry {
	i := sort.Search(len(sortedSlice), func(currentIndex int) bool { return !compare(sortedSlice[currentIndex], newEntry) })
	return append(sortedSlice[:i], append([]leaderBoardEntry{newEntry}, sortedSlice[i:]...)...)
}

func updateSorted(newEntry leaderBoardEntry, oldEntry *leaderBoardEntry) {
	if oldEntry != nil {
		uuid := oldEntry.UUID
		sortedByKills = remove(sortedByKills, uuid)
		sortedByDeaths = remove(sortedByDeaths, uuid)
		sortedByPlaytime = remove(sortedByPlaytime, uuid)
	}

	// -- Daumen
	sortedByKills = insert(sortedByKills, newEntry, func(newEntry leaderBoardEntry, oldEntry leaderBoardEntry) bool {
		return newEntry.Kills > oldEntry.Kills
	})
	sortedByDeaths = insert(sortedByDeaths, newEntry, func(newEntry leaderBoardEntry, oldEntry leaderBoardEntry) bool {
		return newEntry.Deaths > oldEntry.Deaths
	})
	sortedByPlaytime = insert(sortedByPlaytime, newEntry, func(newEntry leaderBoardEntry, oldEntry leaderBoardEntry) bool {
		return newEntry.Playtime > oldEntry.Playtime
	})
}

func GetLeaderBoard(c *gin.Context) {
	leaderboardRWLock.RLock()
	defer leaderboardRWLock.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"kills":    sortedByKills,
		"deaths":   sortedByDeaths,
		"playtime": sortedByPlaytime,
	})
}

func PostLeaderBoard(context *gin.Context) {
	var newEntry leaderBoardEntry
	if err := context.ShouldBindJSON(&newEntry); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	leaderboardRWLock.Lock()
	defer leaderboardRWLock.Unlock()

	var previousEntry *leaderBoardEntry
	if existingEntry, exists := entries[newEntry.UUID]; exists {
		existingEntryCopy := existingEntry
		previousEntry = &existingEntryCopy
	}
	entries[newEntry.UUID] = newEntry
	updateSorted(newEntry, previousEntry)

	context.JSON(http.StatusOK, gin.H{"kills": sortedByKills, "deaths": sortedByDeaths, "playtime": sortedByPlaytime})
}
