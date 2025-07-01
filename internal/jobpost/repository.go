package jobpost

import (
	"gorm.io/gorm"
)

type JobPostRepository interface {
	Create(job *JobPost) error
	Update(job *JobPost) error
	Delete(id string) error
	GetByID(id string) (*JobPost, error)
	GetByEmployer(employerID string) ([]JobPost, error)
	Search(filter *JobPostFilter) ([]JobPost, error)
	GetJobsByIDs(ids []string) ([]JobPost, error)
}

type jobPostRepository struct {
	db *gorm.DB
}

func NewJobPostRepository(db *gorm.DB) JobPostRepository {
	return &jobPostRepository{db}
}

func (r *jobPostRepository) Create(job *JobPost) error {
	return r.db.Create(job).Error
}

func (r *jobPostRepository) Update(job *JobPost) error {
	return r.db.Save(job).Error
}

func (r *jobPostRepository) Delete(id string) error {
	return r.db.Delete(&JobPost{}, "id = ?", id).Error
}

func (r *jobPostRepository) GetByID(id string) (*JobPost, error) {
	var job JobPost
	err := r.db.First(&job, "id = ?", id).Error
	return &job, err
}

func (r *jobPostRepository) GetByEmployer(employerID string) ([]JobPost, error) {
	var jobs []JobPost
	err := r.db.Where("employer_id = ?", employerID).Find(&jobs).Error
	return jobs, err
}

type JobPostFilter struct {
	Location   string
	JobType    string
	Experience string
}

func (r *jobPostRepository) Search(filter *JobPostFilter) ([]JobPost, error) {
	var jobs []JobPost
	query := r.db.Model(&JobPost{})
	if filter.Location != "" {
		query = query.Where("location = ?", filter.Location)
	}
	if filter.JobType != "" {
		query = query.Where("job_type = ?", filter.JobType)
	}
	if filter.Experience != "" {
		query = query.Where("experience = ?", filter.Experience)
	}
	err := query.Find(&jobs).Error
	return jobs, err
}

func (r *jobPostRepository) GetJobsByIDs(ids []string) ([]JobPost, error) {
	var jobs []JobPost
	err := r.db.Where("id IN ?", ids).Find(&jobs).Error
	return jobs, err
}
