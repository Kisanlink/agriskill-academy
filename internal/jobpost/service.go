package jobpost

type JobPostService interface {
	Create(job *JobPost) error
	Update(job *JobPost) error
	Delete(id string) error
	GetByID(id string) (*JobPost, error)
	GetByEmployer(employerID string) ([]JobPost, error)
	Search(filter *JobPostFilter) ([]JobPost, error)
}

type jobPostService struct {
	repo JobPostRepository
}

func NewJobPostService(repo JobPostRepository) JobPostService {
	return &jobPostService{repo}
}

func (s *jobPostService) Create(job *JobPost) error {
	return s.repo.Create(job)
}

func (s *jobPostService) Update(job *JobPost) error {
	return s.repo.Update(job)
}

func (s *jobPostService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *jobPostService) GetByID(id string) (*JobPost, error) {
	return s.repo.GetByID(id)
}

func (s *jobPostService) GetByEmployer(employerID string) ([]JobPost, error) {
	return s.repo.GetByEmployer(employerID)
}

func (s *jobPostService) Search(filter *JobPostFilter) ([]JobPost, error) {
	return s.repo.Search(filter)
}
