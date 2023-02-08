package main

import (
	"fmt"
	"os"
	"reflect"
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
		for range time.Tick(requestInterval) {
			go summonerScheduler.CollectNext()
		}
	}()

	time.Sleep(requestInterval / 2)

	for range time.Tick(requestInterval) {
		go func() {
			if summonerScheduler.Size() <= 1000 {
				matchScheduler.CollectNext()
			}

			if matchScheduler.IsEmpty() {
				exclude := matchScheduler.ListIdsOfCollecterType(reflect.TypeOf(collect.MatchHistoryCollecter{}))
				puuid, lastUpdated, err := db.QueryPuuidMatchesLastUpdated(exclude)
				if err != nil {
					fmt.Print(err)
					return
				}
				matchScheduler.Schedule(collect.NewMatchHistoryCollecter(puuid, lastUpdated), false)
			}
		}()
	}
}
