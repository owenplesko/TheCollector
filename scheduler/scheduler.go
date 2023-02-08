package scheduler

import (
	"reflect"
	"sync"
)

type Collecter interface {
	Id() string
	Collect(priority bool) error
}

type job struct {
	collecter Collecter
	priority  bool
	listeners []chan error
}

func newJob(collecter Collecter, priority bool) *job {
	j := new(job)
	j.collecter = collecter
	j.priority = priority
	return j
}

type Scheduler struct {
	mu            sync.Mutex
	queue         []*job
	priorityQueue []*job
	jobs          map[string]*job
}

func NewScheduler() *Scheduler {
	s := new(Scheduler)
	s.jobs = make(map[string]*job)
	return s
}

func (s *Scheduler) Schedule(collecter Collecter, priority bool) chan error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := collecter.Id()
	job, exists := s.jobs[id]
	if !exists {
		job = newJob(collecter, priority)
		s.jobs[id] = job
		if priority {
			s.priorityQueue = append(s.priorityQueue, job)
		} else {
			s.queue = append(s.queue, job)
		}
	} else if priority && !job.priority {
		// upgrade job priority
		job.priority = true
		// find job in queue
		jobIndex := -1
		for i, j := range s.queue {
			if job == j {
				jobIndex = i
				break
			}
		}
		// remove job from queue and add to priority if found
		if jobIndex != -1 {
			s.queue = append(s.queue[:jobIndex], s.queue[jobIndex+1:]...)
			s.priorityQueue = append(s.priorityQueue, job)
		}
		// if not found job is currently running, no need to queue
	}
	completion := make(chan error)
	job.listeners = append(job.listeners, completion)
	return completion
}

func (s *Scheduler) CollectNext() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.priorityQueue) != 0 {
		// pop queue job from queue
		j := s.priorityQueue[0]
		s.priorityQueue = s.priorityQueue[1:]

		// start working on job
		go func(s *Scheduler, j *job) {
			// collect
			err := j.collecter.Collect(j.priority)
			// send err to job listeners
			for _, listener := range j.listeners {
				listener <- err
			}
			// delete job from scheduler
			s.mu.Lock()
			defer s.mu.Unlock()
			delete(s.jobs, j.collecter.Id())
		}(s, j)
	} else if len(s.queue) != 0 {
		// pop queue job from queue
		j := s.queue[0]
		s.queue = s.queue[1:]

		// start working on job
		go func(s *Scheduler, j *job) {
			// collect
			err := j.collecter.Collect(j.priority)
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
}

func (s *Scheduler) IsEmpty() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.queue) == 0
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
