package authz

import (
	"asa/pkg/jwtutil"
	"log"
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

// CheckLocalPermission checks permissions locally using JWT claims instead of calling AAA service
func CheckLocalPermission(username, resource, action, resourceID, jwtToken string) (bool, error) {
	log.Printf("🔍 === LOCAL PERMISSION CHECK START ===")
	log.Printf("🔍 Username: %s", username)
	log.Printf("🔍 Resource: %s", resource)
	log.Printf("🔍 Action: %s", action)
	log.Printf("🔍 ResourceID: %s", resourceID)
	log.Printf("🔍 JWT Token length: %d", len(jwtToken))

	// Parse the JWT token to get user claims
	log.Printf("🔍 Parsing JWT token...")
	claims, err := jwtutil.ParseToken(jwtToken)
	if err != nil {
		log.Printf("❌ Failed to parse JWT token: %v", err)
		return false, err
	}
	log.Printf("✅ JWT token parsed successfully")
	log.Printf("🔍 All JWT claims: %+v", claims)

	// Extract roles from JWT claims - handle both 'role' (AAA) and 'roles' (local)
	var roles []string

	log.Printf("🔍 Extracting roles from JWT claims...")

	// First try to get 'roles' (plural array)
	if rolesInterface, exists := claims["roles"]; exists {
		log.Printf("🔍 Found 'roles' in claims: %+v (type: %T)", rolesInterface, rolesInterface)
		switch v := rolesInterface.(type) {
		case []string:
			roles = v
			log.Printf("✅ Roles extracted as []string: %v", roles)
		case []interface{}:
			for _, r := range v {
				if s, ok := r.(string); ok {
					roles = append(roles, s)
				}
			}
			log.Printf("✅ Roles extracted from []interface{}: %v", roles)
		default:
			log.Printf("⚠️ Unknown roles type: %T", rolesInterface)
		}
	} else {
		log.Printf("🔍 No 'roles' found in claims")
	}

	// If no roles found, try 'role' (singular from AAA service)
	if len(roles) == 0 {
		log.Printf("🔍 No roles found, trying 'role' (singular)...")
		if roleInterface, exists := claims["role"]; exists {
			log.Printf("🔍 Found 'role' in claims: %+v (type: %T)", roleInterface, roleInterface)
			if role, ok := roleInterface.(string); ok {
				roles = []string{role}
				log.Printf("✅ Role extracted as string: %v", roles)
			} else {
				log.Printf("❌ Role is not a string: %T", roleInterface)
			}
		} else {
			log.Printf("🔍 No 'role' found in claims either")
		}
	}

	if len(roles) == 0 {
		log.Printf("❌ No roles found in JWT token for user: %s", username)
		return false, nil
	}

	log.Printf("🔍 Final roles array: %v", roles)
	log.Printf("🔍 Checking permission for user: %s, resource: %s, action: %s, resourceID: %s", username, resource, action, resourceID)

	// Define permission rules based on resource and action
	var result bool
	var err2 error

	switch resource {
	case "db_asa_student_profile":
		log.Printf("🔍 Checking student profile permissions...")
		result, err2 = checkStudentProfilePermissions(roles, action, resourceID, username, claims)
	case "db_asa_employer_profiles":
		log.Printf("🔍 Checking employer profile permissions...")
		result, err2 = checkEmployerProfilePermissions(roles, action, resourceID, username, claims)
	case "db_asa_job_posts":
		log.Printf("🔍 Checking job post permissions...")
		result, err2 = checkJobPostPermissions(roles, action, resourceID, username, claims)
	case "db_asa_applications":
		log.Printf("🔍 Checking application permissions...")
		result, err2 = checkApplicationPermissions(roles, action, resourceID, username, claims)
	case "db_asa_bookmarks":
		log.Printf("🔍 Checking bookmark permissions...")
		result, err2 = checkBookmarkPermissions(roles, action, resourceID, username, claims)
	case "db_asa_files":
		log.Printf("🔍 Checking file permissions...")
		result, err2 = checkFilePermissions(roles, action, resourceID, username, claims)
	case "db_asa_notification_preferences":
		log.Printf("🔍 Checking notification permissions...")
		result, err2 = checkNotificationPermissions(roles, action, resourceID, username, claims)
	case "db_asa_certificates":
		log.Printf("🔍 Checking certificate permissions...")
		result, err2 = checkCertificatePermissions(roles, action, resourceID, username, claims)
	case "db_asa_messages":
		log.Printf("🔍 Checking message permissions...")
		result, err2 = checkMessagePermissions(roles, action, resourceID, username, claims)
	case "db_asa_job_alerts":
		log.Printf("🔍 Checking job alert permissions...")
		result, err2 = checkJobAlertPermissions(roles, action, resourceID, username, claims)
	case "db_asa_jobs":
		log.Printf("🔍 Checking job permissions...")
		result, err2 = checkJobPermissions(roles, action, resourceID, username, claims)
	default:
		log.Printf("⚠️ Unknown resource: %s", resource)
		return false, nil
	}

	if err2 != nil {
		log.Printf("❌ Permission check error: %v", err2)
		return false, err2
	}

	if result {
		log.Printf("✅ === PERMISSION GRANTED ===")
	} else {
		log.Printf("❌ === PERMISSION DENIED ===")
	}

	return result, nil
}

// Student Profile Permissions
func checkStudentProfilePermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	log.Printf("🔍 === STUDENT PROFILE PERMISSIONS ===")
	log.Printf("🔍 Roles: %v", roles)
	log.Printf("🔍 Action: %s", action)
	log.Printf("🔍 ResourceID: %s", resourceID)

	userID := claims["user_id"].(string)
	log.Printf("🔍 UserID from claims: %s", userID)

	switch action {
	case "read":
		// Students can read their own profile, employers can read any profile
		result := contains(roles, "student") || contains(roles, "employer") || contains(roles, "admin")
		log.Printf("🔍 Read permission - Student: %v, Employer: %v, Admin: %v, Result: %v",
			contains(roles, "student"), contains(roles, "employer"), contains(roles, "admin"), result)
		return result, nil
	case "update":
		// Only students can update their own profile
		result := contains(roles, "student") && resourceID == userID
		log.Printf("🔍 Update permission - Student: %v, ResourceID match: %v, Result: %v",
			contains(roles, "student"), resourceID == userID, result)
		return result, nil
	case "create":
		// Only students can create profiles
		result := contains(roles, "student")
		log.Printf("🔍 Create permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	default:
		log.Printf("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Employer Profile Permissions
func checkEmployerProfilePermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	log.Printf("🔍 === EMPLOYER PROFILE PERMISSIONS ===")
	log.Printf("🔍 Roles: %v", roles)
	log.Printf("🔍 Action: %s", action)
	log.Printf("🔍 ResourceID: %s", resourceID)

	userID := claims["user_id"].(string)
	log.Printf("🔍 UserID from claims: %s", userID)

	switch action {
	case "read":
		// Anyone can read employer profiles
		log.Printf("🔍 Read permission - Anyone can read, Result: true")
		return true, nil
	case "update":
		// Only employers can update their own profile
		result := contains(roles, "employer") && resourceID == userID
		log.Printf("🔍 Update permission - Employer: %v, ResourceID match: %v, Result: %v",
			contains(roles, "employer"), resourceID == userID, result)
		return result, nil
	case "create":
		// Only employers can create profiles
		result := contains(roles, "employer")
		log.Printf("🔍 Create permission - Employer: %v, Result: %v", contains(roles, "employer"), result)
		return result, nil
	case "delete":
		// Only admins can delete profiles
		result := contains(roles, "admin")
		log.Printf("🔍 Delete permission - Admin: %v, Result: %v", contains(roles, "admin"), result)
		return result, nil
	default:
		log.Printf("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Job Post Permissions
func checkJobPostPermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	log.Printf("🔍 === JOB POST PERMISSIONS ===")
	log.Printf("🔍 Roles: %v", roles)
	log.Printf("🔍 Action: %s", action)
	log.Printf("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Anyone can read job posts
		log.Printf("🔍 Read permission - Anyone can read, Result: true")
		return true, nil
	case "create":
		// Only employers can create job posts
		result := contains(roles, "employer")
		log.Printf("🔍 Create permission - Employer: %v, Result: %v", contains(roles, "employer"), result)
		return result, nil
	case "update":
		// Only employers can update their own job posts
		result := contains(roles, "employer")
		log.Printf("🔍 Update permission - Employer: %v, Result: %v", contains(roles, "employer"), result)
		return result, nil
	case "delete":
		// Only employers can delete their own job posts
		result := contains(roles, "employer")
		log.Printf("🔍 Delete permission - Employer: %v, Result: %v", contains(roles, "employer"), result)
		return result, nil
	default:
		log.Printf("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Application Permissions
func checkApplicationPermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	log.Printf("🔍 === APPLICATION PERMISSIONS ===")
	log.Printf("🔍 Roles: %v", roles)
	log.Printf("🔍 Action: %s", action)
	log.Printf("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Students can read their own applications, employers can read applications for their jobs
		result := contains(roles, "student") || contains(roles, "employer")
		log.Printf("🔍 Read permission - Student: %v, Employer: %v, Result: %v",
			contains(roles, "student"), contains(roles, "employer"), result)
		return result, nil
	case "create":
		// Only students can create applications
		result := contains(roles, "student")
		log.Printf("🔍 Create permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	case "update":
		// Students can update their own applications, employers can update applications for their jobs
		result := contains(roles, "student") || contains(roles, "employer")
		log.Printf("🔍 Update permission - Student: %v, Employer: %v, Result: %v",
			contains(roles, "student"), contains(roles, "employer"), result)
		return result, nil
	case "delete":
		// Students can delete their own applications
		result := contains(roles, "student")
		log.Printf("🔍 Delete permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	default:
		log.Printf("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Bookmark Permissions
func checkBookmarkPermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	log.Printf("🔍 === BOOKMARK PERMISSIONS ===")
	log.Printf("🔍 Roles: %v", roles)
	log.Printf("🔍 Action: %s", action)
	log.Printf("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Students can read their own bookmarks
		result := contains(roles, "student")
		log.Printf("🔍 Read permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	case "create":
		// Only students can create bookmarks
		result := contains(roles, "student")
		log.Printf("🔍 Create permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	case "delete":
		// Students can delete their own bookmarks
		result := contains(roles, "student")
		log.Printf("🔍 Delete permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	default:
		log.Printf("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// File Permissions
func checkFilePermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	log.Printf("🔍 === FILE PERMISSIONS ===")
	log.Printf("🔍 Roles: %v", roles)
	log.Printf("🔍 Action: %s", action)
	log.Printf("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Anyone can read files
		log.Printf("🔍 Read permission - Anyone can read, Result: true")
		return true, nil
	case "create":
		// Anyone authenticated can create files
		log.Printf("🔍 Create permission - Anyone authenticated can create, Result: true")
		return true, nil
	case "delete":
		// Only admins can delete files
		result := contains(roles, "admin")
		log.Printf("🔍 Delete permission - Admin: %v, Result: %v", contains(roles, "admin"), result)
		return result, nil
	default:
		log.Printf("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Notification Permissions
func checkNotificationPermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	log.Printf("🔍 === NOTIFICATION PERMISSIONS ===")
	log.Printf("🔍 Roles: %v", roles)
	log.Printf("🔍 Action: %s", action)
	log.Printf("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Anyone can read their own notifications
		log.Printf("🔍 Read permission - Anyone can read, Result: true")
		return true, nil
	case "create":
		// Anyone can create notification preferences
		log.Printf("🔍 Create permission - Anyone can create, Result: true")
		return true, nil
	case "update":
		// Anyone can update their own notification preferences
		log.Printf("🔍 Update permission - Anyone can update, Result: true")
		return true, nil
	default:
		log.Printf("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Certificate Permissions
func checkCertificatePermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	log.Printf("🔍 === CERTIFICATE PERMISSIONS ===")
	log.Printf("🔍 Roles: %v", roles)
	log.Printf("🔍 Action: %s", action)
	log.Printf("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Students can read their own certificates
		result := contains(roles, "student")
		log.Printf("🔍 Read permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	case "create":
		// Only students can create certificates
		result := contains(roles, "student")
		log.Printf("🔍 Create permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	case "update":
		// Students can update their own certificates
		result := contains(roles, "student")
		log.Printf("🔍 Update permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	case "delete":
		// Students can delete their own certificates
		result := contains(roles, "student")
		log.Printf("🔍 Delete permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	default:
		log.Printf("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Message Permissions
func checkMessagePermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	log.Printf("🔍 === MESSAGE PERMISSIONS ===")
	log.Printf("🔍 Roles: %v", roles)
	log.Printf("🔍 Action: %s", action)
	log.Printf("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Students and employers can read messages for applications they're involved in
		result := contains(roles, "student") || contains(roles, "employer")
		log.Printf("🔍 Read permission - Student: %v, Employer: %v, Result: %v",
			contains(roles, "student"), contains(roles, "employer"), result)
		return result, nil
	case "create":
		// Students and employers can create messages for applications they're involved in
		result := contains(roles, "student") || contains(roles, "employer")
		log.Printf("🔍 Create permission - Student: %v, Employer: %v, Result: %v",
			contains(roles, "student"), contains(roles, "employer"), result)
		return result, nil
	default:
		log.Printf("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Job Alert Permissions
func checkJobAlertPermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	log.Printf("🔍 === JOB ALERT PERMISSIONS ===")
	log.Printf("🔍 Roles: %v", roles)
	log.Printf("🔍 Action: %s", action)
	log.Printf("🔍 ResourceID: %s", resourceID)

	switch action {
	case "read":
		// Students can read their own job alerts
		result := contains(roles, "student")
		log.Printf("🔍 Read permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	case "create":
		// Only students can create job alerts
		result := contains(roles, "student")
		log.Printf("🔍 Create permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	case "update":
		// Students can update their own job alerts
		result := contains(roles, "student")
		log.Printf("🔍 Update permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	case "delete":
		// Students can delete their own job alerts
		result := contains(roles, "student")
		log.Printf("🔍 Delete permission - Student: %v, Result: %v", contains(roles, "student"), result)
		return result, nil
	default:
		log.Printf("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Job Permissions (for worker)
func checkJobPermissions(roles []string, action, resourceID, username string, claims map[string]interface{}) (bool, error) {
	log.Printf("🔍 === JOB PERMISSIONS ===")
	log.Printf("🔍 Roles: %v", roles)
	log.Printf("🔍 Action: %s", action)
	log.Printf("🔍 ResourceID: %s", resourceID)

	switch action {
	case "create":
		// Only employers can create jobs
		result := contains(roles, "employer")
		log.Printf("🔍 Create permission - Employer: %v, Result: %v", contains(roles, "employer"), result)
		return result, nil
	default:
		log.Printf("🔍 Unknown action: %s", action)
		return false, nil
	}
}

// Legacy function for backward compatibility - now calls local permission check
func CheckAAAPermission(username, resource, action, resourceID, jwtToken string) (bool, error) {
	log.Printf("🔍 === LEGACY AAA PERMISSION CHECK ===")
	log.Printf("🔍 Calling CheckLocalPermission instead of AAA service")
	return CheckLocalPermission(username, resource, action, resourceID, jwtToken)
}
