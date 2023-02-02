package collect

import (
	s "thecollector/scheduler"
)

var summonerScheduler *s.Scheduler
var matchScheduler *s.Scheduler

func SetSummonerScheduler(s *s.Scheduler) {
	summonerScheduler = s
}

func SetMatchScheduler(s *s.Scheduler) {
	matchScheduler = s
}
