package scheduler

import (
	"reflect"
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
	mu    sync.Mutex
	queue []*job
	jobs  map[string]*job
}

func NewScheduler() *Scheduler {
	s := new(Scheduler)
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
		s.queue = append(s.queue, job)
	}
	completion := make(chan error)
	job.listeners = append(job.listeners, completion)
	return completion
}

func (s *Scheduler) CollectNext() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.queue) < 1 {
		return
	}

	// pop queue job from queue
	j := s.queue[0]
	s.queue = s.queue[1:]

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

func (s *Scheduler) QueueIsEmpty() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.queue) == 0
}

func (s *Scheduler) QueueSize() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.queue)
}

func (s *Scheduler) Size() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.jobs)
}

func (s *Scheduler) ListIdsOfCollecterType(t reflect.Type) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var filteredIds []string
	for id, job := range s.jobs {
		if reflect.TypeOf(job.collecter) == t {
			filteredIds = append(filteredIds, id)
		}
	}
	return filteredIds
}
