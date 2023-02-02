package collect

import (
	"fmt"
	"sync"
	db "thecollector/database"
	"thecollector/riot"
	"time"
)

type MatchHistoryCollecter struct {
	Puuid string
	After int64
}

func NewMatchHistoryCollecter(puuid string, after int64) MatchHistoryCollecter {
	return MatchHistoryCollecter{
		Puuid: puuid,
		After: after,
	}
}

func (c MatchHistoryCollecter) Id() string {
	return c.Puuid
}

func (c MatchHistoryCollecter) Collect() error {
	fmt.Printf("Collecting match history for summoner %v\n", c.Puuid)
	// get match history from riot
	updatedAt := time.Now().Unix()
	history, err := riot.GetMatchHistory(c.Puuid, c.After)
	if err != nil {
		fmt.Print(err)
		return err
	}

	// queue matches not in db
	wg := new(sync.WaitGroup)

	for _, matchId := range history {
		if !db.MatchExists(matchId) {
			wg.Add(1)
			go func(matchId string) {
				defer wg.Done()
				<-matchScheduler.Schedule(NewMatchDetailsCollecter(matchId))
			}(matchId)
		}
	}

	// await matches collected
	wg.Wait()

	// store matches updated time
	return db.UpdateMatchesLastUpdated(c.Puuid, updatedAt)
}
