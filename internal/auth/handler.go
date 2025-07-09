package auth

import (
	"asa/config"
	"asa/pkg/jwtutil"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// POST /auth/signup
func (h *AuthHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid signup request"})
		return
	}

	// Add phone number to the payload
	aaaPayload := map[string]interface{}{
		"username":    req.Email, // email as username
		"password":    req.Password,
		"name":        req.Name,
		"phoneNumber": req.PhoneNumber, // Add phone number here
	}
	body, _ := json.Marshal(aaaPayload)
	resp, err := http.Post(config.AAAServiceBaseURL+"/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Auth service unavailable"})
		return
	}
	defer resp.Body.Close()
	responseBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		c.Data(resp.StatusCode, "application/json", responseBody)
		return
	}

	// Extract tokens from register response headers
	registerToken := resp.Header.Get("token")
	registerRefreshToken := resp.Header.Get("refreshtoken")
	registerUserID := resp.Header.Get("userid")

	var aaaResp struct {
		Success bool `json:"success"`
		Data    struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"username"`
		} `json:"data"`
	}
	json.Unmarshal(responseBody, &aaaResp)

	// Use header values if available, fallback to body
	userID := registerUserID
	if userID == "" {
		userID = aaaResp.Data.ID
	}

	assignRole := req.Role
	assignPayload := map[string]interface{}{
		"user_id": userID,
		"role":    assignRole,
	}
	roleBody, _ := json.Marshal(assignPayload)
	roleResp, err := http.Post(config.AAAServiceBaseURL+"/assign-role", "application/json", bytes.NewBuffer(roleBody))
	if err != nil {
		log.Printf("Failed to assign role: %v", err)
	}

	// Use tokens from register response, fallback to role assignment response
	token := registerToken
	refreshToken := registerRefreshToken

	// If no tokens from register, try role assignment response
	if token == "" && roleResp != nil {
		token = roleResp.Header.Get("token")
		refreshToken = roleResp.Header.Get("refreshtoken")
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Signup successful",
		"user": gin.H{
			"id":    userID,
			"name":  req.Name,
			"email": req.Email,
			"role":  assignRole,
		},
		"token":        token,
		"refreshToken": refreshToken,
	})
}

// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid login request"})
		return
	}

	aaaPayload := map[string]interface{}{
		"username": req.Email,
		"password": req.Password,
	}
	body, _ := json.Marshal(aaaPayload)
	resp, err := http.Post(config.AAAServiceBaseURL+"/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Auth service unavailable"})
		return
	}
	defer resp.Body.Close()
	responseBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Printf("AAA 400 response: %s", responseBody)
		c.Data(resp.StatusCode, "application/json", responseBody)
		return
	}

	// Extract tokens from response headers (try different case variations)
	token := resp.Header.Get("token")
	if token == "" {
		token = resp.Header.Get("Token")
	}
	if token == "" {
		token = resp.Header.Get("TOKEN")
	}

	refreshToken := resp.Header.Get("refreshtoken")
	if refreshToken == "" {
		refreshToken = resp.Header.Get("RefreshToken")
	}
	if refreshToken == "" {
		refreshToken = resp.Header.Get("REFRESHTOKEN")
	}

	userID := resp.Header.Get("userid")
	if userID == "" {
		userID = resp.Header.Get("UserID")
	}
	if userID == "" {
		userID = resp.Header.Get("USERID")
	}

	// Debug: Log what we're getting from headers (for troubleshooting)
	log.Printf("=== AAA LOGIN DEBUG ===")
	log.Printf("Response Status: %d", resp.StatusCode)
	log.Printf("All Response Headers: %v", resp.Header)
	log.Printf("Token from header: '%s'", token)
	log.Printf("RefreshToken from header: '%s'", refreshToken)
	log.Printf("UserID from header: '%s'", userID)
	log.Printf("Response Body: %s", string(responseBody))
	log.Printf("=== END DEBUG ===")

	var aaaResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Email   string `json:"username"`
			Role    string `json:"role"`
			Created string `json:"created"`
		} `json:"data"`
		Token        string `json:"token"`
		RefreshToken string `json:"refreshToken"`
	}
	json.Unmarshal(responseBody, &aaaResp)

	// Use header values if available, fallback to body
	if userID == "" {
		userID = aaaResp.Data.ID
	}

	// If no tokens from headers, try response body
	if token == "" {
		token = aaaResp.Token
	}
	if refreshToken == "" {
		refreshToken = aaaResp.RefreshToken
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "Login successful",
		"id":           userID,
		"name":         aaaResp.Data.Name,
		"email":        aaaResp.Data.Email,
		"role":         aaaResp.Data.Role,
		"created":      aaaResp.Data.Created,
		"token":        token,
		"refreshToken": refreshToken,
	})
}

// GET /auth/verify
func (h *AuthHandler) Verify(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Missing token"})
		return
	}

	// Manually validate the token locally
	claims, err := jwtutil.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid token"})
		return
	}

	// Proceed with the response
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Token is valid", "data": claims})
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
	body, _ := json.Marshal(req)
	resp, err := http.Post(config.AAAServiceBaseURL+"/forgot-password", "application/json", bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Auth service unavailable"})
		return
	}
	defer resp.Body.Close()
	responseBody, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", responseBody)
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
	body, _ := json.Marshal(req)
	resp, err := http.Post(config.AAAServiceBaseURL+"/reset-password", "application/json", bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Auth service unavailable"})
		return
	}
	defer resp.Body.Close()
	responseBody, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", responseBody)
}

// PUT /auth/profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	roles, exists := c.Get("roles")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "No roles found in token"})
		return
	}

	rolesSlice, ok := roles.([]string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid roles format in token"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request: " + err.Error()})
		return
	}

	// Use the auth service to update the profile
	user, err := h.authService.UpdateProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to update profile: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
		"user":    user,
		"roles":   rolesSlice,
	})
}
