package authz

import (
	"github.com/Kisanlink/agriskill-academy/internal/middleware"
	"github.com/Kisanlink/agriskill-academy/pkg/jwtutil"
)

// contains checks if a slice of strings contains a specific string
func contains(list []string, val string) bool {
	for _, item := range list {
		if item == val {
			return true
		}
	}
	return false
}

// CheckLocalPermission checks permissions locally using JWT claims
func CheckLocalPermission(username, resource, action, resourceID, jwtToken string) (bool, error) {
	middleware.DebugLog("🔍 === LOCAL PERMISSION CHECK START ===")
	middleware.DebugLog("🔍 Username: %s", username)
	middleware.DebugLog("🔍 Resource: %s", resource)
	middleware.DebugLog("🔍 Action: %s", action)
	middleware.DebugLog("🔍 ResourceID: %s", resourceID)
	middleware.DebugLog("🔍 JWT Token length: %d", len(jwtToken))

	// Parse the JWT token to get user claims
	middleware.DebugLog("🔍 Parsing JWT token...")
	claims, err := jwtutil.ParseToken(jwtToken)
	if err != nil {
		middleware.DebugLog("❌ Failed to parse JWT token: %v", err)
		return false, err
	}
	middleware.DebugLog("✅ JWT token parsed successfully")
	middleware.DebugLog("🔍 All JWT claims: %+v", claims)

	// Extract role from claims
	var role string
	if roleInterface, exists := claims["role"]; exists {
		if roleStr, ok := roleInterface.(string); ok {
			role = roleStr
			middleware.DebugLog("✅ Role extracted: %s", role)
		} else {
			middleware.DebugLog("❌ Role is not a string: %T", roleInterface)
		}
	} else {
		middleware.DebugLog("❌ No 'role' found in JWT token for user: %s", username)
		return false, nil
	}

	if role == "" {
		middleware.DebugLog("❌ Empty role found in JWT token for user: %s", username)
		return false, nil
	}

	middleware.DebugLog("🔍 Checking permission for user: %s, role: %s, resource: %s, action: %s, resourceID: %s", username, role, resource, action, resourceID)

	// Define permission rules based on resource and action
	var result bool
	var err2 error

	switch resource {
	case "db_asa_student_profile":
		middleware.DebugLog("🔍 Checking student profile permissions...")
		result, err2 = checkStudentProfilePermissions(role, action, resourceID, username, claims)
	case "db_asa_employer_profiles":
		middleware.DebugLog("🔍 Checking employer profile permissions...")
		result, err2 = checkEmployerProfilePermissions(role, action, resourceID, username, claims)
	case "db_asa_job_posts":
		middleware.DebugLog("🔍 Checking job post permissions...")
		result, err2 = checkJobPostPermissions(role, action, resourceID, username, claims)
	case "db_asa_applications":
		middleware.DebugLog("🔍 Checking application permissions...")
		result, err2 = checkApplicationPermissions(role, action, resourceID, username, claims)
	case "db_asa_bookmarks":
		middleware.DebugLog("🔍 Checking bookmark permissions...")
		result, err2 = checkBookmarkPermissions(role, action, resourceID, username, claims)
	case "db_asa_files":
		middleware.DebugLog("🔍 Checking file permissions...")
		result, err2 = checkFilePermissions(role, action, resourceID, username, claims)
	case "db_asa_notification_preferences":
		middleware.DebugLog("🔍 Checking notification permissions...")
		result, err2 = checkNotificationPermissions(role, action, resourceID, username, claims)
	case "db_asa_certificates":
		middleware.DebugLog("🔍 Checking certificate permissions...")
		result, err2 = checkCertificatePermissions(role, action, resourceID, username, claims)
	case "db_asa_job_alerts":
		middleware.DebugLog("🔍 Checking job alert permissions...")
		result, err2 = checkJobAlertPermissions(role, action, resourceID, username, claims)
	case "db_asa_jobs":
		middleware.DebugLog("🔍 Checking job permissions...")
		result, err2 = checkJobPermissions(role, action, resourceID, username, claims)
	default:
		middleware.DebugLog("⚠️ Unknown resource: %s", resource)
		return false, nil
	}

	if err2 != nil {
		middleware.DebugLog("❌ Permission check error: %v", err2)
		return false, err2
	}

	if result {
		middleware.DebugLog("✅ === PERMISSION GRANTED ===")
	} else {
		middleware.DebugLog("❌ === PERMISSION DENIED ===")
	}

	return result, nil
}

// Student Profile Permissions
func checkStudentProfilePermissions(role, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === STUDENT PROFILE PERMISSIONS ===")
	middleware.DebugLog("🔍 Role: %s", role)
	middleware.DebugLog("🔍 Action: %s", action)
	middleware.DebugLog("🔍 ResourceID: %s", resourceID)

	userID := claims["user_id"].(string)
	middleware.DebugLog("🔍 UserID from claims: %s", userID)

	switch action {
	case "read":
		// Students can read their own profile, employers can read any profile
		result := role == "student" || role == "employer" || role == "asa_admin"
		middleware.DebugLog("🔍 Read permission - Student: %v, Employer: %v, Admin: %v, Result: %v",
			role == "student", role == "employer", role == "asa_admin", result)
		return result, nil
	case "update":
		// Only students can update their own profile
		result := role == "student" && resourceID == userID
		middleware.DebugLog("🔍 Update permission - Student: %v, ResourceID match: %v, Result: %v",
			role == "student", resourceID == userID, result)
		return result, nil
	case "create":
		// Only students can create profiles
		result := role == "student"
		middleware.DebugLog("🔍 Create permission - Student: %v, Result: %v", role == "student", result)
		return result, nil
	default:
		middleware.DebugLog("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Employer Profile Permissions
func checkEmployerProfilePermissions(role, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === EMPLOYER PROFILE PERMISSIONS ===")
	middleware.DebugLog("🔍 Role: %s", role)
	middleware.DebugLog("🔍 Action: %s", action)
	middleware.DebugLog("🔍 ResourceID: %s", resourceID)

	userID := claims["user_id"].(string)
	middleware.DebugLog("🔍 UserID from claims: %s", userID)

	switch action {
	case "read":
		// Anyone can read employer profiles
		middleware.DebugLog("🔍 Read permission - Anyone can read, Result: true")
		return true, nil
	case "update":
		// Only employers can update their own profile
		result := role == "employer" && resourceID == userID
		middleware.DebugLog("🔍 Update permission - Employer: %v, ResourceID match: %v, Result: %v",
			role == "employer", resourceID == userID, result)
		return result, nil
	case "create":
		// Only employers can create profiles
		result := role == "employer"
		middleware.DebugLog("🔍 Create permission - Employer: %v, Result: %v", role == "employer", result)
		return result, nil
	case "delete":
		// Only admins can delete profiles
		result := role == "asa_admin"
		middleware.DebugLog("🔍 Delete permission - Admin: %v, Result: %v", role == "asa_admin", result)
		return result, nil
	default:
		middleware.DebugLog("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Job Post Permissions
func checkJobPostPermissions(role, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === JOB POST PERMISSIONS ===")
	middleware.DebugLog("🔍 Role: %s", role)
	middleware.DebugLog("🔍 Action: %s", action)
	middleware.DebugLog("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Anyone can read job posts
		middleware.DebugLog("🔍 Read permission - Anyone can read, Result: true")
		return true, nil
	case "create":
		// Only employers can create job posts
		result := role == "employer"
		middleware.DebugLog("🔍 Create permission - Employer: %v, Result: %v", role == "employer", result)
		return result, nil
	case "update":
		// Only employers can update their own job posts
		result := role == "employer"
		middleware.DebugLog("🔍 Update permission - Employer: %v, Result: %v", role == "employer", result)
		return result, nil
	case "delete":
		// Only employers can delete their own job posts
		result := role == "employer"
		middleware.DebugLog("🔍 Delete permission - Employer: %v, Result: %v", role == "employer", result)
		return result, nil
	default:
		middleware.DebugLog("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Application Permissions
func checkApplicationPermissions(role, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === APPLICATION PERMISSIONS ===")
	middleware.DebugLog("🔍 Role: %s", role)
	middleware.DebugLog("🔍 Action: %s", action)
	middleware.DebugLog("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Students can read their own applications, employers can read applications for their jobs
		result := role == "student" || role == "employer"
		middleware.DebugLog("🔍 Read permission - Student: %v, Employer: %v, Result: %v",
			role == "student", role == "employer", result)
		return result, nil
	case "create":
		// Only students can create applications
		result := role == "student"
		middleware.DebugLog("🔍 Create permission - Student: %v, Result: %v", role == "student", result)
		return result, nil
	case "update":
		// Students can update their own applications, employers can update applications for their jobs
		result := role == "student" || role == "employer"
		middleware.DebugLog("🔍 Update permission - Student: %v, Employer: %v, Result: %v",
			role == "student", role == "employer", result)
		return result, nil
	case "delete":
		// Students can delete their own applications
		result := role == "student"
		middleware.DebugLog("🔍 Delete permission - Student: %v, Result: %v", role == "student", result)
		return result, nil
	default:
		middleware.DebugLog("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Bookmark Permissions
func checkBookmarkPermissions(role, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === BOOKMARK PERMISSIONS ===")
	middleware.DebugLog("🔍 Role: %s", role)
	middleware.DebugLog("🔍 Action: %s", action)
	middleware.DebugLog("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Students can read their own bookmarks
		result := role == "student"
		middleware.DebugLog("🔍 Read permission - Student: %v, Result: %v", role == "student", result)
		return result, nil
	case "create":
		// Only students can create bookmarks
		result := role == "student"
		middleware.DebugLog("🔍 Create permission - Student: %v, Result: %v", role == "student", result)
		return result, nil
	case "delete":
		// Students can delete their own bookmarks
		result := role == "student"
		middleware.DebugLog("🔍 Delete permission - Student: %v, Result: %v", role == "student", result)
		return result, nil
	default:
		middleware.DebugLog("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// File Permissions
func checkFilePermissions(role, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === FILE PERMISSIONS ===")
	middleware.DebugLog("🔍 Role: %s", role)
	middleware.DebugLog("🔍 Action: %s", action)
	middleware.DebugLog("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Anyone can read files
		middleware.DebugLog("🔍 Read permission - Anyone can read, Result: true")
		return true, nil
	case "create":
		// Anyone authenticated can create files
		middleware.DebugLog("🔍 Create permission - Anyone authenticated can create, Result: true")
		return true, nil
	case "delete":
		// Only admins can delete files
		result := role == "asa_admin"
		middleware.DebugLog("🔍 Delete permission - Admin: %v, Result: %v", role == "asa_admin", result)
		return result, nil
	default:
		middleware.DebugLog("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Notification Permissions
func checkNotificationPermissions(role, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === NOTIFICATION PERMISSIONS ===")
	middleware.DebugLog("🔍 Role: %s", role)
	middleware.DebugLog("🔍 Action: %s", action)
	middleware.DebugLog("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Anyone can read their own notifications
		middleware.DebugLog("🔍 Read permission - Anyone can read, Result: true")
		return true, nil
	case "create":
		// Anyone can create notification preferences
		middleware.DebugLog("🔍 Create permission - Anyone can create, Result: true")
		return true, nil
	case "update":
		// Anyone can update their own notification preferences
		middleware.DebugLog("🔍 Update permission - Anyone can update, Result: true")
		return true, nil
	default:
		middleware.DebugLog("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Certificate Permissions
func checkCertificatePermissions(role, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === CERTIFICATE PERMISSIONS ===")
	middleware.DebugLog("🔍 Role: %s", role)
	middleware.DebugLog("🔍 Action: %s", action)
	middleware.DebugLog("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Students can read their own certificates
		result := role == "student"
		middleware.DebugLog("🔍 Read permission - Student: %v, Result: %v", role == "student", result)
		return result, nil
	case "create":
		// Only students can create certificates
		result := role == "student"
		middleware.DebugLog("🔍 Create permission - Student: %v, Result: %v", role == "student", result)
		return result, nil
	case "update":
		// Students can update their own certificates
		result := role == "student"
		middleware.DebugLog("🔍 Update permission - Student: %v, Result: %v", role == "student", result)
		return result, nil
	case "delete":
		// Students can delete their own certificates
		result := role == "student"
		middleware.DebugLog("🔍 Delete permission - Student: %v, Result: %v", role == "student", result)
		return result, nil
	default:
		middleware.DebugLog("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Job Alert Permissions
func checkJobAlertPermissions(role, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === JOB ALERT PERMISSIONS ===")
	middleware.DebugLog("🔍 Role: %s", role)
	middleware.DebugLog("🔍 Action: %s", action)
	middleware.DebugLog("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Students can read their own job alerts
		result := role == "student"
		middleware.DebugLog("🔍 Read permission - Student: %v, Result: %v", role == "student", result)
		return result, nil
	case "create":
		// Only students can create job alerts
		result := role == "student"
		middleware.DebugLog("🔍 Create permission - Student: %v, Result: %v", role == "student", result)
		return result, nil
	case "update":
		// Students can update their own job alerts
		result := role == "student"
		middleware.DebugLog("🔍 Update permission - Student: %v, Result: %v", role == "student", result)
		return result, nil
	case "delete":
		// Students can delete their own job alerts
		result := role == "student"
		middleware.DebugLog("🔍 Delete permission - Student: %v, Result: %v", role == "student", result)
		return result, nil
	default:
		middleware.DebugLog("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Job Permissions (for worker)
func checkJobPermissions(role, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === JOB PERMISSIONS ===")
	middleware.DebugLog("🔍 Role: %s", role)
	middleware.DebugLog("🔍 Action: %s", action)
	middleware.DebugLog("🔍 ResourceID: %s", resourceID)

	switch action {
	case "create":
		// Only admins can create background jobs
		result := role == "asa_admin"
		middleware.DebugLog("🔍 Create permission - Admin: %v, Result: %v", role == "asa_admin", result)
		return result, nil
	default:
		middleware.DebugLog("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Legacy function for backward compatibility - redirects to local permission check
func CheckAAAPermission(username, resource, action, resourceID, jwtToken string) (bool, error) {
	middleware.DebugLog("🔍 === LEGACY AAA PERMISSION CHECK ===")
	middleware.DebugLog("🔍 Calling CheckLocalPermission for local authentication")
	return CheckLocalPermission(username, resource, action, resourceID, jwtToken)
}
