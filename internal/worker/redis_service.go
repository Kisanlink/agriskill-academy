package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"asa/config"

	"github.com/Kisanlink/kisanlink-db/pkg/core/hash"
	"github.com/go-redis/redis/v8"
)

// RedisJobService implements a Redis-based job queue for production
type RedisJobService struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisJobService creates a new Redis-based job service
func NewRedisJobService(addr, password string, db int) (JobService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisJobService{
		client: client,
		ctx:    ctx,
	}, nil
}

// Enqueue adds a job to the Redis queue
func (s *RedisJobService) Enqueue(job *BackgroundJob) error {
	if job.ID == "" {
		id, err := hash.GenerateRandomID("JOB", hash.Medium)
		if err != nil {
			// Fallback to timestamp-based ID if hash generation fails
			job.ID = fmt.Sprintf("JOB_%d", time.Now().UnixNano())
		} else {
			job.ID = id
		}
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

	// Serialize job to JSON
	jobData, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	// Add to pending queue with score as creation time
	score := float64(job.CreatedAt.Unix())
	err = s.client.ZAdd(s.ctx, "jobs:pending", &redis.Z{
		Score:  score,
		Member: job.ID,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to add job to pending queue: %w", err)
	}

	// Store job data
	err = s.client.Set(s.ctx, fmt.Sprintf("job:%s", job.ID), jobData, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to store job data: %w", err)
	}

	return nil
}

// Dequeue retrieves and removes the next job from the Redis queue
func (s *RedisJobService) Dequeue() (*BackgroundJob, error) {
	// Get the oldest job from pending queue
	result, err := s.client.ZPopMin(s.ctx, "jobs:pending").Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // No jobs available
		}
		return nil, fmt.Errorf("failed to dequeue job: %w", err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	jobID := result[0].Member.(string)

	// Get job data
	jobData, err := s.client.Get(s.ctx, fmt.Sprintf("job:%s", jobID)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get job data: %w", err)
	}

	// Deserialize job
	var job BackgroundJob
	if err := json.Unmarshal([]byte(jobData), &job); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	// Update status to running
	job.Status = JobStatusRunning
	now := time.Now()
	job.StartedAt = &now

	// Store updated job data
	updatedJobData, err := json.Marshal(job)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated job: %w", err)
	}

	err = s.client.Set(s.ctx, fmt.Sprintf("job:%s", job.ID), updatedJobData, 24*time.Hour).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to update job data: %w", err)
	}

	// Add to processing set
	err = s.client.SAdd(s.ctx, "jobs:processing", job.ID).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to add job to processing set: %w", err)
	}

	return &job, nil
}

// Complete marks a job as completed
func (s *RedisJobService) Complete(jobID string, result interface{}) error {
	// Get job data
	jobData, err := s.client.Get(s.ctx, fmt.Sprintf("job:%s", jobID)).Result()
	if err != nil {
		return fmt.Errorf("job not found: %w", err)
	}

	// Deserialize job
	var job BackgroundJob
	if err := json.Unmarshal([]byte(jobData), &job); err != nil {
		return fmt.Errorf("failed to unmarshal job: %w", err)
	}

	// Update job status
	job.Status = JobStatusCompleted
	now := time.Now()
	job.CompletedAt = &now
	job.Result = result

	// Store updated job data
	updatedJobData, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal updated job: %w", err)
	}

	err = s.client.Set(s.ctx, fmt.Sprintf("job:%s", job.ID), updatedJobData, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to update job data: %w", err)
	}

	// Remove from processing set
	err = s.client.SRem(s.ctx, "jobs:processing", jobID).Err()
	if err != nil {
		return fmt.Errorf("failed to remove job from processing set: %w", err)
	}

	// Add to completed set
	err = s.client.ZAdd(s.ctx, "jobs:completed", &redis.Z{
		Score:  float64(now.Unix()),
		Member: jobID,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to add job to completed set: %w", err)
	}

	return nil
}

// Fail marks a job as failed and handles retry logic
func (s *RedisJobService) Fail(jobID string, errorMsg string) error {
	// Get job data
	jobData, err := s.client.Get(s.ctx, fmt.Sprintf("job:%s", jobID)).Result()
	if err != nil {
		return fmt.Errorf("job not found: %w", err)
	}

	// Deserialize job
	var job BackgroundJob
	if err := json.Unmarshal([]byte(jobData), &job); err != nil {
		return fmt.Errorf("failed to unmarshal job: %w", err)
	}

	// Increment retry count
	job.RetryCount++

	// Check if we should retry
	if job.RetryCount < job.MaxRetries {
		// Schedule retry with exponential backoff
		backoff := time.Duration(job.RetryCount*job.RetryCount) * time.Second
		nextRetry := time.Now().Add(backoff)
		job.NextRetryAt = &nextRetry
		job.Status = JobStatusRetrying

		// Add back to pending queue with retry time
		score := float64(nextRetry.Unix())
		err = s.client.ZAdd(s.ctx, "jobs:pending", &redis.Z{
			Score:  score,
			Member: jobID,
		}).Err()
		if err != nil {
			return fmt.Errorf("failed to add job back to pending queue: %w", err)
		}
	} else {
		// Max retries exceeded, move to failed
		job.Status = JobStatusFailed
		job.Error = errorMsg

		// Add to failed set
		err = s.client.ZAdd(s.ctx, "jobs:failed", &redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: jobID,
		}).Err()
		if err != nil {
			return fmt.Errorf("failed to add job to failed set: %w", err)
		}
	}

	// Store updated job data
	updatedJobData, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal updated job: %w", err)
	}

	err = s.client.Set(s.ctx, fmt.Sprintf("job:%s", job.ID), updatedJobData, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to update job data: %w", err)
	}

	// Remove from processing set
	err = s.client.SRem(s.ctx, "jobs:processing", jobID).Err()
	if err != nil {
		return fmt.Errorf("failed to remove job from processing set: %w", err)
	}

	return nil
}

// GetJob retrieves a job by ID
func (s *RedisJobService) GetJob(jobID string) (*BackgroundJob, error) {
	jobData, err := s.client.Get(s.ctx, fmt.Sprintf("job:%s", jobID)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("job not found: %s", jobID)
		}
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	var job BackgroundJob
	if err := json.Unmarshal([]byte(jobData), &job); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	return &job, nil
}

// GetQueueStats returns queue statistics
func (s *RedisJobService) GetQueueStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get counts for each queue
	pendingCount, err := s.client.ZCard(s.ctx, "jobs:pending").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get pending count: %w", err)
	}

	processingCount, err := s.client.SCard(s.ctx, "jobs:processing").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get processing count: %w", err)
	}

	completedCount, err := s.client.ZCard(s.ctx, "jobs:completed").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get completed count: %w", err)
	}

	failedCount, err := s.client.ZCard(s.ctx, "jobs:failed").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get failed count: %w", err)
	}

	stats["pending"] = pendingCount
	stats["processing"] = processingCount
	stats["completed"] = completedCount
	stats["failed"] = failedCount
	stats["total"] = pendingCount + processingCount + completedCount + failedCount

	return stats, nil
}

// GetDeadLetterJobs returns failed jobs
func (s *RedisJobService) GetDeadLetterJobs(limit int64) ([]*BackgroundJob, error) {
	// Get failed job IDs
	jobIDs, err := s.client.ZRevRange(s.ctx, "jobs:failed", 0, limit-1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get failed job IDs: %w", err)
	}

	var jobs []*BackgroundJob
	for _, jobID := range jobIDs {
		job, err := s.GetJob(jobID)
		if err != nil {
			continue // Skip if job data is corrupted
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// RetryDeadLetterJob retries a failed job
func (s *RedisJobService) RetryDeadLetterJob(jobID string) error {
	// Get job data
	job, err := s.GetJob(jobID)
	if err != nil {
		return fmt.Errorf("job not found: %w", err)
	}

	if job.Status != JobStatusFailed {
		return fmt.Errorf("job is not in failed status")
	}

	// Reset job for retry
	job.Status = JobStatusPending
	job.RetryCount = 0
	job.Error = ""
	job.NextRetryAt = nil
	job.StartedAt = nil
	job.CompletedAt = nil
	job.Result = nil

	// Store updated job data
	jobData, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	err = s.client.Set(s.ctx, fmt.Sprintf("job:%s", job.ID), jobData, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to update job data: %w", err)
	}

	// Remove from failed set
	err = s.client.ZRem(s.ctx, "jobs:failed", jobID).Err()
	if err != nil {
		return fmt.Errorf("failed to remove job from failed set: %w", err)
	}

	// Add back to pending queue
	score := float64(time.Now().Unix())
	err = s.client.ZAdd(s.ctx, "jobs:pending", &redis.Z{
		Score:  score,
		Member: jobID,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to add job back to pending queue: %w", err)
	}

	return nil
}

// Close closes the Redis connection
func (s *RedisJobService) Close() error {
	return s.client.Close()
}
