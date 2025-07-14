// File: internal/studentprofile/handler.go

package studentprofile

import (
	"asa/pkg/authz"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type StudentProfileHandler struct {
	service StudentProfileService
}

func NewStudentProfileHandler(s StudentProfileService) *StudentProfileHandler {
	return &StudentProfileHandler{s}
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
	allowed, err := authz.CheckAAAPermission(username, "db_asa_student_profile", "read", profileID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	profile, err := h.service.GetProfile(profileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Profile not found"})
		return
	}

	// Rewrite resume field to always use /api/files/serve/
	if profile.Resume != nil {
		profile.Resume = profile.Resume
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
	allowed, err := authz.CheckAAAPermission(username, "db_asa_student_profile", "update", profileID, jwtToken)
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
	err = h.service.UpdateProfile(&req)
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
	allowed, err := authz.CheckAAAPermission(username, "db_asa_student_profile", "read", userID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Profile not found"})
		return
	}

	// Rewrite resume field to always use /api/files/serve/
	if profile.Resume != nil {
		profile.Resume = profile.Resume
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
	allowed, err := authz.CheckAAAPermission(username, "db_asa_certificates", "create", studentID, jwtToken)
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
	studentProfile, err := h.service.GetProfile(studentID)
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
	allowed, err := authz.CheckAAAPermission(username, "db_asa_certificates", "delete", certificateID, jwtToken)
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
	username := c.GetString("username")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_student_profile", "update", userID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	var req UpdateStudentProfileRequest
	contentType := c.GetHeader("Content-Type")

	fmt.Printf("DEBUG: Content-Type: %s\n", contentType)
	fmt.Printf("DEBUG: User ID: %s\n", userID)

	if strings.Contains(contentType, "multipart/form-data") {
		fmt.Printf("DEBUG: Handling multipart form data\n")
		// Handle multipart form data (file upload)
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			fmt.Printf("DEBUG: ParseMultipartForm error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Failed to parse form data",
				"error":   err.Error(),
			})
			return
		}

		// Handle resume file upload
		if resumeFile, err := c.FormFile("resume"); err == nil {
			fmt.Printf("DEBUG: Resume file found - Name: %s, Size: %d\n", resumeFile.Filename, resumeFile.Size)
			// Validate file type
			if !IsValidResumeFile(resumeFile.Filename) {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Invalid file type. Allowed: PDF, DOC, DOCX",
				})
				return
			}

			// Validate file size (10MB max)
			if resumeFile.Size > 10*1024*1024 {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "File size exceeds maximum allowed size (10MB)",
				})
				return
			}

			// Read file into bytes
			file, err := resumeFile.Open()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Failed to read file",
				})
				return
			}
			defer file.Close()

			fileBytes, err := io.ReadAll(file)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Failed to read file",
				})
				return
			}

			// Set resume data in request
			req.Resume = fileBytes
			req.ResumeName = resumeFile.Filename
			req.ResumeType = resumeFile.Header.Get("Content-Type")
			if req.ResumeType == "" {
				req.ResumeType = getMimeTypeFromExtension(resumeFile.Filename)
			}
			req.ResumeSize = resumeFile.Size
			fmt.Printf("DEBUG: Resume data set - Size: %d bytes, Type: %s\n", len(fileBytes), req.ResumeType)
		} else {
			fmt.Printf("DEBUG: No resume file found: %v\n", err)
		}

		// Handle profile photo upload
		if profilePhotoFile, err := c.FormFile("profile_photo"); err == nil {
			fmt.Printf("DEBUG: Profile photo found - Name: %s, Size: %d\n", profilePhotoFile.Filename, profilePhotoFile.Size)
			// Validate file type
			if !IsValidImageFile(profilePhotoFile.Filename) {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Invalid image type. Allowed: JPG, PNG, GIF, WebP",
				})
				return
			}

			// Validate file size (5MB max for images)
			if profilePhotoFile.Size > 5*1024*1024 {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Image size exceeds maximum allowed size (5MB)",
				})
				return
			}

			// Read file into bytes
			file, err := profilePhotoFile.Open()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Failed to read image file",
				})
				return
			}
			defer file.Close()

			fileBytes, err := io.ReadAll(file)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Failed to read image file",
				})
				return
			}

			// Set profile photo data in request
			req.ProfilePhoto = fileBytes
			req.ProfilePhotoName = profilePhotoFile.Filename
			req.ProfilePhotoType = profilePhotoFile.Header.Get("Content-Type")
			if req.ProfilePhotoType == "" {
				req.ProfilePhotoType = getMimeTypeFromExtension(profilePhotoFile.Filename)
			}
			req.ProfilePhotoSize = profilePhotoFile.Size
			fmt.Printf("DEBUG: Profile photo data set - Size: %d bytes, Type: %s\n", len(fileBytes), req.ProfilePhotoType)
		}

		// Handle other form fields
		if name := c.PostForm("name"); name != "" {
			req.Name = name
			fmt.Printf("DEBUG: Name from form: %s\n", name)
		}
		if email := c.PostForm("email"); email != "" {
			req.Email = email
			fmt.Printf("DEBUG: Email from form: %s\n", email)
		}
		if location := c.PostForm("location"); location != "" {
			req.Location = location
			fmt.Printf("DEBUG: Location from form: %s\n", location)
		}
		if phoneNumber := c.PostForm("phone_number"); phoneNumber != "" {
			req.PhoneNumber = phoneNumber
			fmt.Printf("DEBUG: Phone number from form: %s\n", phoneNumber)
		}
		if education := c.PostForm("education"); education != "" {
			req.Education = education
			fmt.Printf("DEBUG: Education from form: %s\n", education)
		}
		if portfolio := c.PostForm("portfolio"); portfolio != "" {
			req.Portfolio = portfolio
			fmt.Printf("DEBUG: Portfolio from form: %s\n", portfolio)
		}
		if linkedin := c.PostForm("linkedin"); linkedin != "" {
			req.Linkedin = linkedin
			fmt.Printf("DEBUG: LinkedIn from form: %s\n", linkedin)
		}
		if github := c.PostForm("github"); github != "" {
			req.Github = github
			fmt.Printf("DEBUG: GitHub from form: %s\n", github)
		}
		if skills := c.PostForm("skills"); skills != "" {
			// Handle skills as comma-separated string and convert to array
			skillsArray := strings.Split(skills, ",")
			for i, skill := range skillsArray {
				skillsArray[i] = strings.TrimSpace(skill)
			}
			req.Skills = Skills(skillsArray)
			fmt.Printf("DEBUG: Skills from form: %s -> %v\n", skills, req.Skills)
		}
	} else {
		// Handle JSON request
		if err := c.ShouldBindJSON(&req); err != nil {
			fmt.Printf("DEBUG: JSON binding error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid request format",
				"error":   err.Error(),
			})
			return
		}
	}

	// Set user ID
	req.UserID = userID

	fmt.Printf("DEBUG: Final request: %+v\n", req)

	// Get existing profile or create new one
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		// Create new profile
		profile = &StudentProfile{
			UserID: userID,
			Name:   req.Name,
			Email:  req.Email,
		}
	}

	// Update profile fields from request
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
		profile.PhoneNumber = req.PhoneNumber
	}
	if req.Education != "" {
		profile.Education = req.Education
	}
	if req.Portfolio != "" {
		profile.Portfolio = req.Portfolio
	}
	if req.Linkedin != "" {
		profile.Linkedin = req.Linkedin
	}
	if req.Github != "" {
		profile.Github = req.Github
	}
	if req.Experience != nil {
		profile.Experience = *req.Experience
	}
	if req.Skills != nil {
		profile.Skills = req.Skills
	}

	// Update file fields
	if req.Resume != nil {
		profile.Resume = req.Resume
		profile.ResumeName = req.ResumeName
		profile.ResumeType = req.ResumeType
		profile.ResumeSize = req.ResumeSize
	}
	if req.ProfilePhoto != nil {
		profile.ProfilePhoto = req.ProfilePhoto
		profile.ProfilePhotoName = req.ProfilePhotoName
		profile.ProfilePhotoType = req.ProfilePhotoType
		profile.ProfilePhotoSize = req.ProfilePhotoSize
	}

	// Update profile
	if profile.ID == "" {
		err = h.service.CreateProfile(profile)
	} else {
		err = h.service.UpdateProfile(profile)
	}

	if err != nil {
		fmt.Printf("DEBUG: Service error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update profile",
			"error":   err.Error(),
		})
		return
	}

	fmt.Printf("DEBUG: Profile updated successfully: %+v\n", profile)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
		"data":    profile,
	})
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
	fmt.Printf("DEBUG: AddMyCertificate called\n")
	username := c.GetString("username")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_student_profile", "update", userID, jwtToken)
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
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		// Create new profile
		profile = &StudentProfile{
			UserID: userID,
			Name:   c.GetString("name"),
			Email:  c.GetString("email"),
		}
		err = h.service.CreateProfile(profile)
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
		File:             req.File,
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
	allowed, err := authz.CheckAAAPermission(username, "db_asa_files", "create", "", jwtToken)
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

	// Read file into bytes
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to read file",
		})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to read file",
		})
		return
	}

	// Get or create student profile
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		// Create new profile
		profile = &StudentProfile{
			UserID: userID,
			Name:   c.GetString("name"),
			Email:  c.GetString("email"),
		}
	}

	// Update resume data
	profile.Resume = fileBytes
	profile.ResumeName = fileHeader.Filename
	profile.ResumeType = fileHeader.Header.Get("Content-Type")
	if profile.ResumeType == "" {
		profile.ResumeType = getMimeTypeFromExtension(fileHeader.Filename)
	}
	profile.ResumeSize = fileHeader.Size

	// Save profile
	if profile.ID == "" {
		err = h.service.CreateProfile(profile)
	} else {
		// Update the profile with resume data
		profile.Resume = fileBytes
		profile.ResumeName = fileHeader.Filename
		profile.ResumeType = profile.ResumeType
		profile.ResumeSize = fileHeader.Size
		err = h.service.UpdateProfile(profile)
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
			"resume_name": profile.ResumeName,
			"resume_type": profile.ResumeType,
			"resume_size": profile.ResumeSize,
			"file_url":    fmt.Sprintf("/api/files/serve/resume/%s", userID),
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
	fmt.Printf("DEBUG: UploadMyCertificate called\n")
	username := c.GetString("username")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_student_profile", "update", userID, jwtToken)
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

	// File system storage is no longer needed - we store binary data directly

	// Get or create student profile
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		// Create new profile
		profile = &StudentProfile{
			UserID: userID,
			Name:   c.GetString("name"),
			Email:  c.GetString("email"),
		}
		err = h.service.CreateProfile(profile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to create profile",
				"error":   err.Error(),
			})
			return
		}
	}

	// Read file into bytes
	file, err := certificateFile.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to read certificate file",
		})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to read certificate file",
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

	// Create certificate record
	certificate := &Certificate{
		StudentProfileID: profile.ID,
		Name:             certificateName,
		File:             fileBytes,
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
	allowed, err := authz.CheckAAAPermission(username, "db_asa_student_profile", "update", userID, jwtToken)
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

	// Read file into bytes
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to read file",
		})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to read file",
		})
		return
	}

	// Get existing profile or create new one
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		// Create new profile
		profile = &StudentProfile{
			UserID: userID,
			Name:   c.GetString("name"),
			Email:  c.GetString("email"),
		}
	}

	// Update resume data
	profile.Resume = fileBytes
	profile.ResumeName = fileHeader.Filename
	profile.ResumeType = fileHeader.Header.Get("Content-Type")
	if profile.ResumeType == "" {
		profile.ResumeType = getMimeTypeFromExtension(fileHeader.Filename)
	}
	profile.ResumeSize = fileHeader.Size

	// Save profile
	if profile.ID == "" {
		err = h.service.CreateProfile(profile)
	} else {
		err = h.service.UpdateProfile(profile)
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
	allowed, err := authz.CheckAAAPermission(username, "db_asa_student_profile", "update", userID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	fmt.Printf("DEBUG: UploadCertificate - UserID: %s\n", userID)

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		fmt.Printf("DEBUG: ParseMultipartForm error: %v\n", err)
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
		fmt.Printf("DEBUG: Certificate file error: %v\n", err)
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

	fmt.Printf("DEBUG: Certificate details - Name: %s, IssueDate: %s, File: %s, Size: %d\n",
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

	// Read file into bytes
	file, err := certificateFile.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to read certificate file",
		})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to read certificate file",
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

	fmt.Printf("DEBUG: File read successfully - Size: %d bytes, Type: %s\n", len(fileBytes), fileType)

	// Get or create student profile
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		fmt.Printf("DEBUG: Profile not found, creating new one\n")
		// Create new profile
		profile = &StudentProfile{
			UserID: userID,
			Name:   c.GetString("name"),
			Email:  c.GetString("email"),
		}
		err = h.service.CreateProfile(profile)
		if err != nil {
			fmt.Printf("DEBUG: Failed to create profile: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to create profile",
				"error":   err.Error(),
			})
			return
		}
	} else {
		fmt.Printf("DEBUG: Existing profile found\n")
	}

	// Create certificate record
	certificate := &Certificate{
		StudentProfileID: profile.ID,
		Name:             certificateName,
		File:             fileBytes,
		FileName:         fileName,
		FileType:         fileType,
		FileSize:         fileSize,
		IssueDate:        issueDate,
	}

	// Save certificate to database
	err = h.service.AddCertificate(certificate)
	if err != nil {
		fmt.Printf("DEBUG: Failed to create certificate record: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save certificate record",
			"error":   err.Error(),
		})
		return
	}

	fmt.Printf("DEBUG: Certificate saved successfully - ID: %s\n", certificate.ID)

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
	fmt.Printf("DEBUG: AddCertificateToProfile called\n")
	username := c.GetString("username")
	userID := c.GetString("user_id")
	jwtToken := getJWT(c)
	allowed, err := authz.CheckAAAPermission(username, "db_asa_student_profile", "update", userID, jwtToken)
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	fmt.Printf("DEBUG: AddCertificateToProfile - UserID: %s\n", userID)

	// Parse multipart form data for file upload
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		fmt.Printf("DEBUG: ParseMultipartForm error: %v\n", err)
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
		fmt.Printf("DEBUG: Certificate file error: %v\n", err)
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

	fmt.Printf("DEBUG: Certificate request - Name: %s, File: %s, IssueDate: %s\n", certificateName, certificateFile.Filename, issueDate)

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

	// Read file into bytes
	file, err := certificateFile.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to read certificate file",
		})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to read certificate file",
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

	// Get or create student profile
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		fmt.Printf("DEBUG: Profile not found, creating new one\n")
		// Create new profile
		profile = &StudentProfile{
			UserID: userID,
			Name:   c.GetString("name"),
			Email:  c.GetString("email"),
		}
		err = h.service.CreateProfile(profile)
		if err != nil {
			fmt.Printf("DEBUG: Failed to create profile: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to create profile",
				"error":   err.Error(),
			})
			return
		}
	} else {
		fmt.Printf("DEBUG: Existing profile found - ID: %s\n", profile.ID)
	}

	// Create certificate record
	certificate := &Certificate{
		StudentProfileID: profile.ID,
		Name:             certificateName,
		File:             fileBytes,
		FileName:         fileName,
		FileType:         fileType,
		FileSize:         fileSize,
		IssueDate:        issueDate,
	}

	// Save certificate to database
	err = h.service.AddCertificate(certificate)
	if err != nil {
		fmt.Printf("DEBUG: Failed to create certificate record: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save certificate record",
			"error":   err.Error(),
		})
		return
	}

	fmt.Printf("DEBUG: Certificate saved successfully - ID: %s\n", certificate.ID)

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
