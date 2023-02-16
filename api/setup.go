package api

import (
	s "thecollector/scheduler"

	"github.com/gofiber/fiber/v2"
)

var prioritySummonerScheduler *s.Scheduler
var priorityMatchScheduler *s.Scheduler
var summonerScheduler *s.Scheduler
var matchScheduler *s.Scheduler

func SetSummonerSchedulers(priorityScheduler *s.Scheduler, scheduler *s.Scheduler) {
	prioritySummonerScheduler = priorityScheduler
	summonerScheduler = scheduler
}

func SetMatchSchedulers(priorityScheduler *s.Scheduler, scheduler *s.Scheduler) {
	priorityMatchScheduler = priorityScheduler
	matchScheduler = scheduler
}

func Start() {
	app := fiber.New()

	app.Get("summoner/:region/:name", SummonerByName)
	app.Get("summoner/:puuid", SummonerByPuuid)
	app.Get("update/summoner/:puuid", UpdateSummoner)

	app.Get("matches/:puuid", MatchHistory)
	app.Get("update/matches/:puuid", UpdateMatchHistory)

	app.Listen(":9090")
}
