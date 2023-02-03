package main

import (
	"fmt"
	"os"
	collect "thecollector/collection"
	"thecollector/config"
	db "thecollector/database"
	"thecollector/riot"
	s "thecollector/scheduler"
	"time"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	riot.Setup(config.Riot.Key, config.Riot.MatchesAfter)
	err = db.Setup(config.Db.Url, config.Db.Port, config.Db.User, config.Db.Password, config.Db.DbName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	requestInterval := time.Duration(float32(config.Riot.RatePeriod)/float32(config.Riot.RateLimit)/config.Riot.RateEfficiency) * time.Millisecond
	summonerScheduler := s.NewScheduler()
	matchScheduler := s.NewScheduler()

	collect.SetSummonerScheduler(summonerScheduler)
	collect.SetMatchScheduler(matchScheduler)

	go func() {
		for {
			summonerScheduler.CollectNext()
			time.Sleep(requestInterval)
		}
	}()

	for {
		go func() {
			if matchScheduler.IsEmpty() {
				puuid, lastUpdated, err := db.QueryPuuidMatchesLastUpdated()
				if err != nil {
					fmt.Print(err)
					return
				}
				matchScheduler.Schedule(collect.NewMatchHistoryCollecter(puuid, lastUpdated))
			}
		}()
		matchScheduler.CollectNext()
		time.Sleep(requestInterval)
	}
}
