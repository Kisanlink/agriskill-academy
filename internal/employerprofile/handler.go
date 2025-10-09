package employerprofile

import (
	"github.com/Kisanlink/agriskill-academy/internal/middleware"
	"github.com/Kisanlink/agriskill-academy/internal/storage"
	"github.com/Kisanlink/agriskill-academy/pkg/authz"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type EmployerProfileHandler struct {
	service EmployerProfileService
	storage storage.StorageService
}

func NewEmployerProfileHandler(s EmployerProfileService, storageSvc storage.StorageService) *EmployerProfileHandler {
	return &EmployerProfileHandler{s, storageSvc}
}

func getJWT(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[7:]
	}
	return ""
}

// IsValidImageFile validates if the file is a valid image type
func IsValidImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	allowedTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, allowedExt := range allowedTypes {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

// getMimeTypeFromExtension returns MIME type based on file extension
func getMimeTypeFromExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

// @Summary Get Employer Profile
// @Description Get a specific employer profile by ID
// @Tags Employer Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param employerId path string true "Employer ID"
// @Success 200 {object} map[string]interface{} "Profile fetched successfully"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 404 {object} map[string]interface{} "Employer profile not found"
// @Router /api/employers/{employerId}/profile [get]
// GET /employers/:employerId/profile
func (h *EmployerProfileHandler) GetProfile(c *gin.Context) {
	username := c.GetString("email")
	employerID := c.Param("employerId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_employer_profiles", "read", employerID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	userID := c.Param("employerId")
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Employer profile not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": profile})
}

// @Summary Get My Employer Profile
// @Description Get the current user's employer profile
// @Tags Employer Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Profile fetched successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 404 {object} map[string]interface{} "Profile not found"
// @Router /api/employers/me/profile [get]
// GET /employers/me/profile
func (h *EmployerProfileHandler) GetMyProfile(c *gin.Context) {
	username := c.GetString("email")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_employer_profiles", "read", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Profile not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": profile})
}

// @Summary Update Employer Profile
// @Description Update a specific employer profile by ID
// @Tags Employer Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param employerId path string true "Employer ID"
// @Param request body UpdateEmployerProfileRequest true "Profile update data"
// @Success 200 {object} map[string]interface{} "Profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Router /api/employers/{employerId}/profile [put]
// PUT /employers/:employerId/profile
func (h *EmployerProfileHandler) UpdateProfile(c *gin.Context) {
	username := c.GetString("email")
	employerID := c.Param("employerId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_employer_profiles", "update", employerID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	userID := c.Param("employerId")
	var req UpdateEmployerProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request format"})
		return
	}

	// Get existing profile first
	existingProfile, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Profile not found"})
		return
	}

	// Update only the fields that are provided in the request (non-empty values)
	if req.CompanyName != "" {
		existingProfile.CompanyName = req.CompanyName
	}
	if req.Industry != "" {
		existingProfile.Industry = req.Industry
	}
	if req.CompanySize != "" {
		existingProfile.CompanySize = req.CompanySize
	}
	if req.WebsiteURL != "" {
		existingProfile.WebsiteURL = req.WebsiteURL
	}
	if req.CompanyDescription != "" {
		existingProfile.CompanyDescription = req.CompanyDescription
	}
	if req.RecruiterName != "" {
		existingProfile.RecruiterName = req.RecruiterName
	}
	if req.Designation != "" {
		existingProfile.Designation = req.Designation
	}
	if req.OfficialEmail != "" {
		existingProfile.OfficialEmail = req.OfficialEmail
	}
	if req.PhoneNumber != "" {
		existingProfile.PhoneNumber = req.PhoneNumber
	}
	if req.LinkedinProfile != "" {
		existingProfile.LinkedinProfile = req.LinkedinProfile
	}
	if req.GSTINNumber != "" {
		existingProfile.GSTINNumber = req.GSTINNumber
	}
	if req.CompanyAddress != "" {
		existingProfile.CompanyAddress = req.CompanyAddress
	}
	if req.City != "" {
		existingProfile.City = req.City
	}
	if req.State != "" {
		existingProfile.State = req.State
	}
	if req.Pincode != "" {
		existingProfile.Pincode = req.Pincode
	}

	// Handle arrays - only update if provided (non-nil)
	if req.JobCategories != nil {
		existingProfile.JobCategories = req.JobCategories
	}
	if req.HiringLocations != nil {
		existingProfile.HiringLocations = req.HiringLocations
	}
	if req.HiringTypes != nil {
		existingProfile.HiringTypes = req.HiringTypes
	}

	// Handle logo fields - only update if provided
	if req.LogoName != "" {
		existingProfile.LogoName = req.LogoName
	}
	if req.LogoType != "" {
		existingProfile.LogoType = req.LogoType
	}
	if req.LogoSize != 0 {
		existingProfile.LogoSize = req.LogoSize
	}
	if req.LogoKey != "" {
		existingProfile.LogoKey = req.LogoKey
	}

	if err := h.service.UpdateProfile(existingProfile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
		"data":    existingProfile,
	})
}

// @Summary Update My Employer Profile
// @Description Update the current user's employer profile
// @Tags Employer Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateEmployerProfileRequest true "Profile update data"
// @Success 200 {object} map[string]interface{} "Profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Router /api/employers/me/profile [put]
// PUT /employers/me/profile
func (h *EmployerProfileHandler) UpdateMyProfile(c *gin.Context) {
	username := c.GetString("email")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)

	// Use user_id from JWT context for permission check instead of URL parameter
	allowed, err := authz.CheckLocalPermission(username, "db_asa_employer_profiles", "update", userID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	var req UpdateEmployerProfileRequest
	contentType := c.GetHeader("Content-Type")

	middleware.DebugLog("🔍 DEBUG: Employer UpdateMyProfile - Content-Type: %s\n", contentType)
	middleware.DebugLog("🔍 DEBUG: User ID: %s\n", userID)

	if strings.Contains(contentType, "multipart/form-data") {
		middleware.DebugLog("🔍 DEBUG: Processing multipart form data for employer\n")

		// Handle logo upload
		if logoFile, err := c.FormFile("logo"); err == nil {
			middleware.DebugLog("🔍 DEBUG: Logo found - Name: %s, Size: %d\n", logoFile.Filename, logoFile.Size)

			// Validate file type
			if !IsValidImageFile(logoFile.Filename) {
				middleware.DebugLog("❌ DEBUG: Invalid logo type: %s\n", logoFile.Filename)
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Invalid image type. Allowed: JPG, PNG, GIF, WebP",
				})
				return
			}

			// Validate file size (5MB max for images)
			if logoFile.Size > 5*1024*1024 {
				middleware.DebugLog("❌ DEBUG: Logo size too large: %d bytes\n", logoFile.Size)
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Image size exceeds maximum allowed size (5MB)",
				})
				return
			}


			// Upload to S3 and set only the S3 key
			key, err := h.storage.SaveImage(logoFile, "employer_logos")
			if err != nil {
				middleware.DebugLog("❌ DEBUG: Failed to upload logo: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Failed to upload logo",
				})
				return
			}
			req.LogoKey = key
			middleware.DebugLog("✅ DEBUG: Logo S3 key set: %s\n", key)
		}

		// Handle other form fields
		if companyName := c.PostForm("company_name"); companyName != "" {
			req.CompanyName = companyName
			middleware.DebugLog("🔍 DEBUG: Company name from form: %s\n", companyName)
		}
		if industry := c.PostForm("industry"); industry != "" {
			req.Industry = industry
			middleware.DebugLog("🔍 DEBUG: Industry from form: %s\n", industry)
		}
		if companySize := c.PostForm("company_size"); companySize != "" {
			req.CompanySize = companySize
			middleware.DebugLog("🔍 DEBUG: Company size from form: %s\n", companySize)
		}
		if websiteURL := c.PostForm("website_url"); websiteURL != "" {
			req.WebsiteURL = websiteURL
			middleware.DebugLog("🔍 DEBUG: Website URL from form: %s\n", websiteURL)
		}
		if companyDescription := c.PostForm("company_description"); companyDescription != "" {
			req.CompanyDescription = companyDescription
			middleware.DebugLog("🔍 DEBUG: Company description from form: %s\n", companyDescription)
		}
		if recruiterName := c.PostForm("recruiter_name"); recruiterName != "" {
			req.RecruiterName = recruiterName
			middleware.DebugLog("🔍 DEBUG: Recruiter name from form: %s\n", recruiterName)
		}
		if designation := c.PostForm("designation"); designation != "" {
			req.Designation = designation
			middleware.DebugLog("🔍 DEBUG: Designation from form: %s\n", designation)
		}
		if officialEmail := c.PostForm("official_email"); officialEmail != "" {
			req.OfficialEmail = officialEmail
			middleware.DebugLog("🔍 DEBUG: Official email from form: %s\n", officialEmail)
		}
		if phoneNumber := c.PostForm("phone_number"); phoneNumber != "" {
			req.PhoneNumber = phoneNumber
			middleware.DebugLog("🔍 DEBUG: Phone number from form: %s\n", phoneNumber)
		}
		if linkedinProfile := c.PostForm("linkedin_profile"); linkedinProfile != "" {
			req.LinkedinProfile = linkedinProfile
			middleware.DebugLog("🔍 DEBUG: LinkedIn profile from form: %s\n", linkedinProfile)
		}
		if gstinNumber := c.PostForm("gstin_number"); gstinNumber != "" {
			req.GSTINNumber = gstinNumber
			middleware.DebugLog("🔍 DEBUG: GSTIN number from form: %s\n", gstinNumber)
		}
		if companyAddress := c.PostForm("company_address"); companyAddress != "" {
			req.CompanyAddress = companyAddress
			middleware.DebugLog("🔍 DEBUG: Company address from form: %s\n", companyAddress)
		}
		if city := c.PostForm("city"); city != "" {
			req.City = city
			middleware.DebugLog("🔍 DEBUG: City from form: %s\n", city)
		}
		if state := c.PostForm("state"); state != "" {
			req.State = state
			middleware.DebugLog("🔍 DEBUG: State from form: %s\n", state)
		}
		if pincode := c.PostForm("pincode"); pincode != "" {
			req.Pincode = pincode
			middleware.DebugLog("🔍 DEBUG: Pincode from form: %s\n", pincode)
		}
	} else {
		// Handle JSON request
		middleware.DebugLog("🔍 DEBUG: Processing JSON request for employer\n")
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.DebugLog("❌ DEBUG: JSON binding error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid request format",
				"error":   err.Error(),
			})
			return
		}
		middleware.DebugLog("✅ DEBUG: JSON request parsed successfully\n")
	}

	// Get existing profile first
	middleware.DebugLog("🔍 DEBUG: Getting employer profile for user ID: %s\n", userID)
	existingProfile, err := h.service.GetProfile(userID)
	if err != nil {
		middleware.DebugLog("❌ DEBUG: Employer profile not found: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Profile not found"})
		return
	}
	middleware.DebugLog("✅ DEBUG: Existing employer profile found\n")

	// Update only the fields that are provided in the request (non-empty values)
	if req.CompanyName != "" {
		existingProfile.CompanyName = req.CompanyName
		middleware.DebugLog("🔍 DEBUG: Updated company name: %s\n", req.CompanyName)
	}
	if req.Industry != "" {
		existingProfile.Industry = req.Industry
		middleware.DebugLog("🔍 DEBUG: Updated industry: %s\n", req.Industry)
	}
	if req.CompanySize != "" {
		existingProfile.CompanySize = req.CompanySize
		middleware.DebugLog("🔍 DEBUG: Updated company size: %s\n", req.CompanySize)
	}
	if req.WebsiteURL != "" {
		existingProfile.WebsiteURL = req.WebsiteURL
		middleware.DebugLog("🔍 DEBUG: Updated website URL: %s\n", req.WebsiteURL)
	}
	if req.CompanyDescription != "" {
		existingProfile.CompanyDescription = req.CompanyDescription
		middleware.DebugLog("🔍 DEBUG: Updated company description\n")
	}
	if req.RecruiterName != "" {
		existingProfile.RecruiterName = req.RecruiterName
		middleware.DebugLog("🔍 DEBUG: Updated recruiter name: %s\n", req.RecruiterName)
	}
	if req.Designation != "" {
		existingProfile.Designation = req.Designation
		middleware.DebugLog("🔍 DEBUG: Updated designation: %s\n", req.Designation)
	}
	if req.OfficialEmail != "" {
		existingProfile.OfficialEmail = req.OfficialEmail
		middleware.DebugLog("🔍 DEBUG: Updated official email: %s\n", req.OfficialEmail)
	}
	if req.PhoneNumber != "" {
		existingProfile.PhoneNumber = req.PhoneNumber
		middleware.DebugLog("🔍 DEBUG: Updated phone number: %s\n", req.PhoneNumber)
	}
	if req.LinkedinProfile != "" {
		existingProfile.LinkedinProfile = req.LinkedinProfile
		middleware.DebugLog("🔍 DEBUG: Updated LinkedIn profile: %s\n", req.LinkedinProfile)
	}
	if req.GSTINNumber != "" {
		existingProfile.GSTINNumber = req.GSTINNumber
		middleware.DebugLog("🔍 DEBUG: Updated GSTIN number: %s\n", req.GSTINNumber)
	}
	if req.CompanyAddress != "" {
		existingProfile.CompanyAddress = req.CompanyAddress
		middleware.DebugLog("🔍 DEBUG: Updated company address: %s\n", req.CompanyAddress)
	}
	if req.City != "" {
		existingProfile.City = req.City
		middleware.DebugLog("🔍 DEBUG: Updated city: %s\n", req.City)
	}
	if req.State != "" {
		existingProfile.State = req.State
		middleware.DebugLog("🔍 DEBUG: Updated state: %s\n", req.State)
	}
	if req.Pincode != "" {
		existingProfile.Pincode = req.Pincode
		middleware.DebugLog("🔍 DEBUG: Updated pincode: %s\n", req.Pincode)
	}

	// Handle arrays - only update if provided (non-nil)
	if req.JobCategories != nil {
		existingProfile.JobCategories = req.JobCategories
		middleware.DebugLog("🔍 DEBUG: Updated job categories: %v\n", req.JobCategories)
	}
	if req.HiringLocations != nil {
		existingProfile.HiringLocations = req.HiringLocations
		middleware.DebugLog("🔍 DEBUG: Updated hiring locations: %v\n", req.HiringLocations)
	}
	if req.HiringTypes != nil {
		existingProfile.HiringTypes = req.HiringTypes
		middleware.DebugLog("🔍 DEBUG: Updated hiring types: %v\n", req.HiringTypes)
	}

	// Update logo key if provided
	if req.LogoKey != "" {
		existingProfile.LogoKey = req.LogoKey
		middleware.DebugLog("🔍 DEBUG: Updated logo - Key: %s\n", req.LogoKey)
	}

	middleware.DebugLog("🔍 DEBUG: Updating employer profile in database\n")
	if err := h.service.UpdateProfile(existingProfile); err != nil {
		middleware.DebugLog("❌ DEBUG: UpdateProfile error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Update failed"})
		return
	}

	middleware.DebugLog("✅ DEBUG: Employer profile updated successfully\n")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
		"data":    existingProfile,
	})
}

// @Summary Create Employer Profile
// @Description Create a new employer profile for the current user
// @Tags Employer Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body EmployerProfile true "Profile creation data"
// @Success 201 {object} map[string]interface{} "Profile created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Router /api/employers/profile [post]
// POST /employers/profile
func (h *EmployerProfileHandler) CreateProfile(c *gin.Context) {
	username := c.GetString("email")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_employer_profiles", "create", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	var req EmployerProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	// ✅ Get user ID from JWT context (set by AuthMiddleware)
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User ID not found in token"})
		return
	}

	// ✅ Inject user ID into the profile model before saving
	req.UserID = userID.(string)

	if err := h.service.CreateProfile(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create profile"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": req})
}

// @Summary Delete Employer Profile
// @Description Delete a specific employer profile by ID
// @Tags Employer Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param employerId path string true "Employer ID"
// @Success 200 {object} map[string]interface{} "Profile deleted successfully"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/employers/{employerId}/profile [delete]
// DELETE /employers/:employerId/profile
func (h *EmployerProfileHandler) DeleteProfile(c *gin.Context) {
	username := c.GetString("email")
	employerID := c.Param("employerId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_employer_profiles", "delete", employerID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	userID := c.Param("employerId")
	if err := h.service.DeleteProfile(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Profile deleted"})
}

// @Summary Upload Employer Logo
// @Description Upload a logo file for the current employer profile
// @Tags Employer Profile
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Logo file (JPG, PNG, GIF, WebP, max 5MB)"
// @Success 200 {object} map[string]interface{} "Logo uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid file type or size"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Router /api/employers/me/logo [post]
// POST /employers/me/logo
func (h *EmployerProfileHandler) UploadLogo(c *gin.Context) {
	username := c.GetString("email")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_employer_profiles", "update", userID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	// Get the file from the request
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Logo file is required",
			"error":   err.Error(),
		})
		return
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	isValid := false
	for _, allowedExt := range allowedTypes {
		if ext == allowedExt {
			isValid = true
			break
		}
	}
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid file type. Allowed: JPG, PNG, GIF, WebP",
		})
		return
	}

	// Validate file size (5MB max)
	if fileHeader.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "File size exceeds maximum allowed size (5MB)",
		})
		return
	}

	// Validate file can be opened (for S3 upload later)
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to read file",
		})
		return
	}
	defer file.Close()

	// Get or create employer profile
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Employer profile not found. Please create a profile first.",
		})
		return
	}

	// Update logo data
	profile.LogoName = fileHeader.Filename
	profile.LogoType = fileHeader.Header.Get("Content-Type")
	if profile.LogoType == "" {
		profile.LogoType = getMimeTypeFromExtension(fileHeader.Filename)
	}
	profile.LogoSize = fileHeader.Size
	// Note: LogoKey will be set by the service when uploading to S3

	// Save profile
	err = h.service.UpdateProfile(profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update profile with logo",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logo uploaded successfully",
		"data": gin.H{
			"logo_name": profile.LogoName,
			"logo_type": profile.LogoType,
			"logo_size": profile.LogoSize,
			"file_url":  fmt.Sprintf("/api/files/serve/logo/%s", userID),
		},
	})
}

// @Summary Upload My Employer Logo
// @Description Upload a logo file to the current user's employer profile
// @Tags Employer Profile
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param logo formData file true "Logo file (JPG, PNG, GIF, WebP, max 5MB)"
// @Success 200 {object} map[string]interface{} "Logo uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid file type or size"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Router /api/employers/me/logo [post]
// POST /employers/me/logo
func (h *EmployerProfileHandler) UploadMyLogo(c *gin.Context) {
	username := c.GetString("email")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_files", "create", "", jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	// Get the file from the request
	fileHeader, err := c.FormFile("logo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Logo file is required",
			"error":   err.Error(),
		})
		return
	}

	// Validate file type
	if !IsValidImageFile(fileHeader.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid file type. Allowed: JPG, PNG, GIF, WebP",
		})
		return
	}

	// Validate file size (5MB max)
	if fileHeader.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "File size exceeds maximum allowed size (5MB)",
		})
		return
	}

	// Upload to S3
	key, err := h.storage.SaveImage(fileHeader, "employer_logos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to upload logo",
		})
		return
	}

	// Get or create employer profile
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Employer profile not found",
		})
		return
	}

	// Update logo key
	profile.LogoKey = key

	// Save profile
	err = h.service.UpdateProfile(profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update profile with logo",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logo uploaded successfully",
		"data": gin.H{
			"file_url": fmt.Sprintf("/api/files/serve/logo/%s", userID),
		},
	})
}

// @Summary Update My Employer Logo
// @Description Update the logo file for the current user's employer profile
// @Tags Employer Profile
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param logo formData file true "Logo file (JPG, PNG, GIF, WebP, max 5MB)"
// @Success 200 {object} map[string]interface{} "Logo updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid file type or size"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Router /api/employers/me/logo [put]
// PUT /employers/me/logo
func (h *EmployerProfileHandler) UpdateMyLogo(c *gin.Context) {
	username := c.GetString("email")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_employer_profiles", "update", userID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	// Get the file from the request
	fileHeader, err := c.FormFile("logo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Logo file is required",
			"error":   err.Error(),
		})
		return
	}

	// Validate file type
	if !IsValidImageFile(fileHeader.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid file type. Allowed: JPG, PNG, GIF, WebP",
		})
		return
	}

	// Validate file size (5MB max)
	if fileHeader.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "File size exceeds maximum allowed size (5MB)",
		})
		return
	}

	// Upload to S3
	key, err := h.storage.SaveImage(fileHeader, "employer_logos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to upload logo",
		})
		return
	}

	// Get existing profile
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Employer profile not found",
		})
		return
	}

	// Update logo key
	profile.LogoKey = key

	// Save profile
	err = h.service.UpdateProfile(profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update profile",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logo updated successfully",
		"data":    profile,
	})
}
