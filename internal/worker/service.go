package worker

import (
	"context"
	"fmt"
	"time"

	"asa/config"

	"github.com/Kisanlink/kisanlink-db/pkg/core/hash"
)

// JobStatus represents the current status of a job
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusRetrying  JobStatus = "retrying"
)

// BackgroundJob represents a background job with full metadata
type BackgroundJob struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Payload     map[string]interface{} `json:"payload"`
	Status      JobStatus              `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	RetryCount  int                    `json:"retry_count"`
	MaxRetries  int                    `json:"max_retries"`
	Error       string                 `json:"error,omitempty"`
	Result      interface{}            `json:"result,omitempty"`
	NextRetryAt *time.Time             `json:"next_retry_at,omitempty"`
	Priority    int                    `json:"priority"` // Higher number = higher priority
}

// generateID creates a unique ID using kisanlink-db hash generation
func generateID() string {
	id, err := hash.GenerateRandomID("JOB", hash.Medium)
	if err != nil {
		// Fallback to timestamp-based ID if hash generation fails
		return fmt.Sprintf("JOB_%d", time.Now().UnixNano())
	}
	return id
}

// JobService interface for background job processing
type JobService interface {
	Enqueue(job *BackgroundJob) error
	Dequeue() (*BackgroundJob, error)
	Complete(jobID string, result interface{}) error
	Fail(jobID string, errorMsg string) error
	GetJob(jobID string) (*BackgroundJob, error)
	GetQueueStats() (map[string]interface{}, error)
	GetDeadLetterJobs(limit int64) ([]*BackgroundJob, error)
	RetryDeadLetterJob(jobID string) error
	Close() error
}

// InMemoryJobService implements a simple in-memory job queue (for development)
// In production, this should be replaced with RedisJobService
type InMemoryJobService struct {
	jobs           []*BackgroundJob
	processingJobs map[string]*BackgroundJob
	failedJobs     []*BackgroundJob
	ctx            context.Context
}

// NewInMemoryJobService creates a new in-memory job service
func NewInMemoryJobService() JobService {
	return &InMemoryJobService{
		jobs:           make([]*BackgroundJob, 0),
		processingJobs: make(map[string]*BackgroundJob),
		failedJobs:     make([]*BackgroundJob, 0),
		ctx:            context.Background(),
	}
}

// Enqueue adds a job to the queue
func (s *InMemoryJobService) Enqueue(job *BackgroundJob) error {
	if job.ID == "" {
		job.ID = generateID()
	}
	if job.CreatedAt.IsZero() {
		job.CreatedAt = time.Now()
	}
	if job.Status == "" {
		job.Status = JobStatusPending
	}
	if job.MaxRetries == 0 {
		job.MaxRetries = config.GetDefaultMaxRetries()
	}

	s.jobs = append(s.jobs, job)
	return nil
}

// Dequeue retrieves and removes the next job from the queue
func (s *InMemoryJobService) Dequeue() (*BackgroundJob, error) {
	if len(s.jobs) == 0 {
		return nil, nil
	}

	// Get the first job (FIFO)
	job := s.jobs[0]
	s.jobs = s.jobs[1:]

	// Mark as processing
	job.Status = JobStatusRunning
	now := time.Now()
	job.StartedAt = &now
	s.processingJobs[job.ID] = job

	return job, nil
}

// Complete marks a job as completed
func (s *InMemoryJobService) Complete(jobID string, result interface{}) error {
	job, exists := s.processingJobs[jobID]
	if !exists {
		return fmt.Errorf("job not found in processing: %s", jobID)
	}

	job.Status = JobStatusCompleted
	now := time.Now()
	job.CompletedAt = &now
	job.Result = result

	// Remove from processing
	delete(s.processingJobs, jobID)
	return nil
}

// Fail marks a job as failed and handles retry logic
func (s *InMemoryJobService) Fail(jobID string, errorMsg string) error {
	job, exists := s.processingJobs[jobID]
	if !exists {
		return fmt.Errorf("job not found in processing: %s", jobID)
	}

	job.Error = errorMsg
	job.RetryCount++

	// Remove from processing
	delete(s.processingJobs, jobID)

	if job.RetryCount <= job.MaxRetries {
		// Retry with exponential backoff
		job.Status = JobStatusRetrying
		backoff := time.Duration(job.RetryCount*job.RetryCount) * time.Second
		nextRetry := time.Now().Add(backoff)
		job.NextRetryAt = &nextRetry

		// Re-enqueue for retry
		return s.Enqueue(job)
	} else {
		// Max retries exceeded - move to failed jobs
		job.Status = JobStatusFailed
		now := time.Now()
		job.CompletedAt = &now
		s.failedJobs = append(s.failedJobs, job)
	}

	return nil
}

// GetJob retrieves a job by ID
func (s *InMemoryJobService) GetJob(jobID string) (*BackgroundJob, error) {
	// Check processing jobs first
	if job, exists := s.processingJobs[jobID]; exists {
		return job, nil
	}

	// Check failed jobs
	for _, job := range s.failedJobs {
		if job.ID == jobID {
			return job, nil
		}
	}

	// Check pending jobs
	for _, job := range s.jobs {
		if job.ID == jobID {
			return job, nil
		}
	}

	return nil, fmt.Errorf("job not found: %s", jobID)
}

// GetQueueStats returns queue statistics
func (s *InMemoryJobService) GetQueueStats() (map[string]interface{}, error) {
	return map[string]interface{}{
		"pending":    len(s.jobs),
		"processing": len(s.processingJobs),
		"failed":     len(s.failedJobs),
	}, nil
}

// GetDeadLetterJobs returns jobs from the failed jobs list
func (s *InMemoryJobService) GetDeadLetterJobs(limit int64) ([]*BackgroundJob, error) {
	if limit <= 0 {
		limit = 10
	}

	if int64(len(s.failedJobs)) <= limit {
		return s.failedJobs, nil
	}

	return s.failedJobs[:limit], nil
}

// RetryDeadLetterJob moves a job from failed jobs back to main queue
func (s *InMemoryJobService) RetryDeadLetterJob(jobID string) error {
	for i, job := range s.failedJobs {
		if job.ID == jobID {
			// Reset job for retry
			job.Status = JobStatusPending
			job.RetryCount = 0
			job.Error = ""
			job.NextRetryAt = nil
			job.StartedAt = nil
			job.CompletedAt = nil

			// Remove from failed jobs
			s.failedJobs = append(s.failedJobs[:i], s.failedJobs[i+1:]...)

			// Re-enqueue
			return s.Enqueue(job)
		}
	}

	return fmt.Errorf("job not found in failed jobs: %s", jobID)
}

// Close closes the job service
func (s *InMemoryJobService) Close() error {
	// Nothing to close for in-memory service
	return nil
}
