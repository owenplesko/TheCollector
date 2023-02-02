package collect

import (
	"fmt"
	"strings"
	"sync"
	db "thecollector/database"
	"thecollector/riot"
)

type MatchDetailsCollecter struct {
	MatchId string
}

func NewMatchDetailsCollecter(matchId string) MatchDetailsCollecter {
	return MatchDetailsCollecter{
		MatchId: matchId,
	}
}

func (c MatchDetailsCollecter) Id() string {
	return c.MatchId
}

func (c MatchDetailsCollecter) Collect() error {
	fmt.Printf("Collecting match %v\n", c.MatchId)
	region := strings.Split(c.MatchId, "_")[0]

	// get match details from riot
	match, err := riot.GetMatchDetails(c.MatchId)
	if err != nil {
		return err
	}

	// queue summoners not in db
	wg := new(sync.WaitGroup)
	errorChan := make(chan error)

	for _, puuid := range match.MetaData.Participants {
		if !db.SummonerExists(puuid) {
			wg.Add(1)
			go func(puuid string) {
				defer wg.Done()
				err := <-summonerScheduler.Schedule(NewSummonerByPuuidCollecter(region, puuid))
				if err != nil {
					errorChan <- err
				}
			}(puuid)
		}
	}

	// await summoners collected
	go func() {
		wg.Wait()
		close(errorChan)
	}()

	for err := range errorChan {
		if err != nil {
			return err
		}
	}

	// store match
	return db.StoreMatch(match)
}
