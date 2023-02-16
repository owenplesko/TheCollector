package main

import (
	"fmt"
	"os"
	"reflect"
	"thecollector/api"
	"thecollector/collection"
	"thecollector/config"
	db "thecollector/database"
	"thecollector/riot"
	s "thecollector/scheduler"
	"thecollector/util"
	"time"
)

func main() {
	// load config
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// set config vars
	requestInterval := time.Duration(float32(config.Riot.RatePeriod)/float32(config.Riot.RateLimit)/config.Riot.RateEfficiency) * time.Millisecond
	riot.Setup(config.Riot.Key, config.Riot.MatchesAfter)
	err = db.Setup(config.Db.Url, config.Db.Port, config.Db.User, config.Db.Password, config.Db.DbName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// create schedulers
	prioritySummonerScheduler := s.NewScheduler()
	priorityMatchScheduler := s.NewScheduler()
	summonerScheduler := s.NewScheduler()
	matchScheduler := s.NewScheduler()

	// set scheduler vars
	collection.SetSummonerSchedulers(prioritySummonerScheduler, summonerScheduler)
	collection.SetMatchSchedulers(priorityMatchScheduler, matchScheduler)
	api.SetSummonerSchedulers(prioritySummonerScheduler, summonerScheduler)
	api.SetMatchSchedulers(priorityMatchScheduler, matchScheduler)

	// summoner collection loop
	util.OnInterval(func() {
		if !prioritySummonerScheduler.QueueIsEmpty() {
			prioritySummonerScheduler.CollectNext()
		} else {
			summonerScheduler.CollectNext()
		}
	}, requestInterval)

	// offset loops
	time.Sleep(requestInterval / 2)

	// match collection loop
	util.OnInterval(func() {
		if !priorityMatchScheduler.QueueIsEmpty() {
			priorityMatchScheduler.CollectNext()
		} else {
			matchScheduler.CollectNext()
		}

		// queue new match history if match details queue is empty
		if matchScheduler.QueueIsEmpty() {
			exclude := matchScheduler.ListIdsOfCollecterType(reflect.TypeOf(collection.MatchHistoryCollecter{}))
			puuid, lastUpdated, err := db.QueryPuuidMatchesLastUpdated(exclude)
			if err != nil {
				fmt.Print(err)
				return
			}
			matchScheduler.Schedule(collection.NewMatchHistoryCollecter(puuid, lastUpdated, false))
		}
	}, requestInterval)

	api.Start()
}
