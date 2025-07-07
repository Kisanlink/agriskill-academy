package auth

import (
	"AGRIJOBS/internal/employerprofile"
	"AGRIJOBS/internal/userprofile"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(email, password, role string) (*User, string, error)
	Signup(req *SignupRequest) (*User, string, error)
	VerifyToken(tokenStr string) (bool, error)
	SendResetLink(email string) error
	ResetPassword(token, newPassword string) error
	UpdateProfile(userID string, req *UpdateProfileRequest) (*User, error)
}

type authService struct {
	repo            UserRepository
	employerRepo    employerprofile.EmployerProfileRepository
	userProfileRepo userprofile.UserProfileRepository
}

func NewAuthService(repo UserRepository, employerRepo employerprofile.EmployerProfileRepository, userProfileRepo userprofile.UserProfileRepository) AuthService {
	return &authService{
		repo:            repo,
		employerRepo:    employerRepo,
		userProfileRepo: userProfileRepo,
	}
}

// JWT Secret
func getSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "secret"
	}
	return secret
}

// Sign JWT Token
func generateToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
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

// Hash password
func hashPassword(pw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(bytes), err
}

// Compare hashed password
func checkPassword(hashed, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pw))
}

// Login
func (s *authService) Login(email, password, role string) (*User, string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil || user.Role != role {
		return nil, "", errors.New("invalid credentials")
	}
	if err := checkPassword(user.Password, password); err != nil {
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
	// 1. Validate password confirmation
	if req.Password != req.ConfirmPassword {
		return nil, "", errors.New("passwords do not match")
	}

	// 2. Check if user exists
	existingUser, _ := s.repo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, "", errors.New("email already registered")
	}

	// 3. Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, "", err
	}

	// 4. Create user
	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
	}
	err = s.repo.Create(user)
	if err != nil {
		return nil, "", err
	}

	// 5. Generate token
	token, err := generateToken(user)
	if err != nil {
		return nil, "", err
	}

	// 6. Create corresponding profile
	switch user.Role {
	case "employer":
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
		if err := s.employerRepo.Create(profile); err != nil {
			return nil, "", fmt.Errorf("failed to create employer profile: %w", err)
		}

	case "student":
		// Build location string safely
		location := ""
		if req.City != "" && req.State != "" {
			location = req.City + ", " + req.State
		} else if req.City != "" {
			location = req.City
		} else if req.State != "" {
			location = req.State
		}

		profile := &userprofile.UserProfile{
			UserID:       user.ID,
			Name:         user.Name,
			Email:        user.Email,
			Location:     location,
			Skills:       []string{},
			Resume:       "",
			ProfilePhoto: "",
		}
		if err := s.userProfileRepo.Create(profile); err != nil {
			return nil, "", fmt.Errorf("failed to create user profile: %w", err)
		}
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
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(getSecret()), nil
	})
	if err != nil || !parsed.Valid {
		return errors.New("invalid or expired token")
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}
	email, ok := claims["email"].(string)
	if !ok {
		return errors.New("invalid token email")
	}
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return err
	}
	hashed, err := hashPassword(newPassword)
	if err != nil {
		return err
	}
	user.Password = hashed
	return s.repo.Update(user)
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

	// Update role-specific profile based on user role
	switch user.Role {
	case "employer":
		if err := s.updateEmployerProfile(userID, req); err != nil {
			return nil, fmt.Errorf("failed to update employer profile: %w", err)
		}
	case "student":
		if err := s.updateStudentProfile(userID, req); err != nil {
			return nil, fmt.Errorf("failed to update student profile: %w", err)
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
	profile, err := s.userProfileRepo.GetByUserID(userID)
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
	if req.PhoneNumber != "" {
		// Note: UserProfile doesn't have phone number field, we might need to add it
		// For now, we'll skip this field
	}
	if req.ProfilePhoto != "" {
		profile.ProfilePhoto = req.ProfilePhoto
	}
	if len(req.Skills) > 0 {
		profile.Skills = req.Skills
	}
	if req.LinkedinProfile != "" {
		// Note: UserProfile doesn't have linkedin field, we might need to add it
		// For now, we'll skip this field
	}

	return s.userProfileRepo.Update(profile)
}
