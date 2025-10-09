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
- User registration and login (student and employer roles only)
- Password hashing with bcrypt
- JWT token generation and validation
- Role-based access control (student, employer, asa_admin)
- Password reset functionality (mock implementation)
- **Admin Account Management:** Secure admin creation through protected endpoints only

**Authorization System:**
```go
type AuthService interface {
    Signup(req *SignupRequest) (*SignupResponse, error)
    Login(username, password string) (*LoginResponse, error)
    GetUserByID(userID string) (*User, error)
    GetCompleteProfile(userID string) (map[string]interface{}, error)
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
  - `010_convert_file_storage_to_binary.sql` - File storage conversion
  - `011_convert_ids_to_uuid.sql` - UUID conversion
  - `012_convert_binary_to_s3_keys.sql` - S3 key conversion
  - `013_remove_all_binary_storage.sql` - Binary storage cleanup
  - `014_add_contact_requests_table.sql` - Contact requests
  - `015_add_username_to_users.sql` - Username field addition

**Key Database Tables:**
```sql
-- Core user management
users (id, email, username, password_hash, role, phone_number, country_code, created_at, updated_at)

-- Profile management
student_profiles (id, user_id, location, phone_number, profile_photo_key, resume_key, 
                 education, portfolio, linkedin, github, experience, skills, certificates, created_at, updated_at)

employer_profiles (id, user_id, company_name, industry, company_size, logo_key, logo_name, 
                  logo_type, logo_size, website_url, company_description, recruiter_name, 
                  designation, official_email, phone_number, linkedin_profile, job_categories, 
                  hiring_locations, hiring_types, gstin_number, company_address, city, 
                  state, pincode, created_at, updated_at)

-- Job management
job_posts (id, employer_id, title, description, requirements, location, 
           job_type, experience_level, salary_min, salary_max, salary_currency, 
           status, created_at, updated_at)

-- Application system
applications (id, job_id, student_id, applied_at, status, cover_letter, 
             resume_file, job_title, company, location, job_type, experience, updated_at)

-- File management
certificates (id, student_profile_id, name, file_key, issue_date, created_at, updated_at)

-- Communication
messages (id, application_id, sender_id, receiver_id, content, sent_at, read_at)

-- Contact requests
contact_requests (id, name, email, subject, message, created_at)
```

#### 3. File Storage System
**Storage Architecture:**
- **AWS S3 Integration:** Primary file storage with local fallback
- **File Keys:** S3 object keys stored in database
- **File Validation:** File type and size validation (configurable max size)
- **Metadata Storage:** File information stored in database with S3 keys

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
  - `admin.go` - Admin role authorization middleware
  - `cors.go` - Cross-origin resource sharing configuration
  - `logger.go` - Request logging and debugging
  - `security.go` - Security headers and rate limiting

**Middleware Chain:**
```go
router.Use(middleware.CORS())
router.Use(middleware.Logger())
router.Use(middleware.Auth()) // For protected routes
router.Use(middleware.AdminAuthMiddleware()) // For admin-only routes
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
JWT_SECRET=your_jwt_secret_key             # JWT signing secret (required)

# Server Configuration
SERVER_PORT=8080                            # Server port (required)
GIN_MODE=debug                             # Set to 'release' for production

# AWS S3 Configuration
AWS_ACCESS_KEY_ID=your_access_key          # AWS access key (required)
AWS_SECRET_ACCESS_KEY=your_secret_key      # AWS secret key (required)
AWS_REGION=your_region                     # AWS region (required)
AWS_S3_BUCKET=your_bucket_name            # S3 bucket name (required)
AWS_S3_ENDPOINT=your_s3_endpoint          # S3 endpoint (optional, for custom endpoints)
AWS_S3_FORCE_PATH_STYLE=false             # Force path style (optional)
AWS_S3_DISABLE_SSL=false                  # Disable SSL (optional)

# Application Configuration
ASA_BASE_URL=https://yourdomain.com        # Base URL for the application (required)

# Email Configuration
SMTP_HOST=your_smtp_host                   # SMTP server host (required)
SMTP_PORT=587                              # SMTP server port (required)
SMTP_USERNAME=your_smtp_username           # SMTP username (required)
SMTP_PASSWORD=your_smtp_password           # SMTP password (required)
SMTP_FROM_EMAIL=noreply@yourdomain.com     # From email address (required)

# Admin Seeding Configuration (Optional - for initial setup)
ASA_RUN_SEED=false                         # Enable admin seeding (true/false)
DEFAULT_ADMIN_EMAIL=admin@agrijobs.com     # Default admin email (required if seeding)
DEFAULT_ADMIN_PASSWORD=admin123            # Default admin password (required if seeding)
DEFAULT_ADMIN_NAME=System Administrator    # Default admin name (required if seeding)
DEFAULT_ADMIN_USERNAME=admin               # Default admin username (required if seeding)

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://app.yourdomain.com

# Rate Limiting
RATE_LIMIT_REQUESTS=100                    # Max requests per window
RATE_LIMIT_WINDOW=1m                       # Time window for rate limiting

# Request Limits
MAX_REQUEST_SIZE=10485760                   # 10MB max request size
MAX_FILE_SIZE=5242880                       # 5MB max file size

# Redis Configuration
REDIS_ADDR=localhost:6379                  # Redis server address
REDIS_PASSWORD=your-redis-password         # Redis password
REDIS_DB=0                                 # Redis database number

# Health Check
HEALTH_CHECK_TIMEOUT=30s                   # Health check timeout
```

### Configuration Management
**Location:** `config/config.go`
**Implementation:** Environment variable loading with validation and required variable enforcement

**Important Notes:**
- **Critical variables** (database, JWT, server, AWS, email) are required and will cause application failure if not set
- **Optional variables** have internal defaults for non-critical configurations
- **Admin seeding** is controlled by `ASA_RUN_SEED` environment variable
- **No hardcoded defaults** for security-sensitive configurations

## API Endpoints Architecture

### Authentication Endpoints
```go
POST   /api/auth/signup          // User registration (student/employer only)
POST   /api/auth/login           // User authentication
POST   /api/auth/forgot-password // Request password reset
POST   /api/auth/reset-password  // Reset password with token
GET    /api/auth/profile         // Get complete user profile with role-specific data
PUT    /api/auth/profile         // Update current user profile
```

### Admin Management Endpoints
```go
POST   /api/admin/create-admin   // Create new admin user (admin-only)
GET    /api/admin/dashboard      // Get admin dashboard analytics
GET    /api/admin/users          // Get user list (admin-only)
DELETE /api/admin/users/:id      // Delete user (admin-only)
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
- **Redis** for job queue and caching
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
# Ensure all required variables are set
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

#### 4. Admin Account Setup (Optional)
```bash
# Enable admin seeding in .env
ASA_RUN_SEED=true
DEFAULT_ADMIN_EMAIL=admin@agrijobs.com
DEFAULT_ADMIN_PASSWORD=admin123
DEFAULT_ADMIN_NAME=System Administrator
DEFAULT_ADMIN_USERNAME=admin

# Start the application - admin account will be created automatically
make run
```

#### 5. Development Server
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

## Security Features

### Admin Account Security
- **Public Signup Restriction:** Only "student" and "employer" roles allowed in public signup
- **Protected Admin Creation:** New admin accounts can only be created by existing admins
- **Initial Admin Seeding:** Optional seeding mechanism for initial setup with environment variable control
- **Role-Based Middleware:** Admin endpoints protected by role verification middleware

### Authentication & Authorization
- **JWT-Based Authentication:** Secure token-based authentication
- **Role-Based Access Control:** Granular permissions based on user roles
- **Password Security:** Bcrypt hashing with configurable complexity
- **Token Validation:** Comprehensive JWT token validation and expiration

### Input Validation & Sanitization
- **Request Validation:** Comprehensive input validation using binding tags
- **SQL Injection Protection:** GORM ORM provides parameterized query protection
- **File Upload Security:** File type and size validation
- **XSS Prevention:** Input sanitization and output encoding

### Security Headers & Rate Limiting
- **Security Headers:** Comprehensive HTTP security headers
- **Rate Limiting:** IP-based rate limiting with configurable limits
- **CORS Protection:** Strict origin validation
- **Request Size Limits:** Configurable request and file size limits

## Deployment Considerations

### Production Configuration
1. Set `GIN_MODE=release`
2. Configure proper CORS origins
3. Use production database with SSL
4. Implement proper logging and monitoring
5. Set up AWS S3 for file storage
6. Configure strong JWT secret
7. **Disable admin seeding** (`ASA_RUN_SEED=false`)
8. Use strong, unique passwords for all admin accounts

### Security Considerations
- JWT token expiration and refresh mechanisms
- File upload validation and sanitization
- SQL injection prevention via GORM
- CORS configuration for production domains
- Rate limiting implementation
- Input validation and sanitization
- Password hashing with bcrypt
- **Admin account management through protected endpoints only**

### Infrastructure Requirements
- PostgreSQL with SSL enabled
- Redis for job queue and caching
- AWS S3 for file storage
- SMTP server for email notifications
- Load balancer (optional)
- CDN for static assets (optional)

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
- **Admin Role Protection**: Admin-only endpoints with middleware protection

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
- [ ] Admin seeding disabled in production
- [ ] Admin endpoints protected by middleware

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
- [ ] Admin accounts properly configured

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
- Efficient S3 key management
- Proper file metadata storage
- File size validation and limits

### API Performance
- Response caching where appropriate
- Pagination for large datasets
- Efficient database queries with GORM

## Troubleshooting

### Common Issues
1. **Database Connection:** Verify PostgreSQL service and credentials
2. **File Uploads:** Check S3 configuration and permissions
3. **JWT Issues:** Verify JWT_SECRET configuration
4. **CORS Issues:** Check CORS_ALLOWED_ORIGINS configuration
5. **Authentication:** Verify user credentials and role assignments
6. **Admin Access:** Ensure admin accounts are properly created and roles assigned
7. **Environment Variables:** Verify all required variables are set

### Admin Account Issues
1. **Cannot Create Admin:** Ensure you're logged in as an existing admin
2. **Seeding Not Working:** Check `ASA_RUN_SEED` and related environment variables
3. **Role Assignment:** Verify user has `asa_admin` role in database

## API Documentation

### Swagger Documentation
The API includes comprehensive Swagger documentation:
- **URL:** `/swagger/index.html` (when running in development)
- **Generated from:** `docs/swagger.go`
- **Auto-generated:** `docs/docs.go` and `docs/swagger.json`

### Authentication Flow
1. **Registration:** `POST /api/auth/signup` (student/employer only)
2. **Login:** `POST /api/auth/login`
3. **Token Usage:** Include `Authorization: Bearer <token>` header
4. **Profile Access:** `GET /api/auth/profile` (returns complete profile data)

### Admin Account Management
1. **Initial Setup:** Use seeding with `ASA_RUN_SEED=true` (development only)
2. **Create New Admin:** `POST /api/admin/create-admin` (admin-only endpoint)
3. **Admin Login:** Use admin credentials to access admin endpoints

### Role-Based Access
- **Student:** Can create profiles, apply for jobs, manage applications
- **Employer:** Can create job posts, manage applications, view profiles
- **Admin:** Full system access for analytics, user management, and admin creation

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
- **Test admin endpoint security thoroughly**