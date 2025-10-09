// File: internal/studentprofile/handler.go

package studentprofile

import (
	"github.com/Kisanlink/agriskill-academy/internal/middleware"
	"github.com/Kisanlink/agriskill-academy/pkg/authz"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Kisanlink/agriskill-academy/internal/storage"

	"github.com/gin-gonic/gin"
)

type StudentProfileHandler struct {
	service StudentProfileService
	storage storage.StorageService
}

func NewStudentProfileHandler(s StudentProfileService, storageSvc storage.StorageService) *StudentProfileHandler {
	return &StudentProfileHandler{s, storageSvc}
}

func getJWT(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[7:]
	}
	return ""
}

// @Summary Get Student Profile
// @Description Get a specific student profile by ID
// @Tags Student Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param studentId path string true "Student ID"
// @Success 200 {object} map[string]interface{} "Profile fetched successfully"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 404 {object} map[string]interface{} "Profile not found"
// @Router /api/students/{studentId}/profile [get]
// GET /students/:studentId/profile
func (h *StudentProfileHandler) GetProfile(c *gin.Context) {
	username := c.GetString("username")
	profileID := c.Param("studentId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_student_profile", "read", profileID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	profile, err := h.service.GetProfile(c.Request.Context(), profileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Profile fetched", "data": profile})
}

// @Summary Update Student Profile
// @Description Update a specific student profile by ID
// @Tags Student Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param studentId path string true "Student ID"
// @Param request body StudentProfile true "Profile update data"
// @Success 200 {object} map[string]interface{} "Profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Router /api/students/{studentId}/profile [put]
// PUT /students/:studentId/profile
func (h *StudentProfileHandler) UpdateProfile(c *gin.Context) {
	username := c.GetString("username")
	profileID := c.Param("studentId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_student_profile", "update", profileID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	userID := c.GetString("user_id")
	if profileID != userID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "You can only update your own profile"})
		return
	}

	var req StudentProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}
	req.UserID = userID
	err = h.service.UpdateProfile(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Profile updated", "data": req})
}

// @Summary Get My Student Profile
// @Description Get the current user's student profile
// @Tags Student Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Profile fetched successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 404 {object} map[string]interface{} "Profile not found"
// @Router /api/students/me/profile [get]
// GET /students/me/profile
func (h *StudentProfileHandler) GetMyProfile(c *gin.Context) {
	username := c.GetString("username")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_student_profile", "read", userID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	profile, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Profile fetched", "data": profile})
}

// @Summary Add Certificate to Student Profile
// @Description Add a certificate to a specific student profile
// @Tags Student Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param studentId path string true "Student ID"
// @Param request body Certificate true "Certificate data"
// @Success 200 {object} map[string]interface{} "Certificate added successfully"
// @Failure 400 {object} map[string]interface{} "Invalid certificate data"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 404 {object} map[string]interface{} "Student profile not found"
// @Router /api/students/{studentId}/certificates [post]
// POST /students/:studentId/certificates
func (h *StudentProfileHandler) AddCertificate(c *gin.Context) {
	username := c.GetString("username")
	studentID := c.Param("studentId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_certificates", "create", studentID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	userID := c.GetString("user_id")
	if studentID != userID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "You can only add certificates to your own profile"})
		return
	}

	var cert Certificate
	if err := c.ShouldBindJSON(&cert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid certificate"})
		return
	}

	// Get the student profile first to get the profile ID
	studentProfile, err := h.service.GetProfile(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Student profile not found"})
		return
	}

	// Set the student profile ID from the actual profile and clear any invalid ID
	cert.StudentProfileID = studentProfile.ID
	cert.ID = "" // Clear ID to let database generate proper UUID

	err = h.service.AddCertificate(&cert)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to add certificate"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Certificate added", "data": cert})
}

// @Summary Delete My Certificate
// @Description Delete a certificate from the current user's profile
// @Tags Student Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param certificateId path string true "Certificate ID"
// @Success 200 {object} map[string]interface{} "Certificate deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 404 {object} map[string]interface{} "Certificate not found"
// @Router /api/students/me/certificates/{certificateId} [delete]
// DELETE /students/me/certificates/:certificateId
func (h *StudentProfileHandler) DeleteMyCertificate(c *gin.Context) {
	username := c.GetString("username")
	userID := c.GetString("user_id")
	certificateID := c.Param("certificateId")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_certificates", "delete", certificateID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	err = h.service.DeleteCertificate(certificateID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Certificate not found or you don't have permission to delete it"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Certificate deleted successfully"})
}

// @Summary Update My Student Profile
// @Description Update the current user's student profile with optional file uploads
// @Tags Student Profile
// @Accept multipart/form-data
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name formData string false "Student name"
// @Param email formData string false "Student email"
// @Param location formData string false "Student location"
// @Param phone_number formData string false "Phone number"
// @Param education formData string false "Education details"
// @Param portfolio formData string false "Portfolio URL"
// @Param linkedin formData string false "LinkedIn profile"
// @Param github formData string false "GitHub profile"
// @Param resume formData file false "Resume file (PDF, DOC, DOCX, max 10MB)"
// @Param profile_photo formData file false "Profile photo (JPG, PNG, GIF, WebP, max 5MB)"
// @Success 200 {object} map[string]interface{} "Profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or file format"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Router /api/students/me/profile [put]
// PUT /students/me/profile
func (h *StudentProfileHandler) UpdateMyProfile(c *gin.Context) {
	middleware.DebugLog("🔍 DEBUG: UpdateMyProfile called\n")

	username := c.GetString("username")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)

	middleware.DebugLog("🔍 DEBUG: Username: %s, UserID: %s\n", username, userID)

	allowed, err := authz.CheckLocalPermission(username, "db_asa_student_profile", "update", userID, jwtToken)
	if err != nil || !allowed {
		middleware.DebugLog("❌ DEBUG: Permission denied - Error: %v, Allowed: %v\n", err, allowed)
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}
	middleware.DebugLog("✅ DEBUG: Permission check passed\n")

	if userID == "" {
		middleware.DebugLog("❌ DEBUG: UserID is empty\n")
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	var req UpdateStudentProfileRequest
	contentType := c.GetHeader("Content-Type")

	middleware.DebugLog("🔍 DEBUG: Content-Type: %s\n", contentType)
	middleware.DebugLog("🔍 DEBUG: User ID: %s\n", userID)

	if strings.Contains(contentType, "multipart/form-data") {
		middleware.DebugLog("🔍 DEBUG: Processing multipart form data\n")

		// Handle profile photo upload
		if profilePhotoFile, err := c.FormFile("profile_photo"); err == nil {
			middleware.DebugLog("🔍 DEBUG: Profile photo found - Name: %s, Size: %d\n", profilePhotoFile.Filename, profilePhotoFile.Size)
			// Validate file type
			if !IsValidImageFile(profilePhotoFile.Filename) {
				middleware.DebugLog("❌ DEBUG: Invalid image type: %s\n", profilePhotoFile.Filename)
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Invalid image type. Allowed: JPG, PNG, GIF, WebP",
				})
				return
			}

			// Validate file size (5MB max for images)
			if profilePhotoFile.Size > 5*1024*1024 {
				middleware.DebugLog("❌ DEBUG: Image size too large: %d bytes\n", profilePhotoFile.Size)
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Image size exceeds maximum allowed size (5MB)",
				})
				return
			}

			// Upload to S3 and set only the S3 key
			key, err := h.storage.SaveImage(profilePhotoFile, "profile_photos")
			if err != nil {
				middleware.DebugLog("❌ DEBUG: Failed to upload profile photo: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Failed to upload profile photo",
				})
				return
			}
			req.ProfilePhotoKey = key
			middleware.DebugLog("✅ DEBUG: Profile photo S3 key set: %s\n", key)
		}

		// Handle resume upload (optional)
		if resumeFile, err := c.FormFile("resume"); err == nil {
			middleware.DebugLog("🔍 DEBUG: Resume found - Name: %s, Size: %d\n", resumeFile.Filename, resumeFile.Size)
			// Validate file type
			if !IsValidResumeFile(resumeFile.Filename) {
				middleware.DebugLog("❌ DEBUG: Invalid resume type: %s\n", resumeFile.Filename)
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Invalid resume type. Allowed: PDF, DOC, DOCX",
				})
				return
			}

			// Validate file size (10MB max for resumes)
			if resumeFile.Size > 10*1024*1024 {
				middleware.DebugLog("❌ DEBUG: Resume size too large: %d bytes\n", resumeFile.Size)
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Resume size exceeds maximum allowed size (10MB)",
				})
				return
			}

			// Upload to S3 and set only the S3 key
			key, err := h.storage.SaveResume(resumeFile, "resumes")
			if err != nil {
				middleware.DebugLog("❌ DEBUG: Failed to upload resume: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Failed to upload resume",
				})
				return
			}
			req.ResumeKey = key
			middleware.DebugLog("✅ DEBUG: Resume S3 key set: %s\n", key)
		}

		// Handle other form fields
		if name := c.PostForm("name"); name != "" {
			req.Name = name
			middleware.DebugLog("🔍 DEBUG: Name from form: %s\n", name)
		}
		if email := c.PostForm("email"); email != "" {
			req.Email = email
			middleware.DebugLog("🔍 DEBUG: Email from form: %s\n", email)
		}
		if location := c.PostForm("location"); location != "" {
			req.Location = location
			middleware.DebugLog("🔍 DEBUG: Location from form: %s\n", location)
		}
		if phoneNumber := c.PostForm("phone_number"); phoneNumber != "" {
			req.PhoneNumber = phoneNumber
			middleware.DebugLog("🔍 DEBUG: Phone number from form: %s\n", phoneNumber)
		}
		if education := c.PostForm("education"); education != "" {
			req.Education = education
			middleware.DebugLog("🔍 DEBUG: Education from form: %s\n", education)
		}
		if portfolio := c.PostForm("portfolio"); portfolio != "" {
			req.Portfolio = portfolio
			middleware.DebugLog("🔍 DEBUG: Portfolio from form: %s\n", portfolio)
		}
		if linkedin := c.PostForm("linkedin"); linkedin != "" {
			req.Linkedin = linkedin
			middleware.DebugLog("🔍 DEBUG: LinkedIn from form: %s\n", linkedin)
		}
		if github := c.PostForm("github"); github != "" {
			req.Github = github
			middleware.DebugLog("🔍 DEBUG: GitHub from form: %s\n", github)
		}

		// Handle skills array from form data
		if skillsStr := c.PostForm("skills"); skillsStr != "" {
			middleware.DebugLog("🔍 DEBUG: Raw skills string from form: %s\n", skillsStr)
			// Parse skills as JSON array from form data
			var skills []string
			if err := json.Unmarshal([]byte(skillsStr), &skills); err == nil {
				req.Skills = PostgreSQLTextArray(skills)
				middleware.DebugLog("🔍 DEBUG: Skills parsed from form: %v\n", skills)
			} else {
				middleware.DebugLog("❌ DEBUG: Failed to parse skills JSON: %v\n", err)
			}
		} else {
			middleware.DebugLog("🔍 DEBUG: No skills field found in form data\n")
		}

		// Handle experience from form data
		if experienceStr := c.PostForm("experience"); experienceStr != "" {
			if exp, err := strconv.ParseFloat(experienceStr, 64); err == nil {
				req.Experience = &exp
				middleware.DebugLog("🔍 DEBUG: Experience from form: %f\n", exp)
			} else {
				middleware.DebugLog("❌ DEBUG: Failed to parse experience: %v\n", err)
			}
		}
	} else {
		// Handle JSON request
		middleware.DebugLog("🔍 DEBUG: Processing JSON request\n")
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

	// Set user ID
	req.UserID = userID

	middleware.DebugLog("🔍 DEBUG: Final request metadata - Name: %s, Email: %s, Skills: %v, SkillsLength: %d, CertificatesCount: %d\n", req.Name, req.Email, req.Skills, len(req.Skills), len(req.Certificates))

	// Get existing profile or create new one
	middleware.DebugLog("🔍 DEBUG: Getting profile for user ID: %s\n", userID)
	profile, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		middleware.DebugLog("🔍 DEBUG: Profile not found, creating new one for user: %s\n", userID)
		// Create new profile
		profile = &StudentProfile{
			UserID: userID,
			Name:   req.Name,
			Email:  req.Email,
		}
		middleware.DebugLog("🔍 DEBUG: New profile created: %+v\n", profile)
	} else {
		middleware.DebugLog("✅ DEBUG: Existing profile found: %+v\n", profile)
	}

	// Update profile fields from request
	middleware.DebugLog("🔍 DEBUG: Updating profile fields\n")
	if req.Name != "" {
		profile.Name = req.Name
		middleware.DebugLog("🔍 DEBUG: Updated name: %s\n", req.Name)
	}
	if req.Email != "" {
		profile.Email = req.Email
		middleware.DebugLog("🔍 DEBUG: Updated email: %s\n", req.Email)
	}
	if req.Location != "" {
		profile.Location = req.Location
		middleware.DebugLog("🔍 DEBUG: Updated location: %s\n", req.Location)
	}
	if req.PhoneNumber != "" {
		profile.PhoneNumber = req.PhoneNumber
		middleware.DebugLog("🔍 DEBUG: Updated phone number: %s\n", req.PhoneNumber)
	}
	if req.Education != "" {
		profile.Education = req.Education
		middleware.DebugLog("🔍 DEBUG: Updated education: %s\n", req.Education)
	}
	if req.Portfolio != "" {
		profile.Portfolio = req.Portfolio
		middleware.DebugLog("🔍 DEBUG: Updated portfolio: %s\n", req.Portfolio)
	}
	if req.Linkedin != "" {
		profile.Linkedin = req.Linkedin
		middleware.DebugLog("🔍 DEBUG: Updated linkedin: %s\n", req.Linkedin)
	}
	if req.Github != "" {
		profile.Github = req.Github
		middleware.DebugLog("🔍 DEBUG: Updated github: %s\n", req.Github)
	}
	if req.Experience != nil {
		profile.Experience = *req.Experience
		middleware.DebugLog("🔍 DEBUG: Updated experience: %f\n", *req.Experience)
	}
	// Always update skills if they were provided (even if empty array)
	if req.Skills != nil {
		profile.Skills = req.Skills
		middleware.DebugLog("🔍 DEBUG: Updated skills: %v (length: %d)\n", req.Skills, len(req.Skills))
	} else {
		middleware.DebugLog("🔍 DEBUG: No skills to update (req.Skills is nil)\n")
	}

	// Update file fields
	if req.ResumeKey != "" {
		profile.ResumeKey = req.ResumeKey
		middleware.DebugLog("🔍 DEBUG: Updated resume - Key: %s\n", req.ResumeKey)
	}
	if req.ProfilePhotoKey != "" { // Changed from ProfilePhotoKey != nil to ProfilePhotoKey != ""
		profile.ProfilePhotoKey = req.ProfilePhotoKey
		middleware.DebugLog("🔍 DEBUG: Updated profile photo - Key: %s\n", req.ProfilePhotoKey)
	}

	// Update profile
	middleware.DebugLog("🔍 DEBUG: Profile ID: %s\n", profile.ID)
	if profile.ID == "" {
		middleware.DebugLog("🔍 DEBUG: Creating new profile\n")
		err = h.service.CreateProfile(c.Request.Context(), profile)
		if err != nil {
			middleware.DebugLog("❌ DEBUG: CreateProfile error: %v\n", err)
		} else {
			middleware.DebugLog("✅ DEBUG: Profile created successfully\n")
		}
	} else {
		middleware.DebugLog("🔍 DEBUG: Updating existing profile\n")
		err = h.service.UpdateProfile(c.Request.Context(), profile)
		if err != nil {
			middleware.DebugLog("❌ DEBUG: UpdateProfile error: %v\n", err)
		} else {
			middleware.DebugLog("✅ DEBUG: Profile updated successfully\n")
		}
	}

	if err != nil {
		middleware.DebugLog("❌ DEBUG: Service error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update profile",
			"error":   err.Error(),
		})
		return
	}

	middleware.DebugLog("✅ DEBUG: Profile operation completed successfully - ID: %s, Name: %s, CertificatesCount: %d\n", profile.ID, profile.Name, len(profile.Certificates))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
		"data":    profile,
	})
	middleware.DebugLog("✅ DEBUG: Response sent successfully\n")
}

// @Summary Add Certificate to My Profile
// @Description Add a certificate to the current user's student profile
// @Tags Student Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateCertificateRequest true "Certificate data"
// @Success 200 {object} map[string]interface{} "Certificate added successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Router /api/students/me/certificates [post]
// POST /students/me/certificates
func (h *StudentProfileHandler) AddMyCertificate(c *gin.Context) {
	middleware.DebugLog("DEBUG: AddMyCertificate called\n")
	username := c.GetString("username")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_student_profile", "update", userID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	var req UpdateCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	// Get or create student profile
	profile, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		// Create new profile
		profile = &StudentProfile{
			UserID: userID,
			Name:   c.GetString("name"),
			Email:  c.GetString("email"),
		}
		err = h.service.CreateProfile(c.Request.Context(), profile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to create profile",
				"error":   err.Error(),
			})
			return
		}
	}

	// Create certificate record
	certificate := &Certificate{
		StudentProfileID: profile.ID,
		Name:             req.Name,
		FileKey:          req.FileKey,
		IssueDate:        req.IssueDate,
	}

	// Save certificate to database
	err = h.service.AddCertificate(certificate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save certificate record",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Certificate added successfully",
		"data":    certificate,
	})
}

// @Summary Upload My Resume
// @Description Upload a resume file to the current user's student profile
// @Tags Student Profile
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param resume formData file true "Resume file (PDF, DOC, DOCX, max 10MB)"
// @Success 200 {object} map[string]interface{} "Resume uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid file type or size"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Router /api/students/me/resume [post]
// POST /students/me/resume
func (h *StudentProfileHandler) UploadMyResume(c *gin.Context) {
	username := c.GetString("username")
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
	fileHeader, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Resume file is required",
			"error":   err.Error(),
		})
		return
	}

	// Validate file type
	if !IsValidResumeFile(fileHeader.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid file type. Allowed: PDF, DOC, DOCX",
		})
		return
	}

	// Validate file size (10MB max)
	if fileHeader.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "File size exceeds maximum allowed size (10MB)",
		})
		return
	}

	// Upload to S3
	key, err := h.storage.SaveResume(fileHeader, "resumes")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to upload resume",
		})
		return
	}

	// Get or create student profile
	profile, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		// Create new profile
		profile = &StudentProfile{
			UserID: userID,
			Name:   c.GetString("name"),
			Email:  c.GetString("email"),
		}
	}

	// Update resume key
	profile.ResumeKey = key

	// Save profile
	if profile.ID == "" {
		err = h.service.CreateProfile(c.Request.Context(), profile)
	} else {
		// Update the profile with resume data
		err = h.service.UpdateProfile(c.Request.Context(), profile)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update profile with resume",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Resume uploaded successfully",
		"data": gin.H{
			"file_url": fmt.Sprintf("/api/files/serve/resume/%s", userID),
		},
	})
}

// POST /students/me/certificate
// @Summary Upload Certificate
// @Description Upload a certificate for the current student profile
// @Tags Student Profile
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Certificate file (PDF, DOC, DOCX, JPG, JPEG, PNG)"
// @Param name formData string true "Certificate name"
// @Param issue_date formData string true "Issue date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{} "Certificate uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Failed to save certificate record"
// @Router /api/students/me/certificate [post]
// @x-swagger-ui true
func (h *StudentProfileHandler) UploadMyCertificate(c *gin.Context) {
	middleware.DebugLog("DEBUG: UploadMyCertificate called\n")
	username := c.GetString("username")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_student_profile", "update", userID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to parse form data",
			"error":   err.Error(),
		})
		return
	}

	// Get certificate file
	certificateFile, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Certificate file is required",
		})
		return
	}

	// Get certificate details from form
	certificateName := c.PostForm("name")
	issueDate := c.PostForm("issue_date")

	if certificateName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Certificate name is required",
		})
		return
	}

	if issueDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Issue date is required",
		})
		return
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(certificateFile.Filename))
	allowedTypes := []string{".pdf", ".doc", ".docx", ".jpg", ".jpeg", ".png"}
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
			"message": "Invalid file type. Allowed: PDF, DOC, DOCX, JPG, JPEG, PNG",
		})
		return
	}

	// Validate file size (10MB max)
	if certificateFile.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "File size exceeds maximum allowed size (10MB)",
		})
		return
	}

	// Upload certificate to S3
	key, err := h.storage.SaveDocument(certificateFile, "certificates")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to upload certificate",
		})
		return
	}

	// Get or create student profile
	profile, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		// Create new profile
		profile = &StudentProfile{
			UserID: userID,
			Name:   c.GetString("name"),
			Email:  c.GetString("email"),
		}
		err = h.service.CreateProfile(c.Request.Context(), profile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to create profile",
				"error":   err.Error(),
			})
			return
		}
	}

	// Get file metadata
	fileName := certificateFile.Filename
	fileType := certificateFile.Header.Get("Content-Type")
	if fileType == "" {
		fileType = getMimeTypeFromExtension(fileName)
	}
	fileSize := certificateFile.Size

	// Create certificate record
	certificate := &Certificate{
		StudentProfileID: profile.ID,
		Name:             certificateName,
		FileKey:          key, // Store S3 key
		FileName:         fileName,
		FileType:         fileType,
		FileSize:         fileSize,
		IssueDate:        issueDate,
	}

	// Save certificate to database
	err = h.service.AddCertificate(certificate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save certificate record",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Certificate uploaded successfully",
		"data": gin.H{
			"certificate": certificate,
			"file_url":    fmt.Sprintf("/api/files/serve/certificate/%s", certificate.ID),
		},
	})
}

// PUT /students/me/resume - Update resume field only
// @Summary Update Resume
// @Description Update the resume file for the current student profile
// @Tags Student Profile
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param resume formData file true "Resume file (PDF, DOC, DOCX)"
// @Success 200 {object} map[string]interface{} "Resume updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Failed to update profile"
// @Router /api/students/me/resume [put]
// @x-swagger-ui true
func (h *StudentProfileHandler) UpdateMyResume(c *gin.Context) {
	username := c.GetString("username")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_student_profile", "update", userID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	// Get the file from the request
	fileHeader, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Resume file is required",
			"error":   err.Error(),
		})
		return
	}

	// Validate file type
	if !IsValidResumeFile(fileHeader.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid file type. Allowed: PDF, DOC, DOCX",
		})
		return
	}

	// Validate file size (10MB max)
	if fileHeader.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "File size exceeds maximum allowed size (10MB)",
		})
		return
	}

	// Upload to S3
	key, err := h.storage.SaveResume(fileHeader, "resumes")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to upload resume",
		})
		return
	}

	// Get existing profile or create new one
	profile, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		// Create new profile
		profile = &StudentProfile{
			UserID: userID,
			Name:   c.GetString("name"),
			Email:  c.GetString("email"),
		}
	}

	// Update resume key
	profile.ResumeKey = key

	// Save profile
	if profile.ID == "" {
		err = h.service.CreateProfile(c.Request.Context(), profile)
	} else {
		err = h.service.UpdateProfile(c.Request.Context(), profile)
	}

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
		"message": "Resume updated successfully",
		"data":    profile,
	})
}

// POST /students/me/certificates
// @Summary Upload Certificate
// @Description Upload a certificate for the current student profile
// @Tags Student Profile
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Certificate file (PDF, DOC, DOCX, JPG, JPEG, PNG)"
// @Param name formData string true "Certificate name"
// @Param issue_date formData string true "Issue date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{} "Certificate uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Failed to save certificate record"
// @Router /api/students/me/certificates [post]
// @x-swagger-ui true
func (h *StudentProfileHandler) UploadCertificate(c *gin.Context) {
	username := c.GetString("username")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_student_profile", "update", userID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	middleware.DebugLog("DEBUG: UploadCertificate - UserID: %s\n", userID)

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		middleware.DebugLog("DEBUG: ParseMultipartForm error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to parse form data",
			"error":   err.Error(),
		})
		return
	}

	// Get certificate file
	certificateFile, err := c.FormFile("file")
	if err != nil {
		middleware.DebugLog("DEBUG: Certificate file error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Certificate file is required",
		})
		return
	}

	// Get certificate details from form
	certificateName := c.PostForm("name")
	issueDate := c.PostForm("issue_date")

	if certificateName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Certificate name is required",
		})
		return
	}

	if issueDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Issue date is required",
		})
		return
	}

	middleware.DebugLog("DEBUG: Certificate details - Name: %s, IssueDate: %s, File: %s, Size: %d\n",
		certificateName, issueDate, certificateFile.Filename, certificateFile.Size)

	// Validate file type
	ext := strings.ToLower(filepath.Ext(certificateFile.Filename))
	allowedTypes := []string{".pdf", ".doc", ".docx", ".jpg", ".jpeg", ".png"}
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
			"message": "Invalid file type. Allowed: PDF, DOC, DOCX, JPG, JPEG, PNG",
		})
		return
	}

	// Validate file size (10MB max)
	if certificateFile.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "File size exceeds maximum allowed size (10MB)",
		})
		return
	}

	// Upload certificate to S3
	key, err := h.storage.SaveDocument(certificateFile, "certificates")
	if err != nil {
		middleware.DebugLog("DEBUG: Failed to upload certificate to S3: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to upload certificate",
		})
		return
	}

	// Get file metadata
	fileName := certificateFile.Filename
	fileType := certificateFile.Header.Get("Content-Type")
	if fileType == "" {
		fileType = getMimeTypeFromExtension(fileName)
	}
	fileSize := certificateFile.Size

	middleware.DebugLog("DEBUG: File uploaded to S3 successfully - Key: %s, Size: %d bytes, Type: %s\n", key, fileSize, fileType)

	// Get or create student profile
	profile, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		middleware.DebugLog("DEBUG: Profile not found, creating new one\n")
		// Create new profile
		profile = &StudentProfile{
			UserID: userID,
			Name:   c.GetString("name"),
			Email:  c.GetString("email"),
		}
		err = h.service.CreateProfile(c.Request.Context(), profile)
		if err != nil {
			middleware.DebugLog("DEBUG: Failed to create profile: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to create profile",
				"error":   err.Error(),
			})
			return
		}
	} else {
		middleware.DebugLog("DEBUG: Existing profile found\n")
	}

	// Create certificate record
	certificate := &Certificate{
		StudentProfileID: profile.ID,
		Name:             certificateName,
		FileKey:          key, // Store S3 key
		FileName:         fileName,
		FileType:         fileType,
		FileSize:         fileSize,
		IssueDate:        issueDate,
	}

	// Save certificate to database
	err = h.service.AddCertificate(certificate)
	if err != nil {
		middleware.DebugLog("DEBUG: Failed to create certificate record: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save certificate record",
			"error":   err.Error(),
		})
		return
	}

	middleware.DebugLog("DEBUG: Certificate saved successfully - ID: %s\n", certificate.ID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Certificate uploaded successfully",
		"data": gin.H{
			"certificate": certificate,
			"file_url":    fmt.Sprintf("/api/files/serve/certificate/%s", certificate.ID),
		},
	})
}

// POST /students/me/certificates/add
// @Summary Add Certificate to Profile
// @Description Add a certificate to the current student profile
// @Tags Student Profile
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Certificate file (PDF, DOC, DOCX, JPG, JPEG, PNG)"
// @Param name formData string true "Certificate name"
// @Param issue_date formData string true "Issue date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{} "Certificate added successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Permission denied"
// @Failure 500 {object} map[string]interface{} "Failed to save certificate record"
// @Router /api/students/me/certificates/add [post]
// @x-swagger-ui true
func (h *StudentProfileHandler) AddCertificateToProfile(c *gin.Context) {
	middleware.DebugLog("DEBUG: AddCertificateToProfile called\n")
	username := c.GetString("username")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckLocalPermission(username, "db_asa_student_profile", "update", userID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	middleware.DebugLog("DEBUG: AddCertificateToProfile - UserID: %s\n", userID)

	// Parse multipart form data for file upload
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		middleware.DebugLog("DEBUG: ParseMultipartForm error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to parse form data",
			"error":   err.Error(),
		})
		return
	}

	// Get certificate file
	certificateFile, err := c.FormFile("file")
	if err != nil {
		middleware.DebugLog("DEBUG: Certificate file error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Certificate file is required",
		})
		return
	}

	// Get certificate details from form
	certificateName := c.PostForm("name")
	issueDate := c.PostForm("issue_date")

	if certificateName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Certificate name is required",
		})
		return
	}

	if issueDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Issue date is required",
		})
		return
	}

	middleware.DebugLog("DEBUG: Certificate request - Name: %s, File: %s, IssueDate: %s\n", certificateName, certificateFile.Filename, issueDate)

	// Validate file type
	ext := strings.ToLower(filepath.Ext(certificateFile.Filename))
	allowedTypes := []string{".pdf", ".doc", ".docx", ".jpg", ".jpeg", ".png"}
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
			"message": "Invalid file type. Allowed: PDF, DOC, DOCX, JPG, JPEG, PNG",
		})
		return
	}

	// Validate file size (10MB max)
	if certificateFile.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "File size exceeds maximum allowed size (10MB)",
		})
		return
	}

	// Upload certificate to S3
	key, err := h.storage.SaveDocument(certificateFile, "certificates")
	if err != nil {
		middleware.DebugLog("DEBUG: Failed to upload certificate to S3: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to upload certificate",
		})
		return
	}

	// Get file metadata
	fileName := certificateFile.Filename
	fileType := certificateFile.Header.Get("Content-Type")
	if fileType == "" {
		fileType = getMimeTypeFromExtension(fileName)
	}
	fileSize := certificateFile.Size

	middleware.DebugLog("DEBUG: File uploaded to S3 successfully - Key: %s\n", key)

	// Get or create student profile
	profile, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		middleware.DebugLog("DEBUG: Profile not found, creating new one\n")
		// Create new profile
		profile = &StudentProfile{
			UserID: userID,
			Name:   c.GetString("name"),
			Email:  c.GetString("email"),
		}
		err = h.service.CreateProfile(c.Request.Context(), profile)
		if err != nil {
			middleware.DebugLog("DEBUG: Failed to create profile: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to create profile",
				"error":   err.Error(),
			})
			return
		}
	} else {
		middleware.DebugLog("DEBUG: Existing profile found - ID: %s\n", profile.ID)
	}

	// Create certificate record
	certificate := &Certificate{
		StudentProfileID: profile.ID,
		Name:             certificateName,
		FileKey:          key, // Store S3 key
		FileName:         fileName,
		FileType:         fileType,
		FileSize:         fileSize,
		IssueDate:        issueDate,
	}

	// Save certificate to database
	err = h.service.AddCertificate(certificate)
	if err != nil {
		middleware.DebugLog("DEBUG: Failed to create certificate record: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save certificate record",
			"error":   err.Error(),
		})
		return
	}

	middleware.DebugLog("DEBUG: Certificate saved successfully - ID: %s\n", certificate.ID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Certificate added successfully",
		"data":    certificate,
	})
}

// Helper function to validate resume file type
func IsValidResumeFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	allowedTypes := []string{".pdf", ".doc", ".docx"}
	for _, allowedExt := range allowedTypes {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

// Helper function to validate document file type
func IsValidDocumentFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	allowedTypes := []string{".pdf", ".doc", ".docx", ".txt", ".rtf"}
	for _, allowedExt := range allowedTypes {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

// Helper function to get MIME type from file extension
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
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".txt":
		return "text/plain"
	case ".rtf":
		return "application/rtf"
	default:
		return "application/octet-stream"
	}
}

// Helper function to validate image file type
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
