# ASA Job Portal Backend - Technical Documentation

## Overview

ASA Job Portal Backend is a Go-based REST API implementing a comprehensive job portal for agricultural students and employers. The system employs JWT-based authentication with role-based access control (RBAC), PostgreSQL database with GORM ORM, and implements local authentication and authorization without external dependencies.

## Architecture Overview

### Service Layer Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   API Gateway   │    │   ASA Backend   │
│   (React/Vue)   │◄──►│   (Optional)    │◄──►│   (Go/Gin)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                       │
                                                       ▼
                                              ┌─────────────────┐
                                              │   PostgreSQL    │
                                              │   Database      │
                                              └─────────────────┘
```

### Core Components

#### 1. Local Authentication & Authorization
The system implements a complete local authentication and authorization system:

**Local Authentication Implementation:**
- **Location:** `internal/auth/` and `pkg/authz/authz.go`
- **Purpose:** Complete user authentication and authorization without external dependencies
- **Implementation:** Custom permission checking based on JWT claims and resource/action mapping

**Authentication Features:**
- User registration and login
- Password hashing with bcrypt
- JWT token generation and validation
- Role-based access control (student, employer, admin)
- Password reset functionality (mock implementation)

**Authorization System:**
```go
type AuthService interface {
    Signup(req *SignupRequest) (*SignupResponse, error)
    Login(username, password string) (*LoginResponse, error)
    GetUserByID(userID string) (*User, error)
    UpdateProfile(userID string, req *UpdateProfileRequest) error
    SendResetLink(email string) error
    ResetPassword(token, newPassword string) error
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

# JWT Configuration
SECRET_KEY=your_jwt_secret_key             # JWT signing secret

# Server Configuration
PORT=3333
GIN_MODE=debug                         # Set to 'release' for production

# CORS Configuration
CORS_ALLOW_ORIGINS=your_cors_urls                   # Comma-separated list of allowed origins

# Job Queue Configuration
JOB_MAX_RETRIES=3                                   # Maximum retry attempts for failed jobs
```

### Configuration Management
**Location:** `config/config.go`
**Implementation:** Environment variable loading with default values and validation

## API Endpoints Architecture

### Authentication Endpoints
```go
POST   /api/auth/signup          // User registration
POST   /api/auth/login           // User authentication
POST   /api/auth/forgot-password // Request password reset
POST   /api/auth/reset-password  // Reset password with token
GET    /api/auth/profile         // Get current user profile
PUT    /api/auth/profile         // Update current user profile
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
5. Set up file storage with CDN integration
6. Configure strong JWT secret

### Security Considerations
- JWT token expiration and refresh mechanisms
- File upload validation and sanitization
- SQL injection prevention via GORM
- CORS configuration for production domains
- Rate limiting implementation
- Input validation and sanitization
- Password hashing with bcrypt

## 🏭 Production Features

### Background Job Queue
The application includes a production-ready job queue system with:

#### ✅ **Features Implemented**
- **Persistent Job Storage**: Jobs are stored in Redis with full metadata
- **Job Retry Mechanism**: Exponential backoff with configurable max retries
- **Job Status Tracking**: Real-time status updates (pending, running, completed, failed, retrying)
- **Dead Letter Queue**: Failed jobs after max retries are moved to DLQ
- **Priority Queue**: Jobs can be prioritized for processing
- **Job Monitoring**: Queue statistics and health monitoring

#### 🔧 **Configuration**
```bash
# Redis Configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=your-redis-password
REDIS_DB=0
```

#### 📊 **Job Queue API**
```bash
# Get queue statistics
GET /api/worker/stats

# Get failed jobs
GET /api/worker/failed

# Retry failed job
POST /api/worker/retry/{job_id}
```

### Security Enhancements

#### ✅ **Security Features Implemented**
- **Rate Limiting**: IP-based rate limiting with configurable limits
- **Input Sanitization**: XSS and injection attack prevention
- **SQL Injection Protection**: Pattern-based detection and blocking
- **File Upload Security**: Type and size validation
- **CORS Configuration**: Strict origin validation
- **Security Headers**: Comprehensive HTTP security headers
- **Request Size Limits**: Configurable request size limits
- **Context Timeouts**: Request timeout protection

#### 🔧 **Security Configuration**
```bash
# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m

# Request Limits
MAX_REQUEST_SIZE=10485760  # 10MB
MAX_FILE_SIZE=5242880      # 5MB

# CORS
ENABLE_CORS=true
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://app.yourdomain.com
```

#### 🛡️ **Security Headers**
- `X-Frame-Options: DENY` - Prevent clickjacking
- `X-Content-Type-Options: nosniff` - Prevent MIME type sniffing
- `X-XSS-Protection: 1; mode=block` - Enable XSS protection
- `Strict-Transport-Security` - HTTPS enforcement
- `Content-Security-Policy` - Resource loading restrictions
- `Referrer-Policy` - Referrer information control
- `Permissions-Policy` - Feature permissions

### Logging & Monitoring

#### ✅ **Monitoring Features Implemented**
- **Structured Logging**: JSON-formatted logs with context
- **Performance Monitoring**: Request latency tracking and alerts
- **Error Tracking**: Comprehensive error logging with context
- **Health Checks**: Application and database health monitoring
- **Request Tracing**: Unique request IDs for tracing
- **Metrics Collection**: Basic application metrics

#### 🔧 **Logging Configuration**
```bash
# Logging
LOG_LEVEL=info
LOG_FORMAT=json
LOG_FILE=/var/log/agrijobs/app.log
LOG_DEVELOPMENT=false
```

#### 📊 **Monitoring Endpoints**
```bash
# Application health
GET /health

# Database health
GET /health/db

# Queue statistics
GET /api/worker/stats
```

#### 📈 **Performance Alerts**
- **Slow Request Detection**: Alerts for requests > 5 seconds
- **Error Rate Monitoring**: Tracks error rates and patterns
- **Queue Monitoring**: Job queue health and performance

### Production Deployment Checklist

#### ✅ **Infrastructure Requirements**
- [ ] PostgreSQL with SSL enabled
- [ ] Redis for job queue and caching
- [ ] AWS S3 for file storage
- [ ] SMTP server for email notifications
- [ ] Load balancer (optional)
- [ ] CDN for static assets (optional)

#### ✅ **Security Checklist**
- [ ] Strong JWT secret configured
- [ ] Rate limiting enabled
- [ ] CORS origins properly configured
- [ ] Security headers enabled
- [ ] File upload validation active
- [ ] Input sanitization enabled
- [ ] SQL injection protection active

#### ✅ **Monitoring Checklist**
- [ ] Structured logging configured
- [ ] Health check endpoints accessible
- [ ] Performance monitoring active
- [ ] Error tracking enabled
- [ ] Request tracing implemented
- [ ] Metrics collection active

#### ✅ **Operational Checklist**
- [ ] Database migrations run
- [ ] Environment variables configured
- [ ] File storage configured
- [ ] Email service configured
- [ ] Job queue initialized
- [ ] Application logs monitored

## Monitoring and Debugging

### Structured Logging
The system includes comprehensive structured logging:
```go
// Info logging
middleware.InfoLog("User registered successfully", userID)

// Error logging with context
middleware.ErrorLog("Database connection failed: %v", err)

// Debug logging (development only)
middleware.DebugLog("Processing request: %+v", request)
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
3. **JWT Issues:** Verify SECRET_KEY configuration
4. **CORS Issues:** Check CORS_ALLOW_ORIGINS configuration
5. **Authentication:** Verify user credentials and role assignments

## API Documentation

### Swagger Documentation
The API includes comprehensive Swagger documentation:
- **URL:** `/swagger/index.html` (when running in development)
- **Generated from:** `docs/swagger.go`
- **Auto-generated:** `docs/docs.go` and `docs/swagger.json`

### Authentication Flow
1. **Registration:** `POST /api/auth/signup`
2. **Login:** `POST /api/auth/login`
3. **Token Usage:** Include `Authorization: Bearer <token>` header
4. **Profile Access:** `GET /api/auth/profile`

### Role-Based Access
- **Student:** Can create profiles, apply for jobs, manage applications
- **Employer:** Can create job posts, manage applications, view profiles
- **Admin:** Full system access for analytics and user management

## Contributing

### Code Style
- Follow Go formatting standards (`gofmt`)
- Use meaningful variable and function names
- Include proper error handling
- Add comments for complex logic

### Testing
- Write unit tests for business logic
- Test API endpoints with proper authentication
- Verify database migrations work correctly

### Security
- Never commit sensitive data (passwords, keys)
- Validate all user inputs
- Use parameterized queries
- Implement proper error handling