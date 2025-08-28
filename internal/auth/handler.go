package auth

import (
	"asa/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// @Summary User Registration
// @Description Register a new user (student or employer) with local authentication
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body SignupRequest true "User registration data"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/auth/signup [post]
// POST /auth/signup
func (h *AuthHandler) Signup(c *gin.Context) {
	middleware.DebugLog("=== LOCAL SIGNUP DEBUG START ===")

	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.DebugLog("❌ Signup validation error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid signup request: " + err.Error()})
		return
	}

	middleware.DebugLog("📝 Signup request received:")
	middleware.DebugLog("   Name: %s", req.Name)
	middleware.DebugLog("   Username: %s", req.Username)
	middleware.DebugLog("   Email: %s", req.Email)
	middleware.DebugLog("   Role: %s", req.Role)

	// Call local auth service
	user, token, err := h.authService.Signup(&req)
	if err != nil {
		middleware.DebugLog("❌ Signup failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	middleware.DebugLog("✅ Signup successful for user: %s", user.ID)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User registered successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
		"token": token,
	})
}

// @Summary User Login
// @Description Login with email and password using local authentication
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]interface{} "Invalid credentials"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/auth/login [post]
// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	middleware.DebugLog("=== LOCAL LOGIN DEBUG START ===")

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.DebugLog("❌ Login validation error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid login request: " + err.Error()})
		return
	}

	middleware.DebugLog("📝 Login request received:")
	middleware.DebugLog("   Username/Email: %s", req.Username)

	// Call local auth service
	user, token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		middleware.DebugLog("❌ Login failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid credentials"})
		return
	}

	middleware.DebugLog("✅ Login successful for user: %s", user.ID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
		"token": token,
	})
}

// @Summary Forgot Password
// @Description Send password reset link (mock implementation)
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Email for password reset"
// @Success 200 {object} map[string]interface{} "Reset link sent"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /api/auth/forgot-password [post]
// POST /auth/forgot-password
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	middleware.DebugLog("=== FORGOT PASSWORD DEBUG START ===")

	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.DebugLog("❌ Forgot password validation error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request: " + err.Error()})
		return
	}

	middleware.DebugLog("📝 Forgot password request received:")
	middleware.DebugLog("   Email: %s", req.Email)

	// Call local auth service
	err := h.authService.SendResetLink(req.Email)
	if err != nil {
		middleware.DebugLog("❌ Forgot password failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	middleware.DebugLog("✅ Forgot password successful")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password reset link sent to your email",
	})
}

// @Summary Reset Password
// @Description Reset password using token (mock implementation)
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Password reset data"
// @Success 200 {object} map[string]interface{} "Password reset successful"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /api/auth/reset-password [post]
// POST /auth/reset-password
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	middleware.DebugLog("=== RESET PASSWORD DEBUG START ===")

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.DebugLog("❌ Reset password validation error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request: " + err.Error()})
		return
	}

	middleware.DebugLog("📝 Reset password request received")

	// Call local auth service
	err := h.authService.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		middleware.DebugLog("❌ Reset password failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	middleware.DebugLog("✅ Reset password successful")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password reset successful",
	})
}

// @Summary Get User Profile
// @Description Get current user's profile information
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "User profile retrieved"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/auth/profile [get]
// GET /auth/profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile retrieved successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// @Summary Update User Profile
// @Description Update current user's profile information
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "Profile update data"
// @Success 200 {object} map[string]interface{} "Profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/auth/profile [put]
// PUT /auth/profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request: " + err.Error()})
		return
	}

	user, err := h.authService.UpdateProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}
