package auth

import (
	"asa/internal/employerprofile"
	"asa/internal/studentprofile"
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

// HashPassword - copy this exact function from AAA service
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword - copy this exact function from AAA service
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

type AuthService interface {
	Login(email, password string) (*User, string, error)
	Signup(req *SignupRequest) (*User, string, error)
	SignupWithID(req *SignupRequest, userID string, phoneNumber string) (*User, string, error)
	VerifyToken(tokenStr string) (bool, error)
	SendResetLink(email string) error
	ResetPassword(token, newPassword string) error
	UpdateProfile(userID string, req *UpdateProfileRequest) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByID(userID string) (*User, error)
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
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		return "secret"
	}
	return secret
}

// Sign JWT Token
func generateToken(user *User) (string, error) {
	// Convert single role to roles array for consistency with middleware
	roles := []string{user.Role}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role, // Keep single role for backward compatibility
		"roles":   roles,     // Add roles array for middleware
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
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

// GetUserByID retrieves a user by ID
func (s *authService) GetUserByID(userID string) (*User, error) {
	return s.repo.FindByID(userID)
}

// ListAllUsers retrieves all users (for debugging)
func (s *authService) ListAllUsers() ([]User, error) {
	return s.repo.ListAllUsers()
}

// Password functions removed - now handled by AAA service

// Login
func (s *authService) Login(email, password string) (*User, string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}
	// Password validation is now handled by AAA service
	// We only need to check if user exists in our local DB
	token, err := generateToken(user)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

// Signup
func (s *authService) Signup(req *SignupRequest) (*User, string, error) {
	fmt.Printf("🔍 AuthService.Signup called with: %+v\n", req)

	// 1. Validate password confirmation
	if req.Password != req.ConfirmPassword {
		fmt.Printf("❌ Password mismatch\n")
		return nil, "", errors.New("passwords do not match")
	}

	// 2. Check if user exists
	fmt.Printf("🔍 Checking if user exists with email: %s\n", req.Email)
	existingUser, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		fmt.Printf("🔍 User not found (expected): %v\n", err)
	} else if existingUser != nil {
		fmt.Printf("❌ User already exists: %+v\n", existingUser)
		return nil, "", errors.New("email already registered")
	}

	// 3. Create user (store hashed password and role locally as well)
	fmt.Printf("🔍 Creating user with name: %s, email: %s\n", req.Name, req.Email)

	// Hash the password using the same method as AAA service
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		fmt.Printf("❌ Failed to hash password: %v\n", err)
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}
	fmt.Printf("🔍 Password hashed successfully using bcrypt.DefaultCost\n")

	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword, // Store hashed password locally
		Role:     req.Role,       // Store role locally
	}
	err = s.repo.Create(user)
	if err != nil {
		fmt.Printf("❌ Failed to create user in DB: %v\n", err)
		return nil, "", err
	}
	fmt.Printf("✅ User created successfully with hashed password: %+v\n", user)

	// 4. Generate token
	token, err := generateToken(user)
	if err != nil {
		return nil, "", err
	}

	// 5. Create corresponding profile based on role from request
	fmt.Printf("🔍 Creating profile for role: %s\n", req.Role)
	roles := []string{req.Role} // Use role from request since it's not stored in DB
	if contains(roles, "employer") {
		fmt.Printf("🔍 Creating employer profile for user: %s\n", user.ID)
		// Build location string safely
		location := ""
		if req.City != "" && req.State != "" {
			location = req.City + ", " + req.State
		} else if req.City != "" {
			location = req.City
		} else if req.State != "" {
			location = req.State
		}

		profile := &employerprofile.EmployerProfile{
			UserID:             user.ID,
			CompanyName:        req.CompanyName,
			WebsiteUrl:         req.Website,
			Industry:           req.IndustryType,
			CompanySize:        req.CompanySize,
			CompanyDescription: "", // optional
			RecruiterName:      req.Name,
			Designation:        "Recruiter",
			OfficialEmail:      req.Email,
			PhoneNumber:        "", // optional
			LinkedinProfile:    "",
			GstinNumber:        req.GstinNumber,
			CompanyAddress:     req.CompanyAddress,
			City:               req.City,
			State:              req.State,
			Pincode:            req.Pincode,
			JobCategories:      []string{},
			HiringLocations:    []string{location},
			HiringTypes:        []string{"full-time"},
		}
		fmt.Printf("🔍 Employer profile data: %+v\n", profile)
		if err := s.employerRepo.Create(profile); err != nil {
			fmt.Printf("❌ Failed to create employer profile: %v\n", err)
			return nil, "", fmt.Errorf("failed to create employer profile: %w", err)
		}
		fmt.Printf("✅ Employer profile created successfully\n")

	} else if contains(roles, "student") {
		fmt.Printf("🔍 Creating student profile for user: %s\n", user.ID)
		// Build location string safely
		location := ""
		if req.City != "" && req.State != "" {
			location = req.City + ", " + req.State
		} else if req.City != "" {
			location = req.City
		} else if req.State != "" {
			location = req.State
		}

		profile := &studentprofile.StudentProfile{
			UserID:       user.ID,
			Name:         user.Name,
			Email:        user.Email,
			Location:     location,
			Skills:       []string{},
			Resume:       "",
			ProfilePhoto: "",
			Experience:   0.0,
			Education:    "",
			Portfolio:    "",
			Linkedin:     "",
			Github:       "",
		}
		fmt.Printf("🔍 Student profile data: %+v\n", profile)
		if err := s.studentProfileRepo.Create(profile); err != nil {
			fmt.Printf("❌ Failed to create student profile: %v\n", err)
			return nil, "", fmt.Errorf("failed to create user profile: %w", err)
		}
		fmt.Printf("✅ Student profile created successfully\n")
	} else {
		fmt.Printf("⚠️ Unknown role: %s, no profile created\n", req.Role)
	}

	return user, token, nil
}

// SignupWithID - creates a user with a specific ID (for AAA integration)
func (s *authService) SignupWithID(req *SignupRequest, userID string, phoneNumber string) (*User, string, error) {
	fmt.Printf("🔍 AuthService.SignupWithID called with: %+v, userID: %s\n", req, userID)

	// 1. Validate password confirmation
	if req.Password != req.ConfirmPassword {
		fmt.Printf("❌ Password mismatch\n")
		return nil, "", errors.New("passwords do not match")
	}

	// 2. Check if user exists with the specific ID
	fmt.Printf("🔍 Checking if user exists with ID: %s\n", userID)
	existingUser, err := s.repo.FindByID(userID)
	if err != nil {
		fmt.Printf("🔍 User not found with ID (expected): %v\n", err)
	} else if existingUser != nil {
		fmt.Printf("❌ User already exists with ID: %+v\n", existingUser)
		return nil, "", errors.New("user ID already exists")
	}

	// 3. Create user with the specified ID
	fmt.Printf("🔍 Creating user with ID: %s, name: %s, email: %s\n", userID, req.Name, req.Email)

	// Hash the password using the same method as AAA service
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		fmt.Printf("❌ Failed to hash password: %v\n", err)
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}
	fmt.Printf("🔍 Password hashed successfully using bcrypt.DefaultCost\n")

	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword, // Store hashed password locally
		Role:     req.Role,       // Store role locally
	}
	err = s.repo.CreateWithID(user, userID)
	if err != nil {
		fmt.Printf("❌ Failed to create user in DB: %v\n", err)
		return nil, "", err
	}
	fmt.Printf("✅ User created successfully with specified ID: %+v\n", user)

	// 4. Generate token
	token, err := generateToken(user)
	if err != nil {
		return nil, "", err
	}

	// 5. Create corresponding profile based on role from request
	fmt.Printf("🔍 Creating profile for role: %s\n", req.Role)
	roles := []string{req.Role} // Use role from request since it's not stored in DB
	if contains(roles, "employer") {
		fmt.Printf("🔍 Creating employer profile for user: %s\n", user.ID)
		// Build location string safely
		location := ""
		if req.City != "" && req.State != "" {
			location = req.City + ", " + req.State
		} else if req.City != "" {
			location = req.City
		} else if req.State != "" {
			location = req.State
		}

		profile := &employerprofile.EmployerProfile{
			UserID:             user.ID,
			CompanyName:        req.CompanyName,
			WebsiteUrl:         req.Website,
			Industry:           req.IndustryType,
			CompanySize:        req.CompanySize,
			CompanyDescription: "", // optional
			RecruiterName:      req.Name,
			Designation:        "Recruiter",
			OfficialEmail:      req.Email,
			PhoneNumber:        phoneNumber, // <-- set here
			LinkedinProfile:    "",
			GstinNumber:        req.GstinNumber,
			CompanyAddress:     req.CompanyAddress,
			City:               req.City,
			State:              req.State,
			Pincode:            req.Pincode,
			JobCategories:      []string{},
			HiringLocations:    []string{location},
			HiringTypes:        []string{"full-time"},
		}
		fmt.Printf("🔍 Employer profile data: %+v\n", profile)
		if err := s.employerRepo.Create(profile); err != nil {
			fmt.Printf("❌ Failed to create employer profile: %v\n", err)
			return nil, "", fmt.Errorf("failed to create employer profile: %w", err)
		}
		fmt.Printf("✅ Employer profile created successfully\n")

	} else if contains(roles, "student") {
		fmt.Printf("🔍 Creating student profile for user: %s\n", user.ID)
		// Build location string safely
		location := ""
		if req.City != "" && req.State != "" {
			location = req.City + ", " + req.State
		} else if req.City != "" {
			location = req.City
		} else if req.State != "" {
			location = req.State
		}

		profile := &studentprofile.StudentProfile{
			UserID:       user.ID,
			Name:         user.Name,
			Email:        user.Email,
			Location:     location,
			PhoneNumber:  phoneNumber, // <-- set here
			Skills:       []string{},
			Resume:       "",
			ProfilePhoto: "",
			Experience:   0.0,
			Education:    "",
			Portfolio:    "",
			Linkedin:     "",
			Github:       "",
		}
		fmt.Printf("🔍 Student profile data: %+v\n", profile)
		if err := s.studentProfileRepo.Create(profile); err != nil {
			fmt.Printf("❌ Failed to create student profile: %v\n", err)
			return nil, "", fmt.Errorf("failed to create user profile: %w", err)
		}
		fmt.Printf("✅ Student profile created successfully\n")
	} else {
		fmt.Printf("⚠️ Unknown role: %s, no profile created\n", req.Role)
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
	fmt.Printf("Password reset requested for %s (user id: %s)\n", user.Email, user.ID)
	return nil
}

// Reset password using token (mock logic)
func (s *authService) ResetPassword(token, newPassword string) error {
	// Password reset is now handled by AAA service
	// This method is kept for backward compatibility but doesn't modify local DB
	return nil
}

// Update profile
func (s *authService) UpdateProfile(userID string, req *UpdateProfileRequest) (*User, error) {
	user, err := s.repo.FindByID(userID)
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
	if err := s.repo.Update(user); err != nil {
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
		profile.GstinNumber = req.GstinNumber
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
		profile.WebsiteUrl = req.Website
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
	if req.ProfilePhoto != "" {
		profile.ProfilePhoto = req.ProfilePhoto
	}
	if req.Resume != "" {
		profile.Resume = req.Resume
	}
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

	return s.studentProfileRepo.Update(profile)
}
