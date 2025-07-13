// Package docs provides Swagger documentation for the ASA Job Portal API
//
// This file contains the main Swagger documentation including:
// - API Info and metadata
// - Security definitions
// - Model definitions
// - API endpoint documentation
//
// @title ASA Job Portal API
// @version 1.0
// @description A comprehensive job portal API for students and employers
// @termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host localhost:3000
// @BasePath /api
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
//
// @tag.name Authentication
// @tag.description User authentication and authorization endpoints
//
// @tag.name Student Profile
// @tag.description Student profile management endpoints
//
// @tag.name Employer Profile
// @tag.description Employer profile management endpoints
//
// @tag.name Job Posts
// @tag.description Job posting and management endpoints
//
// @tag.name Applications
// @tag.description Job application management endpoints
//
// @tag.name Bookmarks
// @tag.description Job bookmark management endpoints
//
// @tag.name File Storage
// @tag.description File upload and management endpoints
//
// @tag.name Notifications
// @tag.description Notification and email management endpoints
//
// @tag.name Admin Analytics
// @tag.description Admin analytics and dashboard endpoints
//
// @tag.name Admin User Management
// @tag.description Admin user management endpoints
//
// @tag.name Background Jobs
// @tag.description Background job processing endpoints
package docs

// API Models

// @Description User registration request
type SignupRequest struct {
	Name            string `json:"name" binding:"required" example:"John Doe"`
	Username        string `json:"username" binding:"required" example:"johndoe"`
	Email           string `json:"email" binding:"required" example:"john@example.com"`
	Password        string `json:"password" binding:"required" example:"password123"`
	ConfirmPassword string `json:"confirm_password" binding:"required" example:"password123"`
	Role            string `json:"role" binding:"required,oneof=student employer" example:"student"`
	PhoneNumber     string `json:"phone_number" binding:"required" example:"9876543210"`
	CountryCode     string `json:"country_code,omitempty" example:"+91"`
	AadhaarNumber   string `json:"aadhaar_number,omitempty" example:"123456789012"`
	CompanyName     string `json:"company_name,omitempty" example:"Tech Corp"`
	GstinNumber     string `json:"gstin_number,omitempty" example:"22AAAAA0000A1Z5"`
}

// @Description User login request
type LoginRequest struct {
	Username string `json:"user_name" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// @Description Profile update request
type UpdateProfileRequest struct {
	Name               string   `json:"name,omitempty" example:"John Doe"`
	Email              string   `json:"email,omitempty" example:"john@example.com"`
	PhoneNumber        string   `json:"phone_number,omitempty" example:"9876543210"`
	Location           string   `json:"location,omitempty" example:"Mumbai, India"`
	ProfilePhoto       string   `json:"profile_photo,omitempty"`
	Bio                string   `json:"bio,omitempty" example:"Experienced software developer"`
	LinkedinProfile    string   `json:"linkedin_profile,omitempty" example:"https://linkedin.com/in/johndoe"`
	Website            string   `json:"website,omitempty" example:"https://johndoe.com"`
	Skills             []string `json:"skills,omitempty" example:"Go,Python,JavaScript"`
	Resume             string   `json:"resume,omitempty"`
	Experience         float64  `json:"experience,omitempty" example:"3.5"`
	Education          string   `json:"education,omitempty" example:"B.Tech Computer Science"`
	Portfolio          string   `json:"portfolio,omitempty" example:"https://github.com/johndoe"`
	Github             string   `json:"github,omitempty" example:"https://github.com/johndoe"`
	CompanyName        string   `json:"company_name,omitempty" example:"Tech Corp"`
	CompanyDescription string   `json:"company_description,omitempty" example:"Leading tech company"`
	Industry           string   `json:"industry,omitempty" example:"Technology"`
	CompanySize        string   `json:"company_size,omitempty" example:"50-100"`
	RecruiterName      string   `json:"recruiter_name,omitempty" example:"Jane Smith"`
	Designation        string   `json:"designation,omitempty" example:"HR Manager"`
	OfficialEmail      string   `json:"official_email,omitempty" example:"hr@techcorp.com"`
	GstinNumber        string   `json:"gstin_number,omitempty" example:"22AAAAA0000A1Z5"`
	CompanyAddress     string   `json:"company_address,omitempty" example:"123 Tech Street"`
	City               string   `json:"city,omitempty" example:"Mumbai"`
	State              string   `json:"state,omitempty" example:"Maharashtra"`
	Pincode            string   `json:"pincode,omitempty" example:"400001"`
	JobCategories      []string `json:"job_categories,omitempty" example:"Software Development,Data Science"`
	HiringLocations    []string `json:"hiring_locations,omitempty" example:"Mumbai,Bangalore"`
	HiringTypes        []string `json:"hiring_types,omitempty" example:"Full-time,Remote"`
}

// @Description Student profile model
type StudentProfile struct {
	ID               string        `json:"id" example:"uuid-string"`
	UserID           string        `json:"user_id" example:"uuid-string"`
	Name             string        `json:"name" example:"John Doe"`
	Email            string        `json:"email" example:"john@example.com"`
	Location         string        `json:"location" example:"Mumbai, India"`
	PhoneNumber      string        `json:"phone_number" example:"9876543210"`
	ProfilePhoto     []byte        `json:"profile_photo"`
	ProfilePhotoName string        `json:"profile_photo_name" example:"profile.jpg"`
	ProfilePhotoType string        `json:"profile_photo_type" example:"image/jpeg"`
	ProfilePhotoSize int64         `json:"profile_photo_size" example:"1024000"`
	Resume           []byte        `json:"resume"`
	ResumeName       string        `json:"resume_name" example:"resume.pdf"`
	ResumeType       string        `json:"resume_type" example:"application/pdf"`
	ResumeSize       int64         `json:"resume_size" example:"2048576"`
	Certificates     []Certificate `json:"certificates"`
	Skills           []string      `json:"skills" example:"Go,Python,JavaScript"`
	Experience       float64       `json:"experience" example:"3.5"`
	Education        string        `json:"education" example:"B.Tech Computer Science"`
	Portfolio        string        `json:"portfolio" example:"https://github.com/johndoe"`
	Linkedin         string        `json:"linkedin" example:"https://linkedin.com/in/johndoe"`
	Github           string        `json:"github" example:"https://github.com/johndoe"`
	CreatedAt        string        `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt        string        `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// @Description Certificate model
type Certificate struct {
	ID               string `json:"id" example:"uuid-string"`
	StudentProfileID string `json:"student_profile_id" example:"uuid-string"`
	Name             string `json:"name" example:"AWS Certified Solutions Architect"`
	File             []byte `json:"file"`
	FileName         string `json:"file_name" example:"aws-certificate.pdf"`
	FileType         string `json:"file_type" example:"application/pdf"`
	FileSize         int64  `json:"file_size" example:"1024000"`
	IssueDate        string `json:"issue_date" example:"2024-01-01"`
}

// @Description Employer profile model
type EmployerProfile struct {
	ID                 string   `json:"id" example:"uuid-string"`
	UserID             string   `json:"user_id" example:"uuid-string"`
	CompanyName        string   `json:"company_name" example:"Tech Corp"`
	CompanyDescription string   `json:"company_description" example:"Leading technology company"`
	Industry           string   `json:"industry" example:"Technology"`
	CompanySize        string   `json:"company_size" example:"50-100"`
	RecruiterName      string   `json:"recruiter_name" example:"Jane Smith"`
	Designation        string   `json:"designation" example:"HR Manager"`
	OfficialEmail      string   `json:"official_email" example:"hr@techcorp.com"`
	GstinNumber        string   `json:"gstin_number" example:"22AAAAA0000A1Z5"`
	CompanyAddress     string   `json:"company_address" example:"123 Tech Street"`
	City               string   `json:"city" example:"Mumbai"`
	State              string   `json:"state" example:"Maharashtra"`
	Pincode            string   `json:"pincode" example:"400001"`
	JobCategories      []string `json:"job_categories" example:"Software Development,Data Science"`
	HiringLocations    []string `json:"hiring_locations" example:"Mumbai,Bangalore"`
	HiringTypes        []string `json:"hiring_types" example:"Full-time,Remote"`
	CreatedAt          string   `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt          string   `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// @Description Job post model
type JobPost struct {
	ID              string   `json:"id" example:"uuid-string"`
	EmployerID      string   `json:"employer_id" example:"uuid-string"`
	Title           string   `json:"title" example:"Senior Software Engineer"`
	Description     string   `json:"description" example:"We are looking for a senior software engineer..."`
	Requirements    string   `json:"requirements" example:"5+ years of experience in Go..."`
	Location        string   `json:"location" example:"Mumbai, India"`
	JobType         string   `json:"job_type" example:"Full-time"`
	ExperienceLevel string   `json:"experience_level" example:"Senior"`
	Salary          string   `json:"salary" example:"₹800,000 - ₹1,200,000"`
	Skills          []string `json:"skills" example:"Go,Python,JavaScript"`
	Benefits        string   `json:"benefits" example:"Health insurance, flexible hours"`
	Status          string   `json:"status" example:"published"`
	CreatedAt       string   `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt       string   `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// @Description Job application model
type Application struct {
	ID             string `json:"id" example:"uuid-string"`
	JobID          string `json:"job_id" example:"uuid-string"`
	StudentID      string `json:"student_id" example:"uuid-string"`
	CoverLetter    string `json:"cover_letter" example:"I am excited to apply..."`
	ResumeFile     []byte `json:"resume_file"`
	ResumeFileName string `json:"resume_file_name" example:"resume.pdf"`
	ResumeFileType string `json:"resume_file_type" example:"application/pdf"`
	ResumeFileSize int64  `json:"resume_file_size" example:"2048576"`
	Status         string `json:"status" example:"applied"`
	AppliedAt      string `json:"applied_at" example:"2024-01-01T00:00:00Z"`
}

// @Description Upload response model
type UploadResponse struct {
	Success  bool   `json:"success" example:"true"`
	Message  string `json:"message" example:"File uploaded successfully"`
	FilePath string `json:"file_path" example:"uploads/resume/file.pdf"`
	FileName string `json:"file_name" example:"file.pdf"`
	FileSize int64  `json:"file_size" example:"1024000"`
	FileType string `json:"file_type" example:"application/pdf"`
	FileURL  string `json:"file_url" example:"/api/files/serve/resume/file.pdf"`
}

// @Description Notification preferences response
type NotificationPreferencesResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Preferences retrieved successfully"`
	Data    interface{} `json:"data"`
}

// @Description Update notification preferences request
type UpdateNotificationPreferencesRequest struct {
	EmailNotifications bool `json:"email_notifications" example:"true"`
	JobAlerts          bool `json:"job_alerts" example:"true"`
	ApplicationUpdates bool `json:"application_updates" example:"true"`
}

// @Description Analytics response
type AnalyticsResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Analytics retrieved successfully"`
	Data    interface{} `json:"data"`
}

// @Description User list request
type UserListRequest struct {
	Page   int    `form:"page" example:"1"`
	Limit  int    `form:"limit" example:"10"`
	Role   string `form:"role" example:"student"`
	Search string `form:"search" example:"john"`
}

// @Description User response
type UserResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Users retrieved successfully"`
	Data    interface{} `json:"data"`
}

// @Description Update user request
type UpdateUserRequest struct {
	Name  string `json:"name,omitempty" example:"John Doe"`
	Email string `json:"email,omitempty" example:"john@example.com"`
	Role  string `json:"role,omitempty" example:"student"`
}

// @Description Create job post request
type CreateJobPostRequest struct {
	Title           string   `json:"title" binding:"required" example:"Senior Software Engineer"`
	Description     string   `json:"description" binding:"required" example:"We are looking for a senior software engineer..."`
	Requirements    string   `json:"requirements" binding:"required" example:"5+ years of experience in Go..."`
	Location        string   `json:"location" binding:"required" example:"Mumbai, India"`
	JobType         string   `json:"job_type" binding:"required" example:"Full-time"`
	ExperienceLevel string   `json:"experience_level" binding:"required" example:"Senior"`
	Salary          string   `json:"salary" binding:"required" example:"₹800,000 - ₹1,200,000"`
	Skills          []string `json:"skills" example:"Go,Python,JavaScript"`
	Benefits        string   `json:"benefits" example:"Health insurance, flexible hours"`
}

// @Description Create draft request
type CreateDraftRequest struct {
	Title           string   `json:"title" binding:"required" example:"Senior Software Engineer"`
	Description     string   `json:"description" example:"We are looking for a senior software engineer..."`
	Requirements    string   `json:"requirements" example:"5+ years of experience in Go..."`
	Location        string   `json:"location" example:"Mumbai, India"`
	JobType         string   `json:"job_type" example:"Full-time"`
	ExperienceLevel string   `json:"experience_level" example:"Senior"`
	Salary          string   `json:"salary" example:"₹800,000 - ₹1,200,000"`
	Skills          []string `json:"skills" example:"Go,Python,JavaScript"`
	Benefits        string   `json:"benefits" example:"Health insurance, flexible hours"`
}

// @Description Update job post request
type UpdateJobPostRequest struct {
	Title           string   `json:"title,omitempty" example:"Senior Software Engineer"`
	Description     string   `json:"description,omitempty" example:"We are looking for a senior software engineer..."`
	Requirements    string   `json:"requirements,omitempty" example:"5+ years of experience in Go..."`
	Location        string   `json:"location,omitempty" example:"Mumbai, India"`
	JobType         string   `json:"job_type,omitempty" example:"Full-time"`
	ExperienceLevel string   `json:"experience_level,omitempty" example:"Senior"`
	Salary          string   `json:"salary,omitempty" example:"₹800,000 - ₹1,200,000"`
	Skills          []string `json:"skills,omitempty" example:"Go,Python,JavaScript"`
	Benefits        string   `json:"benefits,omitempty" example:"Health insurance, flexible hours"`
}

// @Description Update certificate request
type UpdateCertificateRequest struct {
	Name      string `json:"name" binding:"required" example:"AWS Certified Solutions Architect"`
	File      []byte `json:"file" binding:"required"`
	FileName  string `json:"file_name,omitempty" example:"aws-certificate.pdf"`
	FileType  string `json:"file_type,omitempty" example:"application/pdf"`
	FileSize  int64  `json:"file_size,omitempty" example:"1024000"`
	IssueDate string `json:"issue_date" binding:"required" example:"2024-01-01"`
}

// @Description Update student profile request
type UpdateStudentProfileRequest struct {
	UserID           string        `json:"user_id,omitempty" example:"uuid-string"`
	Name             string        `json:"name,omitempty" example:"John Doe"`
	Email            string        `json:"email,omitempty" example:"john@example.com"`
	Location         string        `json:"location,omitempty" example:"Mumbai, India"`
	PhoneNumber      string        `json:"phone_number,omitempty" example:"9876543210"`
	ProfilePhoto     []byte        `json:"profile_photo,omitempty"`
	ProfilePhotoName string        `json:"profile_photo_name,omitempty" example:"profile.jpg"`
	ProfilePhotoType string        `json:"profile_photo_type,omitempty" example:"image/jpeg"`
	ProfilePhotoSize int64         `json:"profile_photo_size,omitempty" example:"1024000"`
	Resume           []byte        `json:"resume,omitempty"`
	ResumeName       string        `json:"resume_name,omitempty" example:"resume.pdf"`
	ResumeType       string        `json:"resume_type,omitempty" example:"application/pdf"`
	ResumeSize       int64         `json:"resume_size,omitempty" example:"2048576"`
	Skills           []string      `json:"skills,omitempty" example:"Go,Python,JavaScript"`
	Experience       *float64      `json:"experience,omitempty" example:"3.5"`
	Education        string        `json:"education,omitempty" example:"B.Tech Computer Science"`
	Portfolio        string        `json:"portfolio,omitempty" example:"https://github.com/johndoe"`
	Linkedin         string        `json:"linkedin,omitempty" example:"https://linkedin.com/in/johndoe"`
	Github           string        `json:"github,omitempty" example:"https://github.com/johndoe"`
	Certificates     []Certificate `json:"certificates,omitempty"`
}
