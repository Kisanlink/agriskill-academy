// File: internal/application/repository.go

package application

import (
	"gorm.io/gorm"
)

type ApplicationRepository interface {
	Create(app *Application) error
	GetByStudent(studentID string) ([]Application, error)
	Delete(appID, studentID string) error
}

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db}
}

func (r *applicationRepository) Create(app *Application) error {
	return r.db.Create(app).Error
}

func (r *applicationRepository) GetByStudent(studentID string) ([]Application, error) {
	var apps []Application
	err := r.db.Where("student_id = ?", studentID).Find(&apps).Error
	return apps, err
}

func (r *applicationRepository) Delete(appID, studentID string) error {
	return r.db.Where("id = ? AND student_id = ?", appID, studentID).Delete(&Application{}).Error
}
