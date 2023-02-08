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

func (c MatchHistoryCollecter) Collect(priority bool) error {
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
				<-matchScheduler.Schedule(NewMatchDetailsCollecter(matchId), priority)
			}(matchId)
		}
	}

	// await matches collected
	wg.Wait()

	// store matches updated time
	err = db.UpdateMatchesLastUpdated(c.Puuid, updatedAt)
	if err != nil {
		fmt.Printf("Error updating matches last updated for %s value %v, %v\n", c.Puuid, updatedAt, err)
	}
	return err
}
