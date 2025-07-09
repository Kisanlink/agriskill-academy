package auth

import (
	"asa/config"
	"asa/pkg/jwtutil"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

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
	log.Printf("=== AAA SIGNUP DEBUG START ===")
	log.Printf("📥 Raw request received")

	// Log the raw request body first
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("❌ Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to read request body"})
		return
	}

	// Restore the body for binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	log.Printf("📥 Raw request body: %s", string(bodyBytes))

	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ Signup validation error: %v", err)
		log.Printf("❌ Request body that failed: %s", string(bodyBytes))
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid signup request: " + err.Error()})
		return
	}

	log.Printf("📝 Signup request received:")
	log.Printf("   Name: %s", req.Name)
	log.Printf("   Username: %s", req.Username)
	log.Printf("   Email: %s", req.Email)
	log.Printf("   Role: %s", req.Role)
	log.Printf("   Phone: %d", req.PhoneNumber)
	log.Printf("   Country Code: %s", req.CountryCode)
	log.Printf("   Aadhaar: %s", req.AadhaarNumber)
	log.Printf("   Company: %s", req.CompanyName)
	log.Printf("   GSTIN: %s", req.GstinNumber)

	// Set default country code if not provided
	countryCode := req.CountryCode
	if countryCode == "" {
		countryCode = "+91"
		log.Printf("📝 Using default country code: %s", countryCode)
	}

	// Build AAA payload with correct field names
	aaaPayload := map[string]interface{}{
		"username":      req.Username, // Use username for AAA service
		"password":      req.Password,
		"mobile_number": req.PhoneNumber, // Number, not string
		"country_code":  countryCode,
	}

	// Add optional fields if provided
	if req.AadhaarNumber != "" {
		aaaPayload["aadhaar_number"] = req.AadhaarNumber
		log.Printf("📝 Adding aadhaar number: %s", req.AadhaarNumber)
	}

	log.Printf("📤 AAA Service URL: %s", config.AAAServiceBaseURL)
	log.Printf("📤 AAA Register payload: %+v", aaaPayload)

	body, _ := json.Marshal(aaaPayload)
	log.Printf("📤 AAA Register request body: %s", string(body))

	resp, err := http.Post(config.AAAServiceBaseURL+"/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("❌ AAA Register request failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Auth service unavailable"})
		return
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	log.Printf("📥 AAA Register response status: %d", resp.StatusCode)
	log.Printf("📥 AAA Register response headers: %+v", resp.Header)
	log.Printf("📥 AAA Register response body: %s", string(responseBody))

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		log.Printf("❌ AAA Register failed with status %d", resp.StatusCode)
		c.Data(resp.StatusCode, "application/json", responseBody)
		return
	}

	// Extract tokens from register response headers
	registerToken := resp.Header.Get("token")
	registerRefreshToken := resp.Header.Get("refreshtoken")
	registerUserID := resp.Header.Get("userid")

	log.Printf("🔑 Tokens from register response:")
	log.Printf("   Token: %s", registerToken)
	log.Printf("   RefreshToken: %s", registerRefreshToken)
	log.Printf("   UserID: %s", registerUserID)

	var aaaResp struct {
		Success bool `json:"success"`
		Data    struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"username"`
		} `json:"data"`
	}
	json.Unmarshal(responseBody, &aaaResp)

	log.Printf("📋 Parsed AAA response:")
	log.Printf("   Success: %v", aaaResp.Success)
	log.Printf("   Data.ID: %s", aaaResp.Data.ID)
	log.Printf("   Data.Name: %s", aaaResp.Data.Name)
	log.Printf("   Data.Email: %s", aaaResp.Data.Email)

	// Use header values if available, fallback to body
	userID := registerUserID
	if userID == "" {
		userID = aaaResp.Data.ID
		log.Printf("⚠️ Using UserID from response body: %s", userID)
	} else {
		log.Printf("✅ Using UserID from header: %s", userID)
	}

	assignRole := req.Role
	assignPayload := map[string]interface{}{
		"user_id": userID,
		"role":    assignRole,
	}

	log.Printf("📤 AAA Assign role payload: %+v", assignPayload)

	roleBody, _ := json.Marshal(assignPayload)
	log.Printf("📤 AAA Assign role request body: %s", string(roleBody))

	roleResp, err := http.Post(config.AAAServiceBaseURL+"/assign-role", "application/json", bytes.NewBuffer(roleBody))
	if err != nil {
		log.Printf("❌ AAA Assign role request failed: %v", err)
	} else {
		defer roleResp.Body.Close()
		roleResponseBody, _ := io.ReadAll(roleResp.Body)
		log.Printf("📥 AAA Assign role response status: %d", roleResp.StatusCode)
		log.Printf("📥 AAA Assign role response headers: %+v", roleResp.Header)
		log.Printf("📥 AAA Assign role response body: %s", string(roleResponseBody))
	}

	// Use tokens from register response, fallback to role assignment response
	token := registerToken
	refreshToken := registerRefreshToken

	// If no tokens from register, try role assignment response
	if token == "" && roleResp != nil {
		token = roleResp.Header.Get("token")
		refreshToken = roleResp.Header.Get("refreshtoken")
		log.Printf("🔄 Using tokens from role assignment response")
	}

	log.Printf("🎯 Final tokens for response:")
	log.Printf("   Token: %s", token)
	log.Printf("   RefreshToken: %s", refreshToken)

	// Create local user profile if AAA registration was successful
	var localUser *User
	if userID != "" {
		log.Printf("📝 Creating local user profile...")
		log.Printf("📝 Auth service: %+v", h.authService)
		log.Printf("📝 Request data: %+v", req)

		// Create local user record using auth service
		user, localToken, err := h.authService.Signup(&req)
		if err != nil {
			log.Printf("❌ Failed to create local user: %v", err)
			log.Printf("❌ Error details: %T: %v", err, err)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create local user profile: " + err.Error()})
			return
		}

		localUser = user
		log.Printf("✅ Local user created successfully:")
		log.Printf("   User ID: %s", user.ID)
		log.Printf("   Name: %s", user.Name)
		log.Printf("   Email: %s", user.Email)
		log.Printf("   Local Token: %s", localToken)
	} else {
		log.Printf("⚠️ No AAA user ID received, skipping local user creation")
	}

	log.Printf("✅ Signup successful - sending response")
	log.Printf("=== AAA SIGNUP DEBUG END ===")

	// Use local user data if available, otherwise fallback to request data
	responseUser := gin.H{
		"id":       userID,       // Use AAA user ID as fallback
		"name":     req.Name,     // Use request name as fallback
		"email":    req.Email,    // Use request email as fallback
		"username": req.Username, // Include username for reference
		"role":     assignRole,
	}

	if localUser != nil {
		responseUser["id"] = localUser.ID
		responseUser["name"] = localUser.Name
		responseUser["email"] = localUser.Email
	}

	// Use local token (which includes role) instead of AAA token
	localToken := ""
	if localUser != nil {
		// Generate local token with role included
		localToken, err = jwtutil.GenerateToken(localUser.ID, localUser.Email, localUser.Role, 72*time.Hour)
		if err != nil {
			log.Printf("❌ Failed to generate local token: %v", err)
			// Fallback to AAA token if local token generation fails
			localToken = token
		} else {
			log.Printf("✅ Generated local token with role: %s", localUser.Role)
		}
	} else {
		// Fallback to AAA token if no local user
		localToken = token
		log.Printf("⚠️ Using AAA token as fallback")
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "Signup successful",
		"user":         responseUser,
		"token":        localToken,   // Use local token with role
		"refreshToken": refreshToken, // Keep AAA refresh token
	})
}

// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	log.Printf("=== AAA LOGIN DEBUG START ===")

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ Login validation error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid login request"})
		return
	}

	log.Printf("📝 Login request received:")
	log.Printf("   Username: %s", req.Username)
	log.Printf("   Password: [HIDDEN]")

	aaaPayload := map[string]interface{}{
		"username": req.Username,
		"password": req.Password,
	}

	log.Printf("📤 AAA Service URL: %s", config.AAAServiceBaseURL)
	log.Printf("📤 AAA Login payload: %+v", aaaPayload)

	body, _ := json.Marshal(aaaPayload)
	log.Printf("📤 AAA Login request body: %s", string(body))

	resp, err := http.Post(config.AAAServiceBaseURL+"/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("❌ AAA Login request failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Auth service unavailable"})
		return
	}
	defer resp.Body.Close()
	responseBody, _ := io.ReadAll(resp.Body)

	log.Printf("📥 AAA Login response status: %d", resp.StatusCode)
	log.Printf("📥 AAA Login response headers: %+v", resp.Header)
	log.Printf("📥 AAA Login response body: %s", string(responseBody))

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ AAA Login failed with status %d", resp.StatusCode)
		log.Printf("❌ AAA Login error response: %s", string(responseBody))
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

	log.Printf("🔑 Tokens from login response headers:")
	log.Printf("   Token: %s", token)
	log.Printf("   RefreshToken: %s", refreshToken)
	log.Printf("   UserID: %s", userID)

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

	log.Printf("📋 Parsed AAA login response:")
	log.Printf("   Success: %v", aaaResp.Success)
	log.Printf("   Message: %s", aaaResp.Message)
	log.Printf("   Data.ID: %s", aaaResp.Data.ID)
	log.Printf("   Data.Name: %s", aaaResp.Data.Name)
	log.Printf("   Data.Email: %s", aaaResp.Data.Email)
	log.Printf("   Data.Role: %s", aaaResp.Data.Role)
	log.Printf("   Data.Created: %s", aaaResp.Data.Created)
	log.Printf("   Token: %s", aaaResp.Token)
	log.Printf("   RefreshToken: %s", aaaResp.RefreshToken)

	// Use header values if available, fallback to body
	if userID == "" {
		userID = aaaResp.Data.ID
		log.Printf("⚠️ Using UserID from response body: %s", userID)
	} else {
		log.Printf("✅ Using UserID from header: %s", userID)
	}

	// If no tokens from headers, try response body
	if token == "" {
		token = aaaResp.Token
		log.Printf("⚠️ Using token from response body: %s", token)
	} else {
		log.Printf("✅ Using token from header: %s", token)
	}

	if refreshToken == "" {
		refreshToken = aaaResp.RefreshToken
		log.Printf("⚠️ Using refresh token from response body: %s", refreshToken)
	} else {
		log.Printf("✅ Using refresh token from header: %s", refreshToken)
	}

	log.Printf("🎯 Final login response data:")
	log.Printf("   UserID: %s", userID)
	log.Printf("   Name: %s", aaaResp.Data.Name)
	log.Printf("   Email: %s", aaaResp.Data.Email)
	log.Printf("   Role: %s", aaaResp.Data.Role)
	log.Printf("   Token: %s", token)
	log.Printf("   RefreshToken: %s", refreshToken)

	log.Printf("✅ Login successful - sending response")
	log.Printf("=== AAA LOGIN DEBUG END ===")

	// Generate local token with role for proper permission checks
	localToken := ""
	if aaaResp.Data.Role != "" {
		// Generate local token with role from AAA response
		localToken, err = jwtutil.GenerateToken(userID, aaaResp.Data.Email, aaaResp.Data.Role, 72*time.Hour)
		if err != nil {
			log.Printf("❌ Failed to generate local token: %v", err)
			// Fallback to AAA token if local token generation fails
			localToken = token
		} else {
			log.Printf("✅ Generated local token with role: %s", aaaResp.Data.Role)
		}
	} else {
		// Fallback to AAA token if no role in response
		localToken = token
		log.Printf("⚠️ Using AAA token as fallback (no role in response)")
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "Login successful",
		"id":           userID,
		"name":         aaaResp.Data.Name,
		"email":        aaaResp.Data.Email,
		"role":         aaaResp.Data.Role,
		"created":      aaaResp.Data.Created,
		"token":        localToken, // Use local token with role
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
