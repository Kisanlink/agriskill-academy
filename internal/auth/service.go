package auth

import (
	"github.com/Kisanlink/agriskill-academy/internal/employerprofile"
	"github.com/Kisanlink/agriskill-academy/internal/middleware"
	"github.com/Kisanlink/agriskill-academy/internal/studentprofile"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// contains checks if a slice of strings contains a specific string
func contains(list []string, val string) bool {
	for _, item := range list {
		if item == val {
			return true
		}
	}
	return false
}

// HashPassword - hash password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword - verify password using bcrypt
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

type AuthService interface {
	Login(username, password string) (*User, string, error)
	Signup(req *SignupRequest) (*User, string, error)
	VerifyToken(tokenStr string) (bool, error)
	SendResetLink(email string) error
	ResetPassword(token, newPassword string) error
	UpdateProfile(userID string, req *UpdateProfileRequest) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByID(userID string) (*User, error)
	GetCompleteProfile(userID string) (map[string]interface{}, error)
	ListAllUsers() ([]User, error)
}

type authService struct {
	repo               UserRepository
	employerRepo       employerprofile.EmployerProfileRepository
	studentProfileRepo studentprofile.StudentProfileRepository
}

func NewAuthService(repo UserRepository, employerRepo employerprofile.EmployerProfileRepository, studentProfileRepo studentprofile.StudentProfileRepository) AuthService {
	return &authService{
		repo:               repo,
		employerRepo:       employerRepo,
		studentProfileRepo: studentProfileRepo,
	}
}

// JWT Secret
func getSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable is required")
	}
	return secret
}

// Sign JWT Token
func generateToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(getSecret()))
}

// Parse and verify token
func (s *authService) VerifyToken(tokenStr string) (bool, error) {
	// Clean up token format
	tokenStr = strings.TrimSpace(tokenStr)
	const bearerPrefix = "Bearer "

	if len(tokenStr) > len(bearerPrefix) && strings.HasPrefix(tokenStr, bearerPrefix) {
		tokenStr = tokenStr[len(bearerPrefix):]
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Ensure token is signed using HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(getSecret()), nil
	})

	if err != nil || !token.Valid {
		return false, err
	}

	return true, nil
}

// GetUserByEmail retrieves a user by email
func (s *authService) GetUserByEmail(email string) (*User, error) {
	return s.repo.FindByEmail(email)
}

// GetUserByUsername retrieves a user by username (using email)
func (s *authService) GetUserByUsername(username string) (*User, error) {
	return s.repo.FindByEmail(username)
}

// GetUserByID retrieves a user by ID
func (s *authService) GetUserByID(userID string) (*User, error) {
	return s.repo.GetByID(context.Background(), userID, &User{})
}

// ListAllUsers retrieves all users (for debugging)
func (s *authService) ListAllUsers() ([]User, error) {
	return s.repo.ListAllUsers()
}

// Login - now validates password locally
func (s *authService) Login(username, password string) (*User, string, error) {
	// Try to find user by username first, then by email
	var user *User
	var err error

	// First try to find by username
	user, err = s.repo.FindByUsername(username)
	if err != nil {
		// If not found by username, try by email
		user, err = s.repo.FindByEmail(username)
		if err != nil {
			return nil, "", errors.New("invalid credentials")
		}
	}

	// Validate password locally
	if err := VerifyPassword(user.Password, password); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := generateToken(user)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

// Signup
func (s *authService) Signup(req *SignupRequest) (*User, string, error) {
	middleware.DebugLog("🔍 AuthService.Signup called with: %+v\n", req)

	// 1. Validate password confirmation
	if req.Password != req.ConfirmPassword {
		middleware.DebugLog("❌ Password mismatch\n")
		return nil, "", errors.New("passwords do not match")
	}

	// 2. Check if user exists
	middleware.DebugLog("🔍 Checking if user exists with email: %s\n", req.Email)
	existingUser, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		middleware.DebugLog("🔍 User not found (expected): %v\n", err)
	} else if existingUser != nil {
		middleware.DebugLog("❌ User already exists: %+v\n", existingUser)
		return nil, "", errors.New("email already registered")
	}

	// 3. Create user (store hashed password and role locally)
	middleware.DebugLog("🔍 Creating user with name: %s, username: %s, email: %s\n", req.Name, req.Username, req.Email)

	// Hash the password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		middleware.DebugLog("❌ Failed to hash password: %v\n", err)
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}
	middleware.DebugLog("🔍 Password hashed successfully using bcrypt.DefaultCost\n")

	user := NewUser()
	user.Name = req.Name
	user.Username = req.Username
	user.Email = req.Email
	user.Password = hashedPassword     // Store hashed password locally
	user.Role = req.Role               // Store role locally
	user.PhoneNumber = req.PhoneNumber // Store phone number
	err = s.repo.Create(context.Background(), user)
	if err != nil {
		middleware.DebugLog("❌ Failed to create user in DB: %v\n", err)
		return nil, "", err
	}
	middleware.DebugLog("✅ User created successfully with hashed password: %+v\n", user)

	// Refresh user to get the generated ID
	if user.ID == "" {
		middleware.DebugLog("⚠️ User ID is empty after creation, trying to fetch user by email\n")
		fetchedUser, err := s.repo.FindByEmail(user.Email)
		if err != nil {
			middleware.DebugLog("❌ Failed to fetch user after creation: %v\n", err)
			return nil, "", fmt.Errorf("failed to fetch user after creation: %w", err)
		}
		user = fetchedUser
		middleware.DebugLog("✅ User refreshed with ID: %s\n", user.ID)
	}

	// 4. Generate token
	token, err := generateToken(user)
	if err != nil {
		return nil, "", err
	}

	// 5. Create corresponding profile based on role from request
	middleware.DebugLog("🔍 Creating profile for role: %s\n", req.Role)
	if req.Role == "employer" {
		middleware.DebugLog("🔍 Creating employer profile for user: %s\n", user.ID)
		// Build location string safely
		location := ""
		if req.City != "" && req.State != "" {
			location = req.City + ", " + req.State
		} else if req.City != "" {
			location = req.City
		} else if req.State != "" {
			location = req.State
		}

		profile := employerprofile.NewEmployerProfile()
		profile.UserID = user.ID
		profile.CompanyName = req.CompanyName
		profile.WebsiteURL = req.Website
		profile.Industry = req.IndustryType
		profile.CompanySize = req.CompanySize
		profile.CompanyDescription = "" // optional
		profile.RecruiterName = req.Name
		profile.Designation = "Recruiter"
		profile.OfficialEmail = req.Email
		profile.PhoneNumber = req.PhoneNumber // Store phone number from request
		profile.LinkedinProfile = ""
		profile.GSTINNumber = req.GstinNumber
		profile.CompanyAddress = req.CompanyAddress
		profile.City = req.City
		profile.State = req.State
		profile.Pincode = req.Pincode
		profile.JobCategories = []string{}
		profile.HiringLocations = []string{location}
		profile.HiringTypes = []string{"full-time"}
		middleware.DebugLog("🔍 Employer profile data: %+v\n", profile)
		middleware.DebugLog("🔍 User ID for profile: %s\n", user.ID)
		middleware.DebugLog("🔍 Profile ID before creation: %s\n", profile.ID)
		if err := s.employerRepo.Create(profile); err != nil {
			middleware.DebugLog("❌ Failed to create employer profile: %v\n", err)
			return nil, "", fmt.Errorf("failed to create employer profile: %w", err)
		}
		middleware.DebugLog("✅ Employer profile created successfully\n")

	} else if req.Role == "student" {
		middleware.DebugLog("🔍 Creating student profile for user: %s\n", user.ID)
		// Build location string safely
		location := ""
		if req.City != "" && req.State != "" {
			location = req.City + ", " + req.State
		} else if req.City != "" {
			location = req.City
		} else if req.State != "" {
			location = req.State
		}

		profile := studentprofile.NewStudentProfile()
		profile.UserID = user.ID
		profile.Name = user.Name
		profile.Email = user.Email
		profile.Location = location
		profile.PhoneNumber = req.PhoneNumber // Store phone number from request
		profile.Skills = []string{}
		profile.ResumeKey = ""       // S3 key for resume file
		profile.ProfilePhotoKey = "" // S3 key for profile photo
		profile.Experience = 0.0
		profile.Education = ""
		profile.Portfolio = ""
		profile.Linkedin = ""
		profile.Github = ""
		middleware.DebugLog("🔍 Student profile data: %+v\n", profile)
		if err := s.studentProfileRepo.Create(context.Background(), profile); err != nil {
			middleware.DebugLog("❌ Failed to create student profile: %v\n", err)
			return nil, "", fmt.Errorf("failed to create user profile: %w", err)
		}
		middleware.DebugLog("✅ Student profile created successfully\n")
	} else {
		middleware.DebugLog("⚠️ Unknown role: %s, no profile created\n", req.Role)
	}

	return user, token, nil
}

// Send password reset link (mock)
func (s *authService) SendResetLink(email string) error {
	// In production: send real email with reset link/token
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return errors.New("email not found")
	}
	middleware.DebugLog("Password reset requested for %s (user id: %s)\n", user.Email, user.ID)
	return nil
}

// Reset password using token (mock logic)
func (s *authService) ResetPassword(token, newPassword string) error {
	// Password reset is now handled locally
	// This method is kept for backward compatibility but doesn't modify local DB
	return nil
}

// Update profile
func (s *authService) UpdateProfile(userID string, req *UpdateProfileRequest) (*User, error) {
	user, err := s.repo.GetByID(context.Background(), userID, &User{})
	if err != nil {
		return nil, err
	}

	// Update basic user fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// Check if email is already taken by another user
		existingUser, _ := s.repo.FindByEmail(req.Email)
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.New("email already taken by another user")
		}
		user.Email = req.Email
	}

	// Update user table
	if err := s.repo.Update(context.Background(), user); err != nil {
		return nil, err
	}

	// Update role-specific profile based on roles from JWT context
	// Note: Role-based logic should be handled by the calling handler
	// which has access to the JWT context with roles
	if err := s.updateEmployerProfile(userID, req); err != nil {
		// Try student profile if employer profile update fails
		if err := s.updateStudentProfile(userID, req); err != nil {
			return nil, fmt.Errorf("failed to update profile: %w", err)
		}
	}

	return user, nil
}

// Update employer profile
func (s *authService) updateEmployerProfile(userID string, req *UpdateProfileRequest) error {
	profile, err := s.employerRepo.GetByUserID(userID)
	if err != nil {
		return err
	}

	// Update employer-specific fields
	if req.CompanyName != "" {
		profile.CompanyName = req.CompanyName
	}
	if req.CompanyDescription != "" {
		profile.CompanyDescription = req.CompanyDescription
	}
	if req.Industry != "" {
		profile.Industry = req.Industry
	}
	if req.CompanySize != "" {
		profile.CompanySize = req.CompanySize
	}
	if req.RecruiterName != "" {
		profile.RecruiterName = req.RecruiterName
	}
	if req.Designation != "" {
		profile.Designation = req.Designation
	}
	if req.OfficialEmail != "" {
		profile.OfficialEmail = req.OfficialEmail
	}
	if req.PhoneNumber != "" {
		profile.PhoneNumber = req.PhoneNumber
	}
	if req.LinkedinProfile != "" {
		profile.LinkedinProfile = req.LinkedinProfile
	}
	if req.GstinNumber != "" {
		profile.GSTINNumber = req.GstinNumber
	}
	if req.CompanyAddress != "" {
		profile.CompanyAddress = req.CompanyAddress
	}
	if req.City != "" {
		profile.City = req.City
	}
	if req.State != "" {
		profile.State = req.State
	}
	if req.Pincode != "" {
		profile.Pincode = req.Pincode
	}
	if req.Website != "" {
		profile.WebsiteURL = req.Website
	}
	if len(req.JobCategories) > 0 {
		profile.JobCategories = req.JobCategories
	}
	if len(req.HiringLocations) > 0 {
		profile.HiringLocations = req.HiringLocations
	}
	if len(req.HiringTypes) > 0 {
		profile.HiringTypes = req.HiringTypes
	}

	return s.employerRepo.Update(profile)
}

// Update student profile
func (s *authService) updateStudentProfile(userID string, req *UpdateProfileRequest) error {
	profile, err := s.studentProfileRepo.GetByUserID(userID)
	if err != nil {
		return err
	}

	// Update student-specific fields
	if req.Name != "" {
		profile.Name = req.Name
	}
	if req.Email != "" {
		profile.Email = req.Email
	}
	if req.Location != "" {
		profile.Location = req.Location
	}
	// Phone number field is not in StudentProfile model, it's handled at user level
	// if req.PhoneNumber != "" {
	//     // Phone number is stored in users table, not student_profiles
	// }
	if req.LinkedinProfile != "" {
		profile.Linkedin = req.LinkedinProfile
	}
	// ProfilePhoto and Resume are now S3 keys, not strings
	// These should be handled by the file upload handlers, not the auth service
	// if req.ProfilePhoto != "" {
	//     profile.ProfilePhoto = req.ProfilePhoto
	// }
	// if req.Resume != "" {
	//     profile.Resume = req.Resume
	// }
	if len(req.Skills) > 0 {
		profile.Skills = req.Skills
	}
	if req.Experience > 0 {
		profile.Experience = req.Experience
	}
	if req.Education != "" {
		profile.Education = req.Education
	}
	if req.Portfolio != "" {
		profile.Portfolio = req.Portfolio
	}
	if req.Github != "" {
		profile.Github = req.Github
	}

	return s.studentProfileRepo.Update(context.Background(), profile)
}

// GetCompleteProfile retrieves complete profile information including role-specific details
func (s *authService) GetCompleteProfile(userID string) (map[string]interface{}, error) {
	// Get basic user information
	user, err := s.repo.GetByID(context.Background(), userID, &User{})
	if err != nil {
		return nil, err
	}

	// Create base profile response
	profile := map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	}

	// Add phone number if available
	if user.PhoneNumber != "" {
		profile["phone_number"] = user.PhoneNumber
	}

	// Add username if available
	if user.Username != "" {
		profile["username"] = user.Username
	}

	// Get role-specific profile data
	switch user.Role {
	case "employer":
		employerProfile, err := s.employerRepo.GetByUserID(userID)
		if err == nil && employerProfile != nil {
			// Add all non-null employer profile fields
			if employerProfile.CompanyName != "" {
				profile["company_name"] = employerProfile.CompanyName
			}
			if employerProfile.Industry != "" {
				profile["industry"] = employerProfile.Industry
			}
			if employerProfile.CompanySize != "" {
				profile["company_size"] = employerProfile.CompanySize
			}
			if employerProfile.LogoKey != "" {
				profile["logo_key"] = employerProfile.LogoKey
			}
			if employerProfile.LogoName != "" {
				profile["logo_name"] = employerProfile.LogoName
			}
			if employerProfile.LogoType != "" {
				profile["logo_type"] = employerProfile.LogoType
			}
			if employerProfile.LogoSize > 0 {
				profile["logo_size"] = employerProfile.LogoSize
			}
			if employerProfile.WebsiteURL != "" {
				profile["website_url"] = employerProfile.WebsiteURL
			}
			if employerProfile.CompanyDescription != "" {
				profile["company_description"] = employerProfile.CompanyDescription
			}
			if employerProfile.RecruiterName != "" {
				profile["recruiter_name"] = employerProfile.RecruiterName
			}
			if employerProfile.Designation != "" {
				profile["designation"] = employerProfile.Designation
			}
			if employerProfile.OfficialEmail != "" {
				profile["official_email"] = employerProfile.OfficialEmail
			}
			if employerProfile.PhoneNumber != "" {
				profile["phone_number"] = employerProfile.PhoneNumber
			}
			if employerProfile.LinkedinProfile != "" {
				profile["linkedin_profile"] = employerProfile.LinkedinProfile
			}
			if len(employerProfile.JobCategories) > 0 {
				profile["job_categories"] = employerProfile.JobCategories
			}
			if len(employerProfile.HiringLocations) > 0 {
				profile["hiring_locations"] = employerProfile.HiringLocations
			}
			if len(employerProfile.HiringTypes) > 0 {
				profile["hiring_types"] = employerProfile.HiringTypes
			}
			if employerProfile.GSTINNumber != "" {
				profile["gstin_number"] = employerProfile.GSTINNumber
			}
			if employerProfile.CompanyAddress != "" {
				profile["company_address"] = employerProfile.CompanyAddress
			}
			if employerProfile.City != "" {
				profile["city"] = employerProfile.City
			}
			if employerProfile.State != "" {
				profile["state"] = employerProfile.State
			}
			if employerProfile.Pincode != "" {
				profile["pincode"] = employerProfile.Pincode
			}
		}

	case "student":
		studentProfile, err := s.studentProfileRepo.GetByUserID(userID)
		if err == nil && studentProfile != nil {
			// Add all non-null student profile fields
			if studentProfile.Location != "" {
				profile["location"] = studentProfile.Location
			}
			if studentProfile.PhoneNumber != "" {
				profile["phone_number"] = studentProfile.PhoneNumber
			}
			if studentProfile.ProfilePhotoKey != "" {
				profile["profile_photo_key"] = studentProfile.ProfilePhotoKey
			}
			if studentProfile.ResumeKey != "" {
				profile["resume_key"] = studentProfile.ResumeKey
			}
			if len(studentProfile.Skills) > 0 {
				profile["skills"] = studentProfile.Skills
			}
			if studentProfile.Experience > 0 {
				profile["experience"] = studentProfile.Experience
			}
			if studentProfile.Education != "" {
				profile["education"] = studentProfile.Education
			}
			if studentProfile.Portfolio != "" {
				profile["portfolio"] = studentProfile.Portfolio
			}
			if studentProfile.Linkedin != "" {
				profile["linkedin"] = studentProfile.Linkedin
			}
			if studentProfile.Github != "" {
				profile["github"] = studentProfile.Github
			}
			if len(studentProfile.Certificates) > 0 {
				profile["certificates"] = studentProfile.Certificates
			}
		}
	}

	return profile, nil
}
