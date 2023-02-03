package scheduler

import (
	"sync"
)

type Collecter interface {
	Id() string
	Collect() error
}

type job struct {
	collecter Collecter
	listeners []chan error
}

func newJob(collecter Collecter) *job {
	j := new(job)
	j.collecter = collecter
	return j
}

type Scheduler struct {
	mu        sync.Mutex
	available []*job
	jobs      map[string]*job
	jobAdded  sync.Cond
}

func NewScheduler() *Scheduler {
	s := new(Scheduler)
	s.jobAdded.L = &s.mu
	s.jobs = make(map[string]*job)
	return s
}

func (s *Scheduler) Schedule(collecter Collecter) chan error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := collecter.Id()
	job, exists := s.jobs[id]
	if !exists {
		job = newJob(collecter)
		s.jobs[id] = job
		s.available = append(s.available, job)
		s.jobAdded.Signal()
	}
	completion := make(chan error)
	job.listeners = append(job.listeners, completion)
	return completion
}

func (s *Scheduler) CollectNext() {
	s.mu.Lock()
	defer s.mu.Unlock()
	// wait for available job

	for len(s.available) == 0 {
		s.jobAdded.Wait()
	}

	// pop available job from queue
	j := s.available[0]
	s.available = s.available[1:]

	// start working on job
	go func(s *Scheduler, j *job) {
		// collect
		err := j.collecter.Collect()
		// send err to job listeners
		for _, listener := range j.listeners {
			listener <- err
		}
		// delete job from scheduler
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.jobs, j.collecter.Id())
	}(s, j)
}

func (s *Scheduler) IsEmpty() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.available) == 0
}
