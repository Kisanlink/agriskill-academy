package bookmark

import (
	"github.com/Kisanlink/agriskill-academy/internal/jobpost"
)

// Assume you have a Job struct somewhere, e.g. in internal/jobpost/model.go
// type Job struct { ... }

// And a JobRepository interface for fetching jobs by IDs
type JobRepository interface {
	GetJobsByIDs(ids []string) ([]jobpost.JobPost, error)
}

type BookmarkService interface {
	Save(userID, jobID string) error
	Remove(userID, jobID string) error
	GetByUser(userID string) ([]jobpost.JobPost, error) // returns Job list now
}

type bookmarkService struct {
	repo    BookmarkRepository
	jobRepo JobRepository
}

func NewBookmarkService(repo BookmarkRepository, jobRepo JobRepository) BookmarkService {
	return &bookmarkService{repo: repo, jobRepo: jobRepo}
}

func (s *bookmarkService) Save(userID, jobID string) error {
	b := NewBookmark() // Use constructor to generate ID
	b.UserID = userID
	b.JobID = jobID
	return s.repo.Save(b)
}

func (s *bookmarkService) Remove(userID, jobID string) error {
	return s.repo.Remove(userID, jobID)
}

func (s *bookmarkService) GetByUser(userID string) ([]jobpost.JobPost, error) {
	bookmarks, err := s.repo.GetByUser(userID)
	if err != nil {
		return nil, err
	}
	if len(bookmarks) == 0 {
		return []jobpost.JobPost{}, nil
	}
	jobIDs := make([]string, 0, len(bookmarks))
	for _, b := range bookmarks {
		jobIDs = append(jobIDs, b.JobID)
	}
	jobs, err := s.jobRepo.GetJobsByIDs(jobIDs)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}
