// File: internal/bookmark/service.go

package bookmark

type BookmarkService interface {
	Save(userID, jobID string) error
	Remove(userID, jobID string) error
	GetByUser(userID string) ([]Bookmark, error)
}

type bookmarkService struct {
	repo BookmarkRepository
}

func NewBookmarkService(repo BookmarkRepository) BookmarkService {
	return &bookmarkService{repo}
}

func (s *bookmarkService) Save(userID, jobID string) error {
	b := &Bookmark{
		UserID: userID,
		JobID:  jobID,
	}
	return s.repo.Save(b)
}

func (s *bookmarkService) Remove(userID, jobID string) error {
	return s.repo.Remove(userID, jobID)
}

func (s *bookmarkService) GetByUser(userID string) ([]Bookmark, error) {
	return s.repo.GetByUser(userID)
}
