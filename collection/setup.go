package collection

import (
	s "thecollector/scheduler"
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
