package api

import (
	"sort"
	"sync"

	"github.com/gin-gonic/gin"
)

type leaderboardEntry struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Kills    uint   `json:"kills"`
	Deaths   uint   `json:"deaths"`
	Playtime uint   `json:"playtime"`
}

type leaderboards struct {
	Entries           map[string]leaderboardEntry
	leaderboardRWLock sync.RWMutex
	sortedByKills     []leaderboardEntry
	sortedByDeaths    []leaderboardEntry
	sortedByPlaytime  []leaderboardEntry
}

var Leaderboards = &leaderboards{
	Entries: make(map[string]leaderboardEntry),
}

func remove(leaderboard []leaderboardEntry, uuid string) []leaderboardEntry {
	for i, entry := range leaderboard {
		if entry.UUID == uuid {
			return append(leaderboard[:i], leaderboard[i+1:]...)
		}
	}
	return leaderboard
}

// -- sort.Search -- Kuss für die gute Documentation
func insert(sortedSlice []leaderboardEntry, newEntry leaderboardEntry, compare func(currenEntry leaderboardEntry, oldEntry leaderboardEntry) bool) []leaderboardEntry {
	i := sort.Search(len(sortedSlice), func(currentIndex int) bool { return !compare(sortedSlice[currentIndex], newEntry) })
	return append(sortedSlice[:i], append([]leaderboardEntry{newEntry}, sortedSlice[i:]...)...)
}

func updateSorted(newEntry leaderboardEntry, oldEntry *leaderboardEntry) {
	if oldEntry != nil {
		uuid := oldEntry.UUID
		Leaderboards.sortedByKills = remove(Leaderboards.sortedByKills, uuid)
		Leaderboards.sortedByDeaths = remove(Leaderboards.sortedByDeaths, uuid)
		Leaderboards.sortedByPlaytime = remove(Leaderboards.sortedByPlaytime, uuid)
	}

	// -- Daumen
	Leaderboards.sortedByKills = insert(Leaderboards.sortedByKills, newEntry, func(newEntry leaderboardEntry, oldEntry leaderboardEntry) bool {
		return newEntry.Kills > oldEntry.Kills
	})
	Leaderboards.sortedByDeaths = insert(Leaderboards.sortedByDeaths, newEntry, func(newEntry leaderboardEntry, oldEntry leaderboardEntry) bool {
		return newEntry.Deaths > oldEntry.Deaths
	})
	Leaderboards.sortedByPlaytime = insert(Leaderboards.sortedByPlaytime, newEntry, func(newEntry leaderboardEntry, oldEntry leaderboardEntry) bool {
		return newEntry.Playtime > oldEntry.Playtime
	})
}

func PostLeaderboard(context *gin.Context) {
	var newEntry leaderboardEntry
	if err := context.ShouldBindJSON(&newEntry); err != nil {
		context.JSON(400, gin.H{"error": "invalid json format"})
		return
	}

	Leaderboards.leaderboardRWLock.Lock()
	defer Leaderboards.leaderboardRWLock.Unlock()

	var previousEntry *leaderboardEntry
	if existingEntry, exists := Leaderboards.Entries[newEntry.UUID]; exists {
		existingEntryCopy := existingEntry
		previousEntry = &existingEntryCopy
	}
	Leaderboards.Entries[newEntry.UUID] = newEntry
	updateSorted(newEntry, previousEntry)

	context.JSON(200, gin.H{
		"kills":    Leaderboards.sortedByKills,
		"deaths":   Leaderboards.sortedByDeaths,
		"playtime": Leaderboards.sortedByPlaytime,
	})
}

/*
func GetLeaderBoard(context *gin.Context) {
	leaderboardRWLock.RLock()
	defer leaderboardRWLock.RUnlock()

	context.JSON(http.StatusOK, gin.H{
		"kills":    sortedByKills,
		"deaths":   sortedByDeaths,
		"playtime": sortedByPlaytime,
	})
}
*/
