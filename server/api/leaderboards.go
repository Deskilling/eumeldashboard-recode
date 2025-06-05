package api

import (
	"net/http"
	"sort"
	"sync"

	"github.com/gin-gonic/gin"
)

type leaderBoardEntry struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Kills    uint   `json:"kills"`
	Deaths   uint   `json:"deaths"`
	Playtime uint   `json:"playtime"`
}

var (
	leaderboardRWLock sync.RWMutex
	entriesMap        = make(map[string]leaderBoardEntry)
	sortedKills       []leaderBoardEntry
	sortedDeaths      []leaderBoardEntry
	sortedPlaytime    []leaderBoardEntry
)

func rebuildSortedSlices() {
	sortedKills = make([]leaderBoardEntry, 0, len(entriesMap))
	sortedDeaths = make([]leaderBoardEntry, 0, len(entriesMap))
	sortedPlaytime = make([]leaderBoardEntry, 0, len(entriesMap))

	for _, entry := range entriesMap {
		sortedKills = append(sortedKills, entry)
		sortedDeaths = append(sortedDeaths, entry)
		sortedPlaytime = append(sortedPlaytime, entry)
	}

	sort.Slice(sortedKills, func(i, j int) bool {
		return sortedKills[i].Kills > sortedKills[j].Kills
	})
	sort.Slice(sortedDeaths, func(i, j int) bool {
		return sortedDeaths[i].Deaths > sortedDeaths[j].Deaths
	})
	sort.Slice(sortedPlaytime, func(i, j int) bool {
		return sortedPlaytime[i].Playtime > sortedPlaytime[j].Playtime
	})
}

func GetLeaderBoard(c *gin.Context) {
	leaderboardRWLock.RLock()
	defer leaderboardRWLock.RUnlock()

	c.JSON(http.StatusOK, gin.H{"leaderboard_kills": sortedKills, "leaderboard_deaths": sortedDeaths, "leaderboard_playtime": sortedPlaytime})
}

func PostLeaderBoard(c *gin.Context) {
	var entry leaderBoardEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	leaderboardRWLock.Lock()
	defer leaderboardRWLock.Unlock()

	entriesMap[entry.UUID] = entry

	rebuildSortedSlices()

	c.JSON(http.StatusOK, gin.H{"leaderboard_kills": sortedKills, "leaderboard_deaths": sortedDeaths, "leaderboard_playtime": sortedPlaytime})
}
