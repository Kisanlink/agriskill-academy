package admin

import (
	"github.com/Kisanlink/agriskill-academy/internal/auth"
	"errors"
)

type AdminService interface {
	GetJobAnalytics() (*JobAnalytics, error)
	GetUserAnalytics() (*UserAnalytics, error)
	GetApplicationAnalytics() (*ApplicationAnalytics, error)
	GetDashboardAnalytics() (*DashboardAnalytics, error)

	// User Management
	GetUsers(req *UserListRequest) (*UserListResponse, error)
	GetUserByID(userID string) (*UserDetailResponse, error)
	UpdateUser(userID string, req *UpdateUserRequest) error
	DeleteUser(userID string) error
	CreateAdmin(req *CreateAdminRequest) (*CreateAdminResponse, error)

	// Company Management
	GetCompanies(req *CompanyListRequest) (*CompanyListResponse, error)
	GetCompanyByID(companyID string) (*CompanyDetailResponse, error)
	UpdateCompany(companyID string, req *UpdateCompanyRequest) error
	DeleteCompany(companyID string) error
	GetCompanyAnalytics() (*CompanyAnalytics, error)

	// Student/Employer Lists
	GetStudents(req *StudentListRequest) (*StudentListResponse, error)
	GetEmployers(req *EmployerListRequest) (*EmployerListResponse, error)
}

type adminService struct {
	repo AdminRepository
}

func NewAdminService(repo AdminRepository) AdminService {
	return &adminService{repo}
}

func (s *adminService) GetJobAnalytics() (*JobAnalytics, error) {
	return s.repo.GetJobAnalytics()
}

func (s *adminService) GetUserAnalytics() (*UserAnalytics, error) {
	return s.repo.GetUserAnalytics()
}

func (s *adminService) GetApplicationAnalytics() (*ApplicationAnalytics, error) {
	return s.repo.GetApplicationAnalytics()
}

func (s *adminService) GetDashboardAnalytics() (*DashboardAnalytics, error) {
	return s.repo.GetDashboardAnalytics()
}

// User Management Methods
func (s *adminService) GetUsers(req *UserListRequest) (*UserListResponse, error) {
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	return s.repo.GetUsers(req)
}

func (s *adminService) GetUserByID(userID string) (*UserDetailResponse, error) {
	return s.repo.GetUserByID(userID)
}

func (s *adminService) UpdateUser(userID string, req *UpdateUserRequest) error {
	// Role validation is now handled by local authentication
	// We don't validate roles in local DB anymore

	// Validate status if provided
	if req.Status != "" && req.Status != "active" && req.Status != "suspended" && req.Status != "deleted" {
		return errors.New("invalid status: must be 'active', 'suspended', or 'deleted'")
	}

	return s.repo.UpdateUser(userID, req)
}

func (s *adminService) DeleteUser(userID string) error {
	return s.repo.DeleteUser(userID)
}

func (s *adminService) CreateAdmin(req *CreateAdminRequest) (*CreateAdminResponse, error) {
	// Hash the password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create admin user
	adminUser := auth.NewUser()
	adminUser.Name = req.Name
	adminUser.Username = req.Username
	adminUser.Email = req.Email
	adminUser.Password = hashedPassword
	adminUser.Role = "asa_admin"

	// Create the admin user in database
	if err := s.repo.CreateAdmin(adminUser); err != nil {
		return nil, err
	}

	response := &CreateAdminResponse{
		Success: true,
		Message: "Admin user created successfully",
	}
	response.User.ID = adminUser.ID
	response.User.Name = adminUser.Name
	response.User.Username = adminUser.Username
	response.User.Email = adminUser.Email
	response.User.Role = adminUser.Role

	return response, nil
}

// Company Management Methods
func (s *adminService) GetCompanies(req *CompanyListRequest) (*CompanyListResponse, error) {
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	return s.repo.GetCompanies(req)
}

func (s *adminService) GetCompanyByID(companyID string) (*CompanyDetailResponse, error) {
	return s.repo.GetCompanyByID(companyID)
}

func (s *adminService) UpdateCompany(companyID string, req *UpdateCompanyRequest) error {
	// Validate status if provided
	if req.Status != "" && req.Status != "active" && req.Status != "suspended" && req.Status != "deleted" {
		return errors.New("invalid status: must be 'active', 'suspended', or 'deleted'")
	}

	return s.repo.UpdateCompany(companyID, req)
}

func (s *adminService) DeleteCompany(companyID string) error {
	return s.repo.DeleteCompany(companyID)
}

func (s *adminService) GetCompanyAnalytics() (*CompanyAnalytics, error) {
	return s.repo.GetCompanyAnalytics()
}

// Student/Employer List Methods
func (s *adminService) GetStudents(req *StudentListRequest) (*StudentListResponse, error) {
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	return s.repo.GetStudents(req)
}

func (s *adminService) GetEmployers(req *EmployerListRequest) (*EmployerListResponse, error) {
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	return s.repo.GetEmployers(req)
}
