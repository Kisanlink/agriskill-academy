package auth

import (
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
	UpdateProfile(userID string, name string) (*User, error)
}

type authService struct {
	repo UserRepository
}

func NewAuthService(repo UserRepository) AuthService {
	return &authService{repo}
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
	_, err := s.repo.FindByEmail(req.Email)
	if err == nil {
		return nil, "", errors.New("email already registered")
	}
	hashedPW, err := hashPassword(req.Password)
	if err != nil {
		return nil, "", err
	}
	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPW,
		Role:     req.Role,
	}
	// Optionally add employer-specific fields from req
	if err := s.repo.Create(user); err != nil {
		return nil, "", err
	}
	token, err := generateToken(user)
	if err != nil {
		return nil, "", err
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

// Update name
func (s *authService) UpdateProfile(userID string, name string) (*User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	user.Name = name
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}
