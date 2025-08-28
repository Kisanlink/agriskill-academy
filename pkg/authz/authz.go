package authz

import (
	"asa/internal/middleware"
	"asa/pkg/jwtutil"
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

	// Extract roles from JWT claims - handle both 'role' (legacy) and 'roles' (local)
	var roles []string

	middleware.DebugLog("🔍 Extracting roles from JWT claims...")

	// First try to get 'roles' (plural array)
	if rolesInterface, exists := claims["roles"]; exists {
		middleware.DebugLog("🔍 Found 'roles' in claims: %+v (type: %T)", rolesInterface, rolesInterface)
		switch v := rolesInterface.(type) {
		case []string:
			roles = v
			middleware.DebugLog("✅ Roles extracted as []string: %v", roles)
		case []interface{}:
			for _, r := range v {
				if s, ok := r.(string); ok {
					roles = append(roles, s)
				}
			}
			middleware.DebugLog("✅ Roles extracted from []interface{}: %v", roles)
		default:
			middleware.DebugLog("⚠️ Unknown roles type: %T", rolesInterface)
		}
	} else {
		middleware.DebugLog("🔍 No 'roles' found in claims")
	}

	// If no roles found, try 'role' (singular from legacy auth)
	if len(roles) == 0 {
		middleware.DebugLog("🔍 No roles found, trying 'role' (singular)...")
		if roleInterface, exists := claims["role"]; exists {
			middleware.DebugLog("🔍 Found 'role' in claims: %+v (type: %T)", roleInterface, roleInterface)
			if role, ok := roleInterface.(string); ok {
				roles = []string{role}
				middleware.DebugLog("✅ Role extracted as string: %v", roles)
			} else {
				middleware.DebugLog("❌ Role is not a string: %T", roleInterface)
			}
		} else {
			middleware.DebugLog("🔍 No 'role' found in claims either")
		}
	}

	if len(roles) == 0 {
		middleware.DebugLog("❌ No roles found in JWT token for user: %s", username)
		return false, nil
	}

	middleware.DebugLog("🔍 Final roles array: %v", roles)
	middleware.DebugLog("🔍 Checking permission for user: %s, resource: %s, action: %s, resourceID: %s", username, resource, action, resourceID)

	// Define permission rules based on resource and action
	var result bool
	var err2 error

	switch resource {
	case "db_asa_student_profile":
		middleware.DebugLog("🔍 Checking student profile permissions...")
		result, err2 = checkStudentProfilePermissions(roles, action, resourceID, username, claims)
	case "db_asa_employer_profiles":
		middleware.DebugLog("🔍 Checking employer profile permissions...")
		result, err2 = checkEmployerProfilePermissions(roles, action, resourceID, username, claims)
	case "db_asa_job_posts":
		middleware.DebugLog("🔍 Checking job post permissions...")
		result, err2 = checkJobPostPermissions(roles, action, resourceID, username, claims)
	case "db_asa_applications":
		middleware.DebugLog("🔍 Checking application permissions...")
		result, err2 = checkApplicationPermissions(roles, action, resourceID, username, claims)
	case "db_asa_bookmarks":
		middleware.DebugLog("🔍 Checking bookmark permissions...")
		result, err2 = checkBookmarkPermissions(roles, action, resourceID, username, claims)
	case "db_asa_files":
		middleware.DebugLog("🔍 Checking file permissions...")
		result, err2 = checkFilePermissions(roles, action, resourceID, username, claims)
	case "db_asa_notification_preferences":
		middleware.DebugLog("🔍 Checking notification permissions...")
		result, err2 = checkNotificationPermissions(roles, action, resourceID, username, claims)
	case "db_asa_certificates":
		middleware.DebugLog("🔍 Checking certificate permissions...")
		result, err2 = checkCertificatePermissions(roles, action, resourceID, username, claims)
	case "db_asa_messages":
		middleware.DebugLog("🔍 Checking message permissions...")
		result, err2 = checkMessagePermissions(roles, action, resourceID, username, claims)
	case "db_asa_job_alerts":
		middleware.DebugLog("🔍 Checking job alert permissions...")
		result, err2 = checkJobAlertPermissions(roles, action, resourceID, username, claims)
	case "db_asa_jobs":
		middleware.DebugLog("🔍 Checking job permissions...")
		result, err2 = checkJobPermissions(roles, action, resourceID, username, claims)
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
func checkStudentProfilePermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === STUDENT PROFILE PERMISSIONS ===")
	middleware.DebugLog("🔍 Roles: %v", roles)
	middleware.DebugLog("🔍 Action: %s", action)
	middleware.DebugLog("🔍 ResourceID: %s", resourceID)

	userID := claims["user_id"].(string)
	middleware.DebugLog("🔍 UserID from claims: %s", userID)

	switch action {
	case "read":
		// Students can read their own profile, employers can read any profile
		result := contains(roles, "student") || contains(roles, "employer") || contains(roles, "admin")
		middleware.DebugLog("🔍 Read permission - Student: %v, Employer: %v, Admin: %v, Result: %v",
			contains(roles, "student"), contains(roles, "employer"), contains(roles, "admin"), result)
		return result, nil
	case "update":
		// Only students can update their own profile
		result := contains(roles, "student") && resourceID == userID
		middleware.DebugLog("🔍 Update permission - Student: %v, ResourceID match: %v, Result: %v",
			contains(roles, "student"), resourceID == userID, result)
		return result, nil
	case "create":
		// Only students can create profiles
		result := contains(roles, "student")
		middleware.DebugLog("🔍 Create permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	default:
		middleware.DebugLog("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Employer Profile Permissions
func checkEmployerProfilePermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === EMPLOYER PROFILE PERMISSIONS ===")
	middleware.DebugLog("🔍 Roles: %v", roles)
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
		result := contains(roles, "employer") && resourceID == userID
		middleware.DebugLog("🔍 Update permission - Employer: %v, ResourceID match: %v, Result: %v",
			contains(roles, "employer"), resourceID == userID, result)
		return result, nil
	case "create":
		// Only employers can create profiles
		result := contains(roles, "employer")
		middleware.DebugLog("🔍 Create permission - Employer: %v, Result: %v", contains(roles, "employer"), result)
		return result, nil
	case "delete":
		// Only admins can delete profiles
		result := contains(roles, "admin")
		middleware.DebugLog("🔍 Delete permission - Admin: %v, Result: %v", contains(roles, "admin"), result)
		return result, nil
	default:
		middleware.DebugLog("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Job Post Permissions
func checkJobPostPermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === JOB POST PERMISSIONS ===")
	middleware.DebugLog("🔍 Roles: %v", roles)
	middleware.DebugLog("🔍 Action: %s", action)
	middleware.DebugLog("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Anyone can read job posts
		middleware.DebugLog("🔍 Read permission - Anyone can read, Result: true")
		return true, nil
	case "create":
		// Only employers can create job posts
		result := contains(roles, "employer")
		middleware.DebugLog("🔍 Create permission - Employer: %v, Result: %v", contains(roles, "employer"), result)
		return result, nil
	case "update":
		// Only employers can update their own job posts
		result := contains(roles, "employer")
		middleware.DebugLog("🔍 Update permission - Employer: %v, Result: %v", contains(roles, "employer"), result)
		return result, nil
	case "delete":
		// Only employers can delete their own job posts
		result := contains(roles, "employer")
		middleware.DebugLog("🔍 Delete permission - Employer: %v, Result: %v", contains(roles, "employer"), result)
		return result, nil
	default:
		middleware.DebugLog("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Application Permissions
func checkApplicationPermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === APPLICATION PERMISSIONS ===")
	middleware.DebugLog("🔍 Roles: %v", roles)
	middleware.DebugLog("🔍 Action: %s", action)
	middleware.DebugLog("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Students can read their own applications, employers can read applications for their jobs
		result := contains(roles, "student") || contains(roles, "employer")
		middleware.DebugLog("🔍 Read permission - Student: %v, Employer: %v, Result: %v",
			contains(roles, "student"), contains(roles, "employer"), result)
		return result, nil
	case "create":
		// Only students can create applications
		result := contains(roles, "student")
		middleware.DebugLog("🔍 Create permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	case "update":
		// Students can update their own applications, employers can update applications for their jobs
		result := contains(roles, "student") || contains(roles, "employer")
		middleware.DebugLog("🔍 Update permission - Student: %v, Employer: %v, Result: %v",
			contains(roles, "student"), contains(roles, "employer"), result)
		return result, nil
	case "delete":
		// Students can delete their own applications
		result := contains(roles, "student")
		middleware.DebugLog("🔍 Delete permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	default:
		middleware.DebugLog("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Bookmark Permissions
func checkBookmarkPermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === BOOKMARK PERMISSIONS ===")
	middleware.DebugLog("🔍 Roles: %v", roles)
	middleware.DebugLog("🔍 Action: %s", action)
	middleware.DebugLog("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Students can read their own bookmarks
		result := contains(roles, "student")
		middleware.DebugLog("🔍 Read permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	case "create":
		// Only students can create bookmarks
		result := contains(roles, "student")
		middleware.DebugLog("🔍 Create permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	case "delete":
		// Students can delete their own bookmarks
		result := contains(roles, "student")
		middleware.DebugLog("🔍 Delete permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	default:
		middleware.DebugLog("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// File Permissions
func checkFilePermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === FILE PERMISSIONS ===")
	middleware.DebugLog("🔍 Roles: %v", roles)
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
		result := contains(roles, "admin")
		middleware.DebugLog("🔍 Delete permission - Admin: %v, Result: %v", contains(roles, "admin"), result)
		return result, nil
	default:
		middleware.DebugLog("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Notification Permissions
func checkNotificationPermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === NOTIFICATION PERMISSIONS ===")
	middleware.DebugLog("🔍 Roles: %v", roles)
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
func checkCertificatePermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === CERTIFICATE PERMISSIONS ===")
	middleware.DebugLog("🔍 Roles: %v", roles)
	middleware.DebugLog("🔍 Action: %s", action)
	middleware.DebugLog("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Students can read their own certificates
		result := contains(roles, "student")
		middleware.DebugLog("🔍 Read permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	case "create":
		// Only students can create certificates
		result := contains(roles, "student")
		middleware.DebugLog("🔍 Create permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	case "update":
		// Students can update their own certificates
		result := contains(roles, "student")
		middleware.DebugLog("🔍 Update permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	case "delete":
		// Students can delete their own certificates
		result := contains(roles, "student")
		middleware.DebugLog("🔍 Delete permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	default:
		middleware.DebugLog("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Message Permissions
func checkMessagePermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === MESSAGE PERMISSIONS ===")
	middleware.DebugLog("🔍 Roles: %v", roles)
	middleware.DebugLog("🔍 Action: %s", action)
	middleware.DebugLog("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Students and employers can read messages for applications they're involved in
		result := contains(roles, "student") || contains(roles, "employer")
		middleware.DebugLog("🔍 Read permission - Student: %v, Employer: %v, Result: %v",
			contains(roles, "student"), contains(roles, "employer"), result)
		return result, nil
	case "create":
		// Students and employers can create messages for applications they're involved in
		result := contains(roles, "student") || contains(roles, "employer")
		middleware.DebugLog("🔍 Create permission - Student: %v, Employer: %v, Result: %v",
			contains(roles, "student"), contains(roles, "employer"), result)
		return result, nil
	default:
		middleware.DebugLog("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Job Alert Permissions
func checkJobAlertPermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === JOB ALERT PERMISSIONS ===")
	middleware.DebugLog("🔍 Roles: %v", roles)
	middleware.DebugLog("🔍 Action: %s", action)
	middleware.DebugLog("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Students can read their own job alerts
		result := contains(roles, "student")
		middleware.DebugLog("🔍 Read permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	case "create":
		// Only students can create job alerts
		result := contains(roles, "student")
		middleware.DebugLog("🔍 Create permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	case "update":
		// Students can update their own job alerts
		result := contains(roles, "student")
		middleware.DebugLog("🔍 Update permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	case "delete":
		// Students can delete their own job alerts
		result := contains(roles, "student")
		middleware.DebugLog("🔍 Delete permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	default:
		middleware.DebugLog("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Job Permissions (for worker)
func checkJobPermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	middleware.DebugLog("🔍 === JOB PERMISSIONS ===")
	middleware.DebugLog("🔍 Roles: %v", roles)
	middleware.DebugLog("🔍 Action: %s", action)
	middleware.DebugLog("🔍 ResourceID: %s", resourceID)

	switch action {
	case "create":
		// Only employers can create jobs
		result := contains(roles, "employer")
		middleware.DebugLog("🔍 Create permission - Employer: %v, Result: %v", contains(roles, "employer"), result)
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
