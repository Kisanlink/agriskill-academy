// File: internal/auth/service.go

package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(email, password, role string) (*User, string, error)
	Signup(user *User, password string) (*User, string, error)
}

type authService struct {
	repo UserRepository
}

func NewAuthService(r UserRepository) AuthService {
	return &authService{r}
}

func (s *authService) Login(email, password, role string) (*User, string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, "", errors.New("user not found")
	}
	if user.Role != role {
		return nil, "", errors.New("invalid role")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid password")
	}
	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (s *authService) Signup(user *User, password string) (*User, string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, "", err
	}
	user.Password = string(hashed)
	err = s.repo.Create(user)
	if err != nil {
		return nil, "", err
	}
	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (s *authService) generateToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret" // default for dev
	}
	return token.SignedString([]byte(secret))
}
