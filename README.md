# ASA Backend - Technical Documentation

## Overview

ASA Backend is a Go-based microservice architecture implementing a comprehensive job portal for agricultural students and employers. The system employs JWT-based authentication with role-based access control (RBAC), PostgreSQL database with GORM ORM, and implements a custom AAA (Authentication, Authorization, and Accounting) service integration pattern.

## Architecture Overview

### Service Layer Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   API Gateway   │    │   ASA Backend   │
│   (React/Vue)   │◄──►│   (Optional)    │◄──►│   (Go/Gin)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                       │
                                                       ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   AAA Service   │    │   PostgreSQL    │
                       │   (External)    │    │   Database      │
                       └─────────────────┘    └─────────────────┘
```

### Core Components

#### 1. Authentication & Authorization (AAA) Integration
The system implements a hybrid AAA approach with both local and external service integration:

**Local AAA Implementation:**
- **Location:** `pkg/authz/authz.go`
- **Purpose:** Fallback authorization when external AAA service is unavailable
- **Implementation:** Custom permission checking based on JWT claims and resource/action mapping

**External AAA Service Integration:**
- **Configuration:** `AAA_SERVICE_URL` environment variable
- **Protocol:** HTTP REST API calls to external AAA service
- **Fallback:** Automatic fallback to local AAA when external service fails

**AAA Service Interface:**
```go
type AAAService interface {
    CheckPermission(username, resource, action, resourceID, token string) (bool, error)
}
```

**Permission Checking Flow:**
1. Extract JWT claims (user_id, roles, email)
2. Determine resource type (db_asa_student_profile, db_asa_applications, etc.)
3. Check action permissions (create, read, update, delete)
4. Validate resource ownership for user-specific operations
5. Return boolean permission result

#### 2. Database Schema & Migrations
**Migration System:**
- **Location:** `migrations/` directory
- **Execution Order:** Sequential execution via Makefile targets
- **Current Migrations:**
  - `001_complete_database_schema.sql` - Core schema
  - `007_fix_messages_timestamp.sql` - Message table fixes
  - `008_rename_user_profiles_to_student_profiles.sql` - Profile table rename
  - `009_add_phone_to_student_profiles.sql` - Phone number addition

**Key Database Tables:**
```sql
-- Core user management
users (id, email, password_hash, role, created_at, updated_at)

-- Profile management
student_profiles (id, user_id, name, email, location, phone_number, 
                 profile_photo, resume, education, portfolio, linkedin, 
                 github, experience, skills, created_at, updated_at)

employer_profiles (id, user_id, company_name, industry, company_size, 
                  website, description, logo, location, created_at, updated_at)

-- Job management
job_posts (id, employer_id, title, description, requirements, location, 
           job_type, experience_level, salary_min, salary_max, salary_currency, 
           status, created_at, updated_at)

-- Application system
applications (id, job_id, student_id, applied_at, status, cover_letter, 
             resume_file, job_title, company, location, job_type, experience, updated_at)

-- File management
certificates (id, student_profile_id, name, file, issue_date, created_at, updated_at)

-- Communication
messages (id, application_id, sender_id, receiver_id, content, sent_at, read_at)
```

#### 3. File Storage System
**Storage Architecture:**
- **Base Directory:** `uploads/`
- **Subdirectories:** `resumes/`, `certificates/`, `documents/`, `images/`
- **File Naming:** `{timestamp}_{userID}_{originalName}.{ext}`
- **Path Storage:** Relative paths stored in database
- **Validation:** File type and size validation (10MB max)

**File Upload Endpoints:**
```go
// General file upload
POST /api/upload/document/{category}

// Profile-specific uploads
POST /api/students/me/resume
POST /api/students/me/certificates/upload
POST /api/students/me/certificates/add

// Application-specific uploads
POST /api/jobs/:id/apply (includes resume upload)
```

#### 4. Middleware Stack
**Middleware Implementation:**
- **Location:** `internal/middleware/`
- **Components:**
  - `auth.go` - JWT token validation and user context injection
  - `cors.go` - Cross-origin resource sharing configuration
  - `logger.go` - Request logging and debugging

**Middleware Chain:**
```go
router.Use(middleware.CORS())
router.Use(middleware.Logger())
router.Use(middleware.Auth()) // For protected routes
```

## Environment Configuration

### Required Environment Variables
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
POSTGRES_USER=your_username
POSTGRES_PASS=your_password
DB_NAME=asa_db
DB_SSLMODE=disable

# AAA Service Configuration
AAA_SERVICE_URL=aaa_service_url  # External AAA service URL
SECRET_KEY=your_aaa_secret             # AAA service secret

# Server Configuration
PORT=3333
GIN_MODE=debug                         # Set to 'release' for production

# CORS Configuration
CORS_ALLOW_ORIGINS=your_cors_urls                   # Comma-separated list of allowed origins
```

### Configuration Management
**Location:** `config/config.go`
**Implementation:** Environment variable loading with default values and validation

## API Endpoints Architecture

### Authentication Endpoints
```go
POST   /api/auth/signup          // User registration
POST   /api/auth/login           // User authentication
POST   /api/auth/logout          // User logout
GET    /api/auth/me              // Current user info
```

### Student Profile Management
```go
GET    /api/students/me/profile              // Get current profile
PUT    /api/students/me/profile              // Update profile (multipart/form-data)
POST   /api/students/me/resume               // Upload resume
POST   /api/students/me/certificates/upload  // Upload certificate with file
POST   /api/students/me/certificates/add     // Add certificate record
```

### Job Management
```go
GET    /api/jobs                    // List jobs with filters
POST   /api/jobs                    // Create job (employer only)
GET    /api/jobs/:id                // Get job details
PUT    /api/jobs/:id                // Update job (employer only)
DELETE /api/jobs/:id                // Delete job (employer only)
POST   /api/jobs/:id/apply          // Apply for job (student only)
```

### Application Management
```go
GET    /api/applications/my                    // Student's applications
GET    /api/jobs/:id/applications              // Job applications (employer)
PUT    /api/applications/:id/status            // Update application status
```

## Development Setup

### Prerequisites
- **Go 1.21+** with modules enabled
- **PostgreSQL 12+** with proper user permissions
- **Git** for version control
- **Make** for build automation

### Installation Steps

#### 1. Repository Setup
```bash
git clone https://github.com/Kisanlink/asa-backend.git
go mod download
```

#### 2. Environment Configuration
```bash
cp .env.example .env
# Edit .env with your database and service configurations
```

#### 3. Database Initialization
```bash
# Create database
createdb asa_db

# Run migrations
make migrate

# Or reset and migrate
make migrate-reset
```

#### 4. Development Server
```bash
# With hot reloading (recommended)
make air

# Standard run
make run
```

## Build System

### Makefile Targets
```bash
# Development
make run              # Run application
make air              # Run with hot reloading
make dev              # Alias for air

# Database Management
make migrate          # Run all migrations sequentially
make migrate-reset    # Reset database and run migrations
make migrate-schema   # Apply only schema migration

# Building
make build            # Build for current platform
make build-linux      # Build for Linux
make build-windows    # Build for Windows
make build-mac        # Build for macOS

# Code Quality
make fmt              # Format code
make tidy             # Tidy Go modules

# Maintenance
make clean            # Clean build artifacts
make clean-uploads    # Clean uploaded files
make setup            # Setup development environment
```

## Deployment Considerations

### Production Configuration
1. Set `GIN_MODE=release`
2. Configure proper CORS origins
3. Use production database with SSL
4. Implement proper logging and monitoring
5. Configure external AAA service URL
6. Set up file storage with CDN integration

### Security Considerations
- JWT token expiration and refresh mechanisms
- File upload validation and sanitization
- SQL injection prevention via GORM
- CORS configuration for production domains
- Rate limiting implementation
- Input validation and sanitization

## Monitoring and Debugging

### Debug Logging
The system includes comprehensive debug logging throughout:
```go
fmt.Printf("DEBUG: Operation details - %+v\n", data)
```

### Log Levels
- **DEBUG:** Detailed operation logging
- **INFO:** General application events
- **ERROR:** Error conditions and exceptions

## Performance Optimization

### Database Optimization
- Proper indexing on frequently queried columns
- Connection pooling configuration
- Query optimization for complex joins

### File Storage Optimization
- Efficient file naming and organization
- Proper directory structure
- File size validation and limits

### API Performance
- Response caching where appropriate
- Pagination for large datasets
- Efficient database queries with GORM

## Troubleshooting

### Common Issues
1. **Database Connection:** Verify PostgreSQL service and credentials
2. **File Uploads:** Check directory permissions and storage space
3. **AAA Service:** Verify external service availability and configuration
4. **CORS Issues:** Check CORS_ALLOW_ORIGINS configuration