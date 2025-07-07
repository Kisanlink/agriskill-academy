package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(s AuthService) *AuthHandler {
	return &AuthHandler{s}
}

// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid login request"})
		return
	}
	user, token, err := h.service.Login(req.Email, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Login successful", "user": user, "token": token})
}

// POST /auth/signup
func (h *AuthHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid signup request"})
		return
	}

	user, token, err := h.service.Signup(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Signup successful",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
		"token": token,
	})

	fmt.Println("Role:", req.Role)
	fmt.Println("Employer profile creation triggered")

}

// GET /auth/verify
func (h *AuthHandler) Verify(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing token"})
		return
	}
	valid, err := h.service.VerifyToken(token)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid or expired token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Token is valid"})
}

// POST /auth/forgot-password
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid email"})
		return
	}
	if err := h.service.SendResetLink(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Reset link sent"})
}

// POST /auth/reset-password
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}
	if err := h.service.ResetPassword(req.Token, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Password reset successful"})
}

// PUT /auth/profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request: " + err.Error()})
		return
	}

	// Validate that at least one field is provided
	if req.Name == "" && req.Email == "" && req.PhoneNumber == "" && req.Location == "" &&
		req.ProfilePhoto == "" && req.Bio == "" && req.LinkedinProfile == "" && req.Website == "" &&
		len(req.Skills) == 0 && req.CompanyName == "" && req.CompanyDescription == "" &&
		req.Industry == "" && req.CompanySize == "" && req.RecruiterName == "" &&
		req.Designation == "" && req.OfficialEmail == "" && req.GstinNumber == "" &&
		req.CompanyAddress == "" && req.City == "" && req.State == "" && req.Pincode == "" &&
		len(req.JobCategories) == 0 && len(req.HiringLocations) == 0 && len(req.HiringTypes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "At least one field must be provided for update"})
		return
	}

	user, err := h.service.UpdateProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Get updated profile information based on user role
	var profile interface{}
	switch user.Role {
	case "employer":
		employerProfile, err := h.service.(*authService).employerRepo.GetByUserID(userID)
		if err == nil {
			profile = employerProfile
		}
	case "student":
		studentProfile, err := h.service.(*authService).userProfileRepo.GetByUserID(userID)
		if err == nil {
			profile = studentProfile
		}
	}

	response := ProfileResponse{
		Success: true,
		Message: "Profile updated successfully",
		User:    user,
		Profile: profile,
	}

	c.JSON(http.StatusOK, response)
}
