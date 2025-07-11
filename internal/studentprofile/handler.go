// File: internal/studentprofile/handler.go

package studentprofile

import (
	"asa/pkg/authz"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
	if profile.Resume != "" {
		if !strings.HasPrefix(profile.Resume, "http") {
			trimmed := profile.Resume
			if strings.HasPrefix(trimmed, "/") {
				trimmed = trimmed[1:]
			}
			if strings.HasPrefix(trimmed, "uploads/") {
				trimmed = trimmed[len("uploads/"):]
			}
			profile.Resume = "http://localhost:3000/api/files/serve/" + trimmed
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Profile fetched", "data": profile})
}

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
	if profile.Resume != "" {
		if !strings.HasPrefix(profile.Resume, "http") {
			trimmed := profile.Resume
			if strings.HasPrefix(trimmed, "/") {
				trimmed = trimmed[1:]
			}
			if strings.HasPrefix(trimmed, "uploads/") {
				trimmed = trimmed[len("uploads/"):]
			}
			profile.Resume = "http://localhost:3000/api/files/serve/" + trimmed
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Profile fetched", "data": profile})
}

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

	// Check if this is a multipart form (file upload) or JSON
	contentType := c.GetHeader("Content-Type")
	fmt.Printf("DEBUG: UpdateMyProfile - Content-Type: %s\n", contentType)
	fmt.Printf("DEBUG: UpdateMyProfile - UserID: %s\n", userID)
	fmt.Printf("DEBUG: UpdateMyProfile - All headers: %+v\n", c.Request.Header)

	// Log request body for debugging
	body, err := c.GetRawData()
	if err == nil {
		fmt.Printf("DEBUG: Raw request body: %s\n", string(body))
		// Restore the body for further processing
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	var req UpdateStudentProfileRequest

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

			// Generate unique filename
			timestamp := time.Now().UnixNano()
			ext := filepath.Ext(resumeFile.Filename)
			baseName := strings.TrimSuffix(resumeFile.Filename, ext)
			safeBaseName := strings.ReplaceAll(baseName, " ", "_")
			filename := fmt.Sprintf("%d_%s%s", timestamp, safeBaseName, ext)

			// Create uploads/resumes directory if it doesn't exist
			uploadDir := "uploads/resumes"
			if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Failed to create upload directory",
				})
				return
			}

			// Save file
			dst := filepath.Join(uploadDir, filename)
			if err := c.SaveUploadedFile(resumeFile, dst); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Failed to save file",
					"error":   err.Error(),
				})
				return
			}

			// Set resume path in request
			resumePath := filepath.Join("resumes", filename)
			req.Resume = resumePath
			fmt.Printf("DEBUG: Resume path set: %s\n", resumePath)
		} else {
			fmt.Printf("DEBUG: No resume file found: %v\n", err)
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
		if profilePhoto := c.PostForm("profile_photo"); profilePhoto != "" {
			req.ProfilePhoto = profilePhoto
			fmt.Printf("DEBUG: Profile photo from form: %s\n", profilePhoto)
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
		if experienceStr := c.PostForm("experience"); experienceStr != "" {
			if experience, err := strconv.ParseFloat(experienceStr, 64); err == nil {
				req.Experience = &experience
				fmt.Printf("DEBUG: Experience from form: %f\n", experience)
			}
		}
		if skillsStr := c.PostForm("skills"); skillsStr != "" {
			// Parse skills as comma-separated string
			skills := strings.Split(skillsStr, ",")
			for i, skill := range skills {
				skills[i] = strings.TrimSpace(skill)
			}
			req.Skills = skills
			fmt.Printf("DEBUG: Skills from form: %v\n", skills)
		}
	} else {
		fmt.Printf("DEBUG: Handling JSON data\n")
		// Handle JSON data
		if err := c.ShouldBindJSON(&req); err != nil {
			fmt.Printf("DEBUG: JSON binding error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid request",
				"error":   err.Error(),
				"details": "Check that all required fields are present and properly formatted",
			})
			return
		}
		fmt.Printf("DEBUG: JSON request parsed successfully: %+v\n", req)
	}

	// Get existing profile or create new one
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		fmt.Printf("DEBUG: Profile not found, creating new one\n")
		// Create new profile
		profile = &StudentProfile{
			UserID: userID,
			Name:   c.GetString("name"),
			Email:  c.GetString("email"),
		}
	} else {
		fmt.Printf("DEBUG: Existing profile found\n")
	}

	// Update profile fields
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
	if req.ProfilePhoto != "" {
		profile.ProfilePhoto = req.ProfilePhoto
	}
	if req.Resume != "" {
		profile.Resume = req.Resume
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

	fmt.Printf("DEBUG: About to save profile: %+v\n", profile)

	// Save profile
	if profile.ID == "" {
		err = h.service.CreateProfile(profile)
		fmt.Printf("DEBUG: Creating new profile\n")
	} else {
		err = h.service.UpdateProfile(profile)
		fmt.Printf("DEBUG: Updating existing profile\n")
	}

	if err != nil {
		fmt.Printf("DEBUG: Profile save error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update profile",
			"error":   err.Error(),
		})
		return
	}

	fmt.Printf("DEBUG: Profile saved successfully\n")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
		"data":    profile,
	})
}

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

	// Generate unique filename
	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(fileHeader.Filename)
	baseName := strings.TrimSuffix(fileHeader.Filename, ext)
	safeBaseName := strings.ReplaceAll(baseName, " ", "_")
	filename := fmt.Sprintf("%d_%s%s", timestamp, safeBaseName, ext)

	// Create uploads/resumes directory if it doesn't exist
	uploadDir := "uploads/resumes"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create upload directory",
		})
		return
	}

	// Save file
	dst := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(fileHeader, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save file",
			"error":   err.Error(),
		})
		return
	}

	// Update student profile with resume path
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		// If profile doesn't exist, create it
		profile = &StudentProfile{
			UserID: userID,
			Name:   c.GetString("name"),
			Email:  c.GetString("email"),
		}
	}

	// Set resume path (relative to uploads directory)
	profile.Resume = filepath.Join("resumes", filename)

	// Update or create profile
	if profile.ID == "" {
		err = h.service.CreateProfile(profile)
	} else {
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
			"resume_path": profile.Resume,
			"file_name":   filename,
			"file_size":   fileHeader.Size,
		},
	})
}

// POST /students/me/certificate
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

	// Generate unique filename
	timestamp := time.Now().UnixNano()
	ext = filepath.Ext(certificateFile.Filename)
	baseName := strings.TrimSuffix(certificateFile.Filename, ext)
	safeBaseName := strings.ReplaceAll(baseName, " ", "_")
	filename := fmt.Sprintf("%d_%s%s", timestamp, safeBaseName, ext)

	// Create uploads/certificates directory if it doesn't exist
	uploadDir := "uploads/certificates"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create upload directory",
		})
		return
	}

	// Save file
	dst := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(certificateFile, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save file",
			"error":   err.Error(),
		})
		return
	}

	// Set certificate file path
	certificatePath := filepath.Join("certificates", filename)

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
		Name:             certificateName,
		File:             certificatePath,
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
			"file_path":   certificatePath,
		},
	})
}

// PUT /students/me/resume - Update resume field only
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

	var req struct {
		ResumePath string `json:"resume_path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
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

	// Update resume path
	profile.Resume = req.ResumePath

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

	// Generate unique filename
	timestamp := time.Now().UnixNano()
	ext = filepath.Ext(certificateFile.Filename)
	baseName := strings.TrimSuffix(certificateFile.Filename, ext)
	safeBaseName := strings.ReplaceAll(baseName, " ", "_")
	filename := fmt.Sprintf("%d_%s%s", timestamp, safeBaseName, ext)

	// Create uploads/certificates directory if it doesn't exist
	uploadDir := "uploads/certificates"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create upload directory",
		})
		return
	}

	// Save file
	dst := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(certificateFile, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save file",
			"error":   err.Error(),
		})
		return
	}

	// Set certificate file path
	certificatePath := filepath.Join("certificates", filename)
	fmt.Printf("DEBUG: Certificate path set: %s\n", certificatePath)

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
		File:             certificatePath,
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
			"file_path":   certificatePath,
		},
	})
}

// POST /students/me/certificates/add
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

	// Parse JSON request
	var req struct {
		Name      string `json:"name" binding:"required"`
		File      string `json:"file" binding:"required"`
		IssueDate string `json:"issue_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("DEBUG: JSON binding error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	fmt.Printf("DEBUG: Certificate request - Name: %s, File: %s, IssueDate: %s\n", req.Name, req.File, req.IssueDate)

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
		Name:             req.Name,
		File:             req.File,
		IssueDate:        req.IssueDate,
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
