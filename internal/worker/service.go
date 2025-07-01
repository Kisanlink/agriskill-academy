// File: internal/worker/service.go

package worker

import (
	"fmt"
	"time"
)

// Example Job struct for background jobs
type Job struct {
	ID        string
	Type      string
	Payload   map[string]interface{}
	CreatedAt time.Time
}

// JobService manages a queue of jobs (for demonstration; in production, use real queues or tools)
type JobService interface {
	Enqueue(job *Job) error
	ProcessNext() error
}

type inMemoryJobService struct {
	jobs chan *Job
}

func NewInMemoryJobService(buffer int) JobService {
	s := &inMemoryJobService{
		jobs: make(chan *Job, buffer),
	}
	go s.workerLoop()
	return s
}

func (s *inMemoryJobService) Enqueue(job *Job) error {
	select {
	case s.jobs <- job:
		return nil
	default:
		return fmt.Errorf("job queue is full")
	}
}

func (s *inMemoryJobService) ProcessNext() error {
	// For demonstration; real implementation would block or fetch from queue
	select {
	case job := <-s.jobs:
		fmt.Printf("Processing job: %+v\n", job)
		return nil
	default:
		return fmt.Errorf("no job to process")
	}
}

func (s *inMemoryJobService) workerLoop() {
	for job := range s.jobs {
		fmt.Printf("Worker processing job: %+v\n", job)
		// Simulate job processing
		time.Sleep(1 * time.Second)
	}
}
