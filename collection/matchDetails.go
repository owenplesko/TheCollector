package collection

import (
	"fmt"
	"strings"
	"sync"
	db "thecollector/database"
	"thecollector/riot"
)

type MatchDetailsCollecter struct {
	Priority bool
	MatchId  string
}

func NewMatchDetailsCollecter(matchId string, priority bool) MatchDetailsCollecter {
	return MatchDetailsCollecter{
		Priority: priority,
		MatchId:  matchId,
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
		fmt.Printf("ERROR failed to get match %v from riot: %v\n", c.MatchId, err)
		return err
	}

	// queue summoners not in db
	wg := new(sync.WaitGroup)
	errorChan := make(chan error)

	for _, puuid := range match.MetaData.Participants {
		if !db.SummonerPuuidExists(puuid) {
			wg.Add(1)
			go func(puuid string) {
				defer wg.Done()
				var err error
				if c.Priority {
					err = <-prioritySummonerScheduler.Schedule(NewSummonerByPuuidCollecter(region, puuid))
				} else {
					err = <-summonerScheduler.Schedule(NewSummonerByPuuidCollecter(region, puuid))
				}
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
			fmt.Printf("Error failed to get summoner for match %v, %v\n", c.MatchId, err)
			return err
		}
	}

	// store match
	err = db.StoreMatch(match)
	if err != nil {
		fmt.Printf("Error inserting match %v, %v\n", c.MatchId, err)
	}
	return err
}
