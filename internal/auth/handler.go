package auth

import (
	"asa/config"
	"asa/pkg/jwtutil"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// @Summary User Registration
// @Description Register a new user (student or employer) with the AAA service and create local profile
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body SignupRequest true "User registration data"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 500 {object} map[string]interface{} "Internal server error or AAA service unavailable"
// @Router /api/auth/signup [post]
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

	// Convert phone number from string to number for AAA service
	phoneNumber, err := strconv.ParseInt(req.PhoneNumber, 10, 64)
	if err != nil {
		log.Printf("❌ Invalid phone number format: %s", req.PhoneNumber)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid phone number format"})
		return
	}

	// Prepare payload for AAA service
	aaaPayload := map[string]interface{}{
		"username":      req.Username,
		"password":      req.Password,
		"mobile_number": phoneNumber,
		"country_code":  countryCode,
	}

	// Add aadhaar number if provided, otherwise send null
	if req.AadhaarNumber != "" {
		aaaPayload["aadhaar_number"] = req.AadhaarNumber
		log.Printf("📝 Adding aadhaar number: %s", req.AadhaarNumber)
	} else {
		aaaPayload["aadhaar_number"] = nil
		log.Printf("📝 No aadhaar number provided, sending null")
	}

	log.Printf("📤 AAA Service URL: %s", config.AAAServiceBaseURL)
	log.Printf("📤 AAA Register payload: %+v", aaaPayload)

	body, _ := json.Marshal(aaaPayload)
	log.Printf("📤 AAA Register request body: %s", string(body))

	resp, err := http.Post(config.AAAServiceBaseURL+"/api/v1/register", "application/json", bytes.NewBuffer(body))
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

	// === UPDATED 409 HANDLING ===
	if resp.StatusCode == http.StatusConflict {
		log.Printf("🔄 HTTP 409 Conflict: AAA user exists → lookup via gRPC by phone")

		grpcClient := h.authService.(*authService).grpcClient
		if grpcClient == nil {
			log.Printf("❌ gRPC client not initialized")
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "gRPC client not available"})
			return
		}

		ctx := c.Request.Context()
		respGet, err := grpcClient.GetUserByMobileNumber(ctx, req.PhoneNumber)
		if err != nil {
			log.Printf("❌ gRPC GetUserByMobileNumber failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch existing user from AAA"})
			return
		}
		if respGet == nil || respGet.Data == nil {
			log.Printf("❌ AAA gRPC GetUserByMobileNumber returned no Data")
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Internal error: missing AAA user data"})
			return
		}

		userID := respGet.Data.Id
		phoneNumberStr := strconv.FormatUint(respGet.Data.MobileNumber, 10)

		user, token, err := h.authService.SignupWithID(&req, userID, phoneNumberStr)
		if err != nil {
			log.Printf("❌ Failed to create local user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to link existing user: " + err.Error()})
			return
		}

		log.Printf("✅ Local user linked successfully: %+v", user)
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "User linked successfully",
			"user": gin.H{
				"id":       user.ID,
				"name":     user.Name,
				"email":    user.Email,
				"role":     user.Role,
				"username": req.Username,
			},
			"token":                token,
			"aaa_user_id":          userID,
			"is_existing_aaa_user": true,
		})
		return
	}
	// === END UPDATED 409 HANDLING ===

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
			ID           string `json:"id"`
			Name         string `json:"name"`
			Email        string `json:"username"`
			MobileNumber int64  `json:"mobile_number"`
			CountryCode  string `json:"country_code"`
		} `json:"data"`
		Error []string `json:"error"`
	}
	json.Unmarshal(responseBody, &aaaResp)

	log.Printf("📋 Parsed AAA response:")
	log.Printf("   Success: %v", aaaResp.Success)
	log.Printf("   Data.ID: %s", aaaResp.Data.ID)
	log.Printf("   Data.Name: %s", aaaResp.Data.Name)
	log.Printf("   Data.Email: %s", aaaResp.Data.Email)
	log.Printf("   Data.MobileNumber: %d", aaaResp.Data.MobileNumber)
	log.Printf("   Data.CountryCode: %s", aaaResp.Data.CountryCode)
	log.Printf("   Error: %v", aaaResp.Error)

	// Existing user second‐level conflict handling (now unreachable due to early return above)
	var userID string
	if resp.StatusCode == http.StatusConflict {
		// ... old conflict logic omitted ...
	} else {
		userID = registerUserID
		if userID == "" {
			userID = aaaResp.Data.ID
			log.Printf("⚠️ Using UserID from response body: %s", userID)
		} else {
			log.Printf("✅ Using UserID from header: %s", userID)
		}
	}

	phoneNumberStr := strconv.FormatInt(int64(aaaResp.Data.MobileNumber), 10)
	if phoneNumberStr == "0" {
		phoneNumberStr = req.PhoneNumber
	}

	log.Printf("📝 Creating local user with AAA user ID: %s", userID)
	log.Printf("📝 Phone number for local storage: %s", phoneNumberStr)

	user, token, err := h.authService.SignupWithID(&req, userID, phoneNumberStr)
	if err != nil {
		log.Printf("❌ Failed to create local user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create local user: " + err.Error()})
		return
	}

	log.Printf("✅ Local user created successfully: %+v", user)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"name":     user.Name,
			"email":    user.Email,
			"role":     user.Role,
			"username": req.Username,
		},
		"token":                token,
		"aaa_user_id":          userID,
		"is_existing_aaa_user": resp.StatusCode == http.StatusConflict,
	})
}

// @Summary User Login
// @Description Authenticate user with AAA service and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]interface{} "Invalid credentials"
// @Failure 401 {object} map[string]interface{} "Authentication failed"
// @Failure 500 {object} map[string]interface{} "Internal server error or AAA service unavailable"
// @Router /api/auth/login [post]
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

	resp, err := http.Post(config.AAAServiceBaseURL+"/api/v1/login", "application/json", bytes.NewBuffer(body))
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
			Roles   []struct {
				RoleName string `json:"role_name"`
			} `json:"roles"`
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
	log.Printf("   Data.Roles count: %d", len(aaaResp.Data.Roles))
	for i, role := range aaaResp.Data.Roles {
		log.Printf("   Data.Roles[%d]: %s", i, role.RoleName)
	}
	log.Printf("   Token: %s", aaaResp.Token)
	log.Printf("   RefreshToken: %s", aaaResp.RefreshToken)

	// Extract roles from AAA response
	var roles []string
	if len(aaaResp.Data.Roles) > 0 {
		for _, role := range aaaResp.Data.Roles {
			roles = append(roles, role.RoleName)
		}
		log.Printf("✅ Extracted roles from AAA response: %v", roles)
	} else if aaaResp.Data.Role != "" {
		// Fallback to single role if roles array is empty
		roles = []string{aaaResp.Data.Role}
		log.Printf("⚠️ Using single role as fallback: %v", roles)
	} else {
		log.Printf("⚠️ No roles found in AAA response")
	}

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

	// Try to get additional user details from local database
	// Since we can't access repo directly from interface, we'll use the email from AAA response
	// and assume the user exists in local DB (since they logged in successfully)
	var userName, userEmail string
	var userCreated time.Time

	log.Printf("🔍 === LOCAL USER LOOKUP DEBUG ===")
	log.Printf("🔍 User ID from AAA: %s", userID)
	log.Printf("🔍 AAA Email (username): %s", aaaResp.Data.Email)

	// Debug: List all users in database to see what's available
	log.Printf("🔍 === DATABASE DEBUG ===")
	allUsers, err := h.authService.ListAllUsers()
	if err != nil {
		log.Printf("❌ Failed to list users: %v", err)
	} else {
		log.Printf("✅ Database contains %d users", len(allUsers))
	}
	log.Printf("🔍 === DATABASE DEBUG COMPLETE ===")

	// Try to get local user data using the user ID from AAA response
	// The AAA response email field is actually the username, not the real email
	localUser, err := h.authService.GetUserByID(userID)
	if err != nil {
		log.Printf("❌ Local user lookup failed: %v", err)
		log.Printf("❌ Error type: %T", err)
	} else if localUser != nil {
		userName = localUser.Name
		userEmail = localUser.Email
		userCreated = localUser.CreatedAt
		log.Printf("✅ Found local user data:")
		log.Printf("   ID: %s", localUser.ID)
		log.Printf("   Name: %s", userName)
		log.Printf("   Email: %s", userEmail)
		log.Printf("   Created: %s", userCreated.Format("2006-01-02T15:04:05Z"))
		log.Printf("   Role: %s", localUser.Role)
	} else {
		log.Printf("⚠️ Local user lookup returned nil user")
	}

	if err != nil || localUser == nil {
		log.Printf("⚠️ No local user data found for user ID: %s", userID)
		log.Printf("⚠️ Falling back to AAA response data")
		// Fallback to AAA response data
		userName = aaaResp.Data.Name
		userEmail = aaaResp.Data.Email // This is actually the username
		log.Printf("⚠️ Fallback values - Name: %s, Email: %s", userName, userEmail)
	}

	log.Printf("🔍 === LOCAL USER LOOKUP COMPLETE ===")

	// Generate local token with roles for proper permission checks
	localToken := ""
	primaryRole := ""
	if len(roles) > 0 {
		// Generate local token with roles from AAA response
		// Use the first role as primary role for backward compatibility
		primaryRole = roles[0]
		localToken, err = jwtutil.GenerateToken(userID, userEmail, primaryRole, 72*time.Hour)
		if err != nil {
			log.Printf("❌ Failed to generate local token: %v", err)
			// Fallback to AAA token if local token generation fails
			localToken = token
		} else {
			log.Printf("✅ Generated local token with primary role: %s (all roles: %v)", primaryRole, roles)
		}
	} else {
		// Fallback to AAA token if no roles in response
		localToken = token
		log.Printf("⚠️ Using AAA token as fallback (no roles in response)")
	}

	// Use local user data if available, otherwise fallback to AAA response
	responseName := userName
	responseEmail := userEmail

	// If we have local user data, use it; otherwise fallback to request data
	if localUser != nil {
		responseName = localUser.Name
		responseEmail = localUser.Email
		log.Printf("✅ Using local user data - Name: %s, Email: %s", responseName, responseEmail)
	} else {
		log.Printf("⚠️ No local user data, using fallback - Name: %s, Email: %s", responseName, responseEmail)
	}

	// Handle created date
	var responseCreated string
	if !userCreated.IsZero() {
		responseCreated = userCreated.Format("2006-01-02T15:04:05Z")
	} else if aaaResp.Data.Created != "" {
		responseCreated = aaaResp.Data.Created
	} else {
		// Use current time as fallback for created date
		responseCreated = time.Now().Format("2006-01-02T15:04:05Z")
		log.Printf("⚠️ No created date available, using current time: %s", responseCreated)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"message":       "Login successful",
		"id":            userID,
		"name":          responseName,
		"email":         responseEmail,
		"role":          primaryRole, // Use extracted primary role instead of empty Data.Role
		"created":       responseCreated,
		"token":         localToken, // Use local token with role
		"refresh_token": refreshToken,
	})
}

// @Summary Verify JWT Token
// @Description Verify the validity of a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Token is valid"
// @Failure 401 {object} map[string]interface{} "Invalid or missing token"
// @Router /api/auth/verify [get]
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

// @Summary Forgot Password
// @Description Send password reset email to user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Email for password reset"
// @Success 200 {object} map[string]interface{} "Password reset email sent"
// @Failure 400 {object} map[string]interface{} "Invalid email"
// @Failure 500 {object} map[string]interface{} "Internal server error or AAA service unavailable"
// @Router /api/auth/forgot-password [post]
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

// @Summary Reset Password
// @Description Reset user password using reset token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Reset password data"
// @Success 200 {object} map[string]interface{} "Password reset successful"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error or AAA service unavailable"
// @Router /api/auth/reset-password [post]
// POST /auth/reset-password
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
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

// @Summary Update User Profile
// @Description Update user profile information
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "Profile update data"
// @Success 200 {object} map[string]interface{} "Profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/auth/profile [put]
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

// gRPC handlers removed - using HTTP AAA service integration instead
