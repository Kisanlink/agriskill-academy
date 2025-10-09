# ASA Job Portal Backend - Technical Architecture Documentation

## Table of Contents
1. [Project Overview](#project-overview)
2. [Architecture Overview](#architecture-overview)
3. [Technology Stack](#technology-stack)
4. [Project Structure](#project-structure)
5. [Database Architecture](#database-architecture)
6. [API Architecture](#api-architecture)
7. [Authentication & Authorization](#authentication--authorization)
8. [File Storage System](#file-storage-system)
9. [Middleware Components](#middleware-components)
10. [Service Layer Architecture](#service-layer-architecture)
11. [Development Workflow](#development-workflow)
12. [Deployment & Configuration](#deployment--configuration)
13. [Security Features](#security-features)
14. [Performance & Monitoring](#performance--monitoring)

---

## Project Overview

**ASA Job Portal Backend** is a comprehensive REST API built in Go that serves as the backend for an agricultural job portal connecting students with employers. The system implements a modern microservices-inspired architecture with clean separation of concerns, robust authentication, and scalable file storage.

### Key Features
- **Multi-role Authentication**: Student, Employer, and Admin roles
- **Job Management**: Complete job posting and application workflow
- **Profile Management**: Detailed student and employer profiles
- **File Storage**: AWS S3 integration for documents, images, and certificates
- **Real-time Communication**: Messaging system between employers and students
- **Admin Dashboard**: Comprehensive analytics and user management
- **Background Jobs**: Redis-based job queue for async processing

---

## Architecture Overview

### High-Level Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Load Balancer │    │   ASA Backend   │
│   (React/Vue)   │◄──►│   (Optional)    │◄──►│   (Go/Gin)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                       │
                                                       ▼
                                              ┌─────────────────┐
                                              │   PostgreSQL    │
                                              │   Database      │
                                              └─────────────────┘
                                                       │
                                                       ▼
                                              ┌─────────────────┐
                                              │   AWS S3        │
                                              │   File Storage  │
                                              └─────────────────┘
                                                       │
                                                       ▼
                                              ┌─────────────────┐
                                              │   Redis         │
                                              │   Job Queue     │
                                              └─────────────────┘
```

### Design Patterns
- **Clean Architecture**: Separation of handlers, services, and repositories
- **Dependency Injection**: Services injected into handlers
- **Repository Pattern**: Data access abstraction
- **Middleware Pattern**: Cross-cutting concerns (auth, logging, CORS)
- **Factory Pattern**: Service and handler creation

---

## Technology Stack

### Core Technologies
| Technology | Version | Purpose |
|------------|---------|---------|
| **Go** | 1.24.4 | Primary programming language |
| **Gin** | 1.10.1 | HTTP web framework |
| **GORM** | 1.30.0 | ORM for database operations |
| **PostgreSQL** | Latest | Primary database |
| **Redis** | Latest | Job queue and caching |
| **AWS S3** | Latest | File storage |

### Key Dependencies
| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/Kisanlink/kisanlink-db` | v0.1.9 | Custom database utilities and ID generation |
| `github.com/aws/aws-sdk-go` | v1.53.0 | AWS S3 integration |
| `github.com/golang-jwt/jwt/v5` | v5.2.2 | JWT token handling |
| `github.com/joho/godotenv` | v1.5.1 | Environment variable management |
| `github.com/swaggo/swag` | v1.16.4 | API documentation generation |
| `golang.org/x/crypto` | v0.39.0 | Password hashing and encryption |
| `go.uber.org/zap` | v1.24.0 | Structured logging |

### Development Tools
| Tool | Purpose |
|------|---------|
| **Air** | Hot reload for development |
| **Swagger** | API documentation |
| **Make** | Build automation |
| **PostgreSQL Client** | Database management |

---

## Project Structure

```
agrijobs/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── config/
│   └── config.go                   # Configuration management
├── internal/                       # Private application code
│   ├── admin/                      # Admin management module
│   ├── application/                # Job application module
│   ├── auth/                       # Authentication module
│   ├── bookmark/                   # Job bookmarking module
│   ├── contact/                    # Contact form module
│   ├── employerapplication/        # Employer application management
│   ├── employerprofile/            # Employer profile management
│   ├── jobpost/                    # Job posting module
│   ├── middleware/                 # HTTP middleware
│   ├── notification/               # Notification system
│   ├── storage/                    # File storage module
│   ├── studentprofile/             # Student profile management
│   └── worker/                     # Background job processing
├── pkg/                           # Public packages
│   ├── authz/                     # Authorization utilities
│   ├── db/                        # Database utilities
│   ├── jwtutil/                   # JWT utilities
│   └── mailutil/                  # Email utilities
├── migrations/                    # Database migration files
├── scripts/                       # Database migration scripts
├── docs/                          # Swagger documentation
├── uploads/                       # Local file storage (fallback)
├── go.mod                         # Go module definition
├── go.sum                         # Go module checksums
├── Makefile                       # Build automation
└── README.md                      # Project documentation
```

### Module Structure Pattern
Each module follows a consistent structure:
```
module/
├── handler.go     # HTTP request handlers
├── model.go       # Data models and structs
├── repository.go  # Data access layer
├── routes.go      # Route definitions
└── service.go     # Business logic layer
```

---

## Database Architecture

### Database Technology
- **GORM** ORM for Go integration
- **Custom ID Generation** using kisanlink-db package

### Core Tables

#### Users Table
| Column | Type | Description |
|--------|------|-------------|
| id | VARCHAR(255) | Custom generated ID (Primary Key) |
| name | TEXT | User's full name |
| username | TEXT UNIQUE | Separate username field |
| email | TEXT UNIQUE NOT NULL | User's email address |
| password | TEXT NOT NULL | Bcrypt hashed password |
| role | TEXT NOT NULL | User role (student, employer, asa_admin) |
| phone_number | TEXT | User's phone number |
| avatar_key | TEXT | S3 key for avatar image |
| created_at | TIMESTAMP WITH TIME ZONE | Record creation timestamp |
| updated_at | TIMESTAMP WITH TIME ZONE | Record update timestamp |

#### Student Profiles Table
| Column | Type | Description |
|--------|------|-------------|
| id | VARCHAR(255) | Custom generated ID (Primary Key) |
| user_id | VARCHAR(255) NOT NULL | Foreign key to users table |
| name | TEXT NOT NULL | Student's full name |
| email | TEXT NOT NULL | Student's email address |
| location | TEXT | Student's location |
| phone_number | TEXT | Student's phone number |
| profile_photo_key | TEXT | S3 key for profile photo |
| resume_key | TEXT | S3 key for resume file |
| skills | TEXT[] | PostgreSQL array of skills |
| experience | FLOAT | Years of experience |
| education | TEXT | Education details |
| portfolio | TEXT | Portfolio URL |
| linkedin | TEXT | LinkedIn profile URL |
| github | TEXT | GitHub profile URL |
| created_at | TIMESTAMP WITH TIME ZONE | Record creation timestamp |
| updated_at | TIMESTAMP WITH TIME ZONE | Record update timestamp |

#### Employer Profiles Table
| Column | Type | Description |
|--------|------|-------------|
| id | VARCHAR(255) | Custom generated ID (Primary Key) |
| user_id | VARCHAR(255) NOT NULL | Foreign key to users table |
| company_name | TEXT NOT NULL | Company name |
| industry | TEXT NOT NULL | Industry sector |
| company_size | TEXT NOT NULL | Company size category |
| logo_key | TEXT | S3 key for company logo |
| website_url | TEXT | Company website URL |
| company_description | TEXT | Company description |
| recruiter_name | TEXT | Recruiter's name |
| designation | TEXT | Recruiter's designation |
| official_email | TEXT | Official email address |
| phone_number | TEXT | Company phone number |
| linkedin_profile | TEXT | Company LinkedIn profile |
| job_categories | TEXT[] | PostgreSQL array of job categories |
| hiring_locations | TEXT[] | PostgreSQL array of hiring locations |
| hiring_types | TEXT[] | PostgreSQL array of hiring types |
| gstin_number | TEXT | GSTIN number |
| company_address | TEXT | Company address |
| city | TEXT | Company city |
| state | TEXT | Company state |
| pincode | TEXT | Company pincode |
| created_at | TIMESTAMP WITH TIME ZONE | Record creation timestamp |
| updated_at | TIMESTAMP WITH TIME ZONE | Record update timestamp |

#### Job Posts Table
| Column | Type | Description |
|--------|------|-------------|
| id | VARCHAR(255) | Custom generated ID (Primary Key) |
| employer_id | VARCHAR(255) NOT NULL | Foreign key to employer profiles |
| title | TEXT NOT NULL | Job title |
| description | TEXT NOT NULL | Job description |
| requirements | TEXT | Job requirements |
| location | TEXT | Job location |
| job_type | TEXT | Type of job (full-time, part-time, etc.) |
| experience_level | TEXT | Required experience level |
| salary_min | INTEGER | Minimum salary |
| salary_max | INTEGER | Maximum salary |
| salary_currency | TEXT | Salary currency |
| status | TEXT DEFAULT 'draft' | Job status (draft, published, closed) |
| applications_count | INTEGER DEFAULT 0 | Number of applications received |
| created_at | TIMESTAMP WITH TIME ZONE | Record creation timestamp |
| updated_at | TIMESTAMP WITH TIME ZONE | Record update timestamp |

#### Applications Table
| Column | Type | Description |
|--------|------|-------------|
| id | VARCHAR(255) | Custom generated ID (Primary Key) |
| job_id | VARCHAR(255) NOT NULL | Foreign key to job posts |
| student_id | VARCHAR(255) NOT NULL | Foreign key to student profiles |
| applied_at | TIMESTAMP WITH TIME ZONE | Application submission timestamp |
| status | TEXT DEFAULT 'applied' | Application status (applied, reviewing, shortlisted, etc.) |
| cover_letter | TEXT | Cover letter content |
| resume_key | TEXT | S3 key for application resume |
| job_title | TEXT | Job title (denormalized for performance) |
| company | TEXT | Company name (denormalized for performance) |
| location | TEXT | Job location (denormalized for performance) |
| job_type | TEXT | Job type (denormalized for performance) |
| experience | TEXT | Experience level (denormalized for performance) |
| created_at | TIMESTAMP WITH TIME ZONE | Record creation timestamp |
| updated_at | TIMESTAMP WITH TIME ZONE | Record update timestamp |

#### Certificates Table
| Column | Type | Description |
|--------|------|-------------|
| id | VARCHAR(255) | Custom generated ID (Primary Key) |
| student_profile_id | VARCHAR(255) NOT NULL | Foreign key to student profiles |
| name | TEXT NOT NULL | Certificate name |
| file_key | TEXT | S3 key for certificate file |
| issue_date | TEXT | Certificate issue date |
| created_at | TIMESTAMP WITH TIME ZONE | Record creation timestamp |
| updated_at | TIMESTAMP WITH TIME ZONE | Record update timestamp |

#### Messages Table
| Column | Type | Description |
|--------|------|-------------|
| id | VARCHAR(255) | Custom generated ID (Primary Key) |
| application_id | VARCHAR(255) NOT NULL | Foreign key to applications |
| sender_id | VARCHAR(255) NOT NULL | Foreign key to users (sender) |
| receiver_id | VARCHAR(255) NOT NULL | Foreign key to users (receiver) |
| content | TEXT NOT NULL | Message content |
| sent_at | TIMESTAMP WITH TIME ZONE | Message sent timestamp |
| read_at | TIMESTAMP WITH TIME ZONE | Message read timestamp |

#### Contact Requests Table
| Column | Type | Description |
|--------|------|-------------|
| id | VARCHAR(255) | Custom generated ID (Primary Key) |
| first_name | TEXT NOT NULL | Contact person's first name |
| last_name | TEXT NOT NULL | Contact person's last name |
| email | TEXT NOT NULL | Contact person's email |
| phone | TEXT | Contact person's phone |
| subject | TEXT NOT NULL | Contact request subject |
| message | TEXT NOT NULL | Contact request message |
| status | TEXT DEFAULT 'new' | Request status (new, in_progress, resolved) |
| created_at | TIMESTAMP WITH TIME ZONE | Record creation timestamp |
| updated_at | TIMESTAMP WITH TIME ZONE | Record update timestamp |

#### Bookmarks Table
| Column | Type | Description |
|--------|------|-------------|
| id | VARCHAR(255) | Custom generated ID (Primary Key) |
| user_id | VARCHAR(255) NOT NULL | Foreign key to users |
| job_id | VARCHAR(255) NOT NULL | Foreign key to job posts |
| created_at | TIMESTAMP WITH TIME ZONE | Record creation timestamp |
| updated_at | TIMESTAMP WITH TIME ZONE | Record update timestamp |

#### Job Alerts Table
| Column | Type | Description |
|--------|------|-------------|
| id | VARCHAR(255) | Custom generated ID (Primary Key) |
| user_id | VARCHAR(255) NOT NULL | Foreign key to users |
| title | TEXT | Alert title |
| keywords | TEXT[] | PostgreSQL array of keywords |
| location | TEXT | Preferred location |
| job_type | TEXT | Preferred job type |
| experience_level | TEXT | Preferred experience level |
| salary_min | INTEGER | Minimum salary expectation |
| is_active | BOOLEAN DEFAULT true | Alert active status |
| created_at | TIMESTAMP WITH TIME ZONE | Record creation timestamp |
| updated_at | TIMESTAMP WITH TIME ZONE | Record update timestamp |

#### Notification Preferences Table
| Column | Type | Description |
|--------|------|-------------|
| id | VARCHAR(255) | Custom generated ID (Primary Key) |
| user_id | VARCHAR(255) NOT NULL | Foreign key to users |
| email_notifications | BOOLEAN DEFAULT true | Email notification preference |
| sms_notifications | BOOLEAN DEFAULT false | SMS notification preference |
| job_alerts | BOOLEAN DEFAULT true | Job alert notifications |
| application_updates | BOOLEAN DEFAULT true | Application update notifications |
| marketing_emails | BOOLEAN DEFAULT false | Marketing email preference |
| created_at | TIMESTAMP WITH TIME ZONE | Record creation timestamp |
| updated_at | TIMESTAMP WITH TIME ZONE | Record update timestamp |


---

## API Architecture

### API Design Principles
- **RESTful Design**: Standard HTTP methods and status codes
- **JSON Communication**: All requests/responses in JSON format
- **Consistent Response Format**: Standardized success/error responses
- **Swagger Documentation**: Auto-generated API documentation

### Response Format
The API uses a consistent response format with three main fields:
- **success**: Boolean indicating operation success
- **message**: Human-readable message describing the result
- **data**: The actual response data (when successful)

### Error Format
Error responses follow a similar structure:
- **success**: Always false for error responses
- **message**: User-friendly error description
- **error**: Detailed technical error information (for debugging)

### Core API Endpoints

#### Authentication Endpoints
- **POST /api/auth/signup** - User registration
- **POST /api/auth/login** - User authentication
- **POST /api/auth/forgot-password** - Password reset request
- **POST /api/auth/reset-password** - Password reset
- **GET /api/auth/profile** - Get user profile
- **PUT /api/auth/profile** - Update user profile

#### Job Management Endpoints
- **GET /api/jobs** - List all jobs (public)
- **GET /api/jobs/featured** - Get featured jobs
- **GET /api/jobs/recent** - Get recent jobs
- **GET /api/jobs/trending** - Get trending jobs
- **GET /api/jobs/:id** - Get job details
- **POST /api/jobs** - Create job (employer)
- **PUT /api/jobs/:id** - Update job (employer)
- **DELETE /api/jobs/:id** - Delete job (employer)
- **POST /api/jobs/search** - Search jobs
- **POST /api/jobs/advanced-search** - Advanced job search

#### Student Profile Endpoints
- **GET /api/students/me/profile** - Get current profile
- **PUT /api/students/me/profile** - Update profile
- **POST /api/students/me/resume** - Upload resume
- **POST /api/students/me/certificates** - Add certificate
- **POST /api/students/me/certificates/upload** - Upload certificate file
- **DELETE /api/students/me/certificates/:id** - Delete certificate

#### Application Management Endpoints
- **GET /api/applications/my** - Student's applications
- **POST /api/jobs/:id/apply** - Apply for job
- **PUT /api/applications/:id/status** - Update application status
- **GET /api/jobs/:id/applications** - Job applications (employer)

#### Admin Endpoints
- **GET /api/admin/dashboard** - Admin dashboard
- **GET /api/admin/users** - User management
- **POST /api/admin/create-admin** - Create admin user
- **DELETE /api/admin/users/:id** - Delete user

#### File Storage Endpoints
- **POST /api/upload/:folder** - Upload files
- **GET /api/files/serve/:type/:id** - Serve files
- **DELETE /api/files/*filePath** - Delete files

---

## Authentication & Authorization

### Authentication System
The system implements **local authentication** without external dependencies:

#### JWT Token Structure
JWT tokens contain the following claims:
- **user_id**: Unique user identifier
- **username**: User's username
- **email**: User's email address
- **name**: User's full name
- **roles**: Array of user roles
- **iat**: Token issued at timestamp
- **exp**: Token expiration timestamp

#### Password Security
- **Bcrypt Hashing**: Passwords hashed with bcrypt
- **Salt Rounds**: Configurable salt rounds (default: 12)
- **No Plain Text**: Passwords never stored in plain text

#### Role-Based Access Control (RBAC)
- **student**: Can apply for jobs, manage profile
- **employer**: Can post jobs, manage applications
- **asa_admin**: Full system access

### Authorization Flow
1. **Token Validation**: JWT token parsed and validated
2. **Role Extraction**: User roles extracted from token
3. **Permission Check**: Resource-specific permissions validated
4. **Context Setting**: User information set in request context

### Permission System
The permission system uses a resource-action model where:
- **User identifier**: Username or user ID
- **Resource type**: Database table or resource (e.g., "db_asa_student_profile")
- **Action**: Operation type (create, read, update, delete)
- **Resource ID**: Specific resource identifier
- **JWT token**: Authentication token for validation

---

## File Storage System

### Storage Architecture
- **Primary**: AWS S3 for production
- **Fallback**: Local file system for development
- **Metadata**: File information stored in database

### File Types Supported
- **Images**: JPG, PNG, GIF, WebP (max 5MB)
- **Documents**: PDF, DOC, DOCX (max 10MB)
- **Resumes**: PDF, DOC, DOCX (max 10MB)
- **Certificates**: PDF, DOC, DOCX, JPG, JPEG, PNG (max 10MB)

### S3 Integration
The S3 integration includes configuration for:
- **S3Region**: AWS region for the S3 bucket
- **S3Bucket**: Name of the S3 bucket
- **S3Endpoint**: Custom S3 endpoint (for S3-compatible services)
- **S3ForcePathStyle**: Path-style addressing configuration
- **S3DisableSSL**: SSL configuration for development
- **S3AccessKeyID**: AWS access key ID
- **S3SecretAccessKey**: AWS secret access key

### File Upload Flow
1. **Validation**: File type and size validation
2. **Upload**: File uploaded to S3
3. **Key Storage**: S3 key stored in database
4. **Metadata**: File metadata stored in database
5. **Response**: File URL returned to client

### File Serving
Files are served through dedicated endpoints:
- **GET /api/files/serve/resume/:user_id** - Serve user resume
- **GET /api/files/serve/certificate/:certificate_id** - Serve certificate file
- **GET /api/files/serve/profile-photo/:user_id** - Serve profile photo
- **GET /api/files/serve/logo/:employer_id** - Serve company logo

---

## Middleware Components

### Core Middleware Stack
The application uses a comprehensive middleware stack in the following order:
1. **RequestIDMiddleware** - Request tracking with unique IDs
2. **StructuredLoggingMiddleware** - JSON-formatted logging
3. **PerformanceMonitoringMiddleware** - Response time tracking
4. **ErrorLoggingMiddleware** - Error handling and logging
5. **SecurityHeadersMiddleware** - Security headers injection
6. **RateLimitMiddleware** - Request rate limiting
7. **InputSanitizationMiddleware** - Input validation and sanitization
8. **SQLInjectionProtectionMiddleware** - SQL injection prevention
9. **RequestSizeLimitMiddleware** - Request size limits
10. **ContextTimeoutMiddleware** - Request timeout management
11. **CORSMiddleware** - Cross-origin resource sharing
12. **AuthMiddleware** - Authentication and authorization

### Authentication Middleware
- **JWT Validation**: Token parsing and validation
- **Role Extraction**: User roles from token claims
- **Context Setting**: User information in request context
- **Error Handling**: Comprehensive error responses

### Security Middleware
- **Rate Limiting**: Configurable request limits
- **Input Sanitization**: XSS and injection protection
- **Security Headers**: HSTS, CSP, X-Frame-Options
- **Request Size Limits**: Prevent large payload attacks

### Logging Middleware
- **Structured Logging**: JSON-formatted logs
- **Request Tracking**: Unique request IDs
- **Performance Metrics**: Response time tracking
- **Error Logging**: Detailed error information

---

## Service Layer Architecture

### Service Interface Pattern
Each service implements a consistent interface with standard CRUD operations:
- **Create**: Create new entities with validation
- **GetByID**: Retrieve entities by unique identifier
- **Update**: Update existing entities with validation
- **Delete**: Remove entities with proper cleanup
- **List**: Retrieve multiple entities with filtering and pagination

### Key Services

#### AuthService
- User registration and login
- JWT token generation
- Password management
- Profile management

#### JobPostService
- Job creation and management
- Job search and filtering
- Job recommendations
- Job alerts

#### ApplicationService
- Job application processing
- Application status management
- Resume handling
- Application analytics

#### StudentProfileService
- Profile management
- Certificate handling
- Skill management
- File uploads

#### EmployerProfileService
- Company profile management
- Logo handling
- Hiring preferences
- Company analytics

### Repository Pattern
Data access is abstracted through repositories with standard operations:
- **Create**: Insert new entities into the database
- **GetByID**: Retrieve entities by unique identifier
- **Update**: Modify existing entities in the database
- **Delete**: Remove entities from the database
- **List**: Query multiple entities with filtering capabilities

---

## Development Workflow

### Build System
The project uses a comprehensive Makefile for build automation:

**Development Commands:**
- **make run** - Build and run the application
- **make air** - Run with hot reload
- **make build** - Build the application

**Database Commands:**
- **make migrate** - Apply all database migrations
- **make migrate-reset** - Reset database and apply migrations

**Utility Commands:**
- **make test** - Run tests
- **make clean** - Clean build artifacts
- **make setup** - Initial setup

### Hot Reload Development
The project supports hot reload development using Air:
- Install Air: `go install github.com/cosmtrek/air@latest`
- Run with hot reload: `make air`

### Database Management
Database operations are handled through Makefile commands:
- **make migrate** - Apply all database migrations
- **make migrate-reset** - Reset database and apply migrations
- **make debug-migration** - Debug migration issues

### Testing
The project supports comprehensive testing:
- **make test** - Run all tests
- **go test ./internal/auth/...** - Run specific package tests

---

## Deployment & Configuration

### Environment Variables
The application uses environment variables for configuration:

#### Required Variables
**Database Configuration:**
- DB_HOST - Database host address
- DB_PORT - Database port number
- DB_NAME - Database name
- POSTGRES_USER - Database username
- POSTGRES_PASS - Database password

**JWT Configuration:**
- JWT_SECRET - Secret key for JWT token signing

**Server Configuration:**
- SERVER_PORT - Application server port
- ASA_BASE_URL - Base URL for the application

**AWS S3 Configuration:**
- AWS_REGION - AWS region for S3
- AWS_S3_BUCKET - S3 bucket name
- AWS_ACCESS_KEY_ID - AWS access key
- AWS_SECRET_ACCESS_KEY - AWS secret key

**Redis Configuration:**
- REDIS_ADDR - Redis server address
- REDIS_PASSWORD - Redis password
- REDIS_DB - Redis database number

#### Optional Variables
**Email Configuration:**
- MAIL_FROM - Sender email address
- MAIL_HOST - SMTP server host
- MAIL_PORT - SMTP server port
- MAIL_PASS - Email password

**Logging Configuration:**
- LOG_LEVEL - Logging level (debug, info, warn, error)
- LOG_OUTPUT_PATH - Log file path
- LOG_FORMAT - Log format (json, text)
- LOG_DEVELOPMENT - Development mode flag

**Security Configuration:**
- RATE_LIMIT_REQUESTS - Rate limit requests per window
- RATE_LIMIT_WINDOW - Rate limit time window
- MAX_REQUEST_SIZE - Maximum request size
- ENABLE_CORS - CORS enablement flag
- ALLOWED_ORIGINS - Allowed CORS origins

### Configuration Management
The application uses a structured configuration system with the following categories:

**Database Configuration:**
- DBHost, DBPort, DBUser, DBPassword, DBName, DBSSLMode

**Authentication Configuration:**
- JWTSecret for token signing

**Server Configuration:**
- ServerPort and GinMode for application settings

**AWS S3 Configuration:**
- AWSRegion, AWSS3Bucket, AWSAccessKeyID, AWSSecretKey

**Security Configuration:**
- RateLimitRequests, RateLimitWindow, AllowedOrigins, MaxRequestSize, EnableCORS

### Docker Support
The application can be containerized using a multi-stage Docker build:
- **Builder Stage**: Uses Go 1.24.4 Alpine image for compilation
- **Runtime Stage**: Uses Alpine Linux for minimal runtime environment
- **Dependencies**: Downloads Go modules and builds the application
- **Security**: Includes CA certificates for secure connections

---

## Security Features

### Authentication Security
- **JWT Tokens**: Secure token-based authentication
- **Password Hashing**: Bcrypt with configurable salt rounds
- **Token Expiration**: Configurable token lifetime
- **Role-Based Access**: Granular permission system

### Input Validation
- **Request Validation**: Comprehensive input validation
- **SQL Injection Protection**: Parameterized queries
- **XSS Protection**: Input sanitization
- **File Upload Validation**: Type and size validation

### Security Headers
The application implements comprehensive security headers:
- **X-Content-Type-Options**: Prevents MIME type sniffing
- **X-Frame-Options**: Prevents clickjacking attacks
- **X-XSS-Protection**: Enables XSS filtering
- **Strict-Transport-Security**: Enforces HTTPS connections
- **Content-Security-Policy**: Controls resource loading

### Rate Limiting
- **Request Limiting**: Configurable rate limits
- **IP-Based Limiting**: Per-IP request limits
- **Endpoint-Specific**: Different limits for different endpoints

### CORS Configuration
The application implements configurable CORS settings:
- **AllowOrigins**: Configured allowed origins for cross-origin requests
- **AllowMethods**: Supported HTTP methods (GET, POST, PUT, DELETE, OPTIONS)
- **AllowHeaders**: Allowed request headers (Origin, Content-Type, Authorization)
- **ExposeHeaders**: Headers exposed to the client
- **AllowCredentials**: Support for credentials in cross-origin requests
- **MaxAge**: Cache duration for preflight requests

---

## Performance & Monitoring

### Logging System
- **Structured Logging**: JSON-formatted logs with Zap
- **Log Levels**: Debug, Info, Warn, Error, Fatal
- **Request Tracking**: Unique request IDs for tracing
- **Performance Metrics**: Response time tracking

### Performance Monitoring
The application implements comprehensive performance monitoring:
- **Request Timing**: Tracks request start and end times
- **Duration Measurement**: Calculates request processing duration
- **Method Tracking**: Records HTTP method for each request
- **Path Monitoring**: Tracks request paths for analysis
- **Status Code Logging**: Records response status codes
- **Structured Logging**: Uses Zap for efficient JSON logging

### Health Checks
The application provides comprehensive health check endpoints:
- **GET /health** - Basic application health status
- **GET /health/db** - Database connectivity and health status

### Background Job Processing
- **Redis Queue**: Redis-based job queue
- **Job Types**: Email sending, file processing, analytics
- **Retry Logic**: Configurable retry attempts
- **Dead Letter Queue**: Failed job handling

### Database Optimization
- **Connection Pooling**: GORM connection pool configuration
- **Query Optimization**: Efficient database queries
- **Indexing**: Strategic database indexes
- **Migration System**: Versioned database migrations

---

## Data Flow Architecture

### Request Processing Flow
The application follows a layered request processing pattern:

1. **HTTP Layer**: Gin router receives incoming requests
2. **Middleware Stack**: Security, logging, and authentication middleware process requests
3. **Handler Layer**: Route-specific handlers validate and parse requests
4. **Service Layer**: Business logic processing and validation
5. **Repository Layer**: Data access and database operations
6. **Response Layer**: Structured JSON responses with consistent formatting

### Authentication Flow
The authentication system implements a stateless JWT-based approach:

1. **Login Process**: User credentials validated against database
2. **Token Generation**: JWT token created with user claims and roles
3. **Token Storage**: Client stores token for subsequent requests
4. **Request Authentication**: Middleware validates token on each request
5. **Authorization**: Role-based permissions checked for resource access
6. **Context Propagation**: User information passed through request context

### File Upload Architecture
The file storage system implements a robust upload and serving mechanism:

1. **Upload Validation**: File type, size, and security validation
2. **S3 Integration**: Files uploaded to AWS S3 with unique keys
3. **Database Storage**: File metadata and S3 keys stored in database
4. **Access Control**: File access controlled through authentication
5. **CDN Integration**: Files served through optimized endpoints
6. **Cleanup Process**: Orphaned files cleaned up through background jobs

---

## Scalability Architecture

### Horizontal Scaling Design
The application is designed for horizontal scaling with several key architectural decisions:

**Stateless Design**: All application state is externalized to databases and caches, allowing multiple instances to run independently without shared state.

**Database Connection Pooling**: GORM connection pools are configured to handle multiple concurrent connections efficiently, supporting high-traffic scenarios.

**Redis Integration**: Background job processing and caching are handled through Redis, providing distributed processing capabilities.

**Load Balancer Ready**: The application can be deployed behind load balancers with session affinity not required due to stateless design.

### Performance Optimization Strategies

**Database Optimization**:
- Strategic indexing on frequently queried columns
- Query optimization through GORM's query builder
- Connection pooling for efficient database resource utilization
- Read replicas support for read-heavy operations

**Caching Strategy**:
- Redis-based caching for frequently accessed data
- Application-level caching for static configuration
- File serving optimization through CDN integration
- Database query result caching for expensive operations

**Background Processing**:
- Asynchronous job processing for non-critical operations
- Email sending, file processing, and analytics handled in background
- Retry mechanisms for failed operations
- Dead letter queues for problematic jobs

---

## Security Architecture

### Multi-Layer Security Model
The application implements defense in depth with multiple security layers:

**Network Security**:
- HTTPS enforcement through security headers
- CORS configuration for cross-origin request control
- Rate limiting to prevent abuse and DDoS attacks
- Request size limits to prevent resource exhaustion

**Application Security**:
- Input validation and sanitization at multiple layers
- SQL injection prevention through parameterized queries
- XSS protection through output encoding
- CSRF protection through token validation

**Authentication Security**:
- JWT tokens with configurable expiration
- Password hashing using bcrypt with salt
- Role-based access control with granular permissions
- Session management through secure token storage

**Data Security**:
- Sensitive data encryption at rest
- Secure file storage with access controls
- Database connection encryption
- Audit logging for security events

### Authorization Model
The authorization system implements a flexible permission-based model:

**Resource-Based Permissions**: Each resource type has specific permission sets (create, read, update, delete).

**Role-Based Access**: Users are assigned roles that determine their access levels across different resources.

**Context-Aware Authorization**: Permissions can be context-dependent, such as users only being able to modify their own profiles.

**Hierarchical Permissions**: Admin roles have elevated permissions that can override standard access controls.

---

## Integration Architecture

### External Service Integration
The application integrates with several external services through well-defined interfaces:

**AWS S3 Integration**:
- File storage abstraction layer for easy provider switching
- Configurable endpoints for different S3-compatible services
- Automatic retry mechanisms for failed uploads
- Metadata management for file organization

**Email Service Integration**:
- SMTP configuration for transactional emails
- Template-based email generation
- Queue-based email processing for reliability
- Delivery status tracking and retry logic

**Redis Integration**:
- Job queue management for background processing
- Caching layer for performance optimization
- Session storage for distributed applications
- Pub/Sub messaging for real-time features

### API Integration Patterns
The application follows consistent patterns for external API integration:

**Circuit Breaker Pattern**: Prevents cascading failures when external services are unavailable.

**Retry Logic**: Automatic retry with exponential backoff for transient failures.

**Timeout Management**: Configurable timeouts to prevent hanging requests.

**Error Handling**: Graceful degradation when external services fail.

---

## Monitoring and Observability

### Logging Architecture
The application implements comprehensive logging for observability:

**Structured Logging**: All logs are in JSON format for easy parsing and analysis.

**Log Levels**: Debug, Info, Warn, Error, and Fatal levels for different severity levels.

**Request Tracing**: Unique request IDs for tracking requests across service boundaries.

**Performance Logging**: Response time and resource usage tracking.

**Security Logging**: Authentication and authorization events for security monitoring.

### Health Monitoring
The application provides multiple health check endpoints:

**Application Health**: Basic application status and version information.

**Database Health**: Database connectivity and query performance monitoring.

**External Service Health**: Status of integrated services like S3 and Redis.

**Resource Health**: Memory usage, CPU utilization, and connection pool status.

### Metrics Collection
The application collects various metrics for performance monitoring:

**Request Metrics**: Response times, request counts, and error rates.

**Database Metrics**: Query performance, connection pool usage, and transaction rates.

**Business Metrics**: User registrations, job applications, and system usage patterns.

**Infrastructure Metrics**: Server resource utilization and external service performance.

---

## Deployment Architecture

### Environment Management
The application supports multiple deployment environments:

**Development Environment**: Local development with hot reload and debug logging.

**Staging Environment**: Production-like environment for testing and validation.

**Production Environment**: Optimized configuration with security hardening.

**Configuration Management**: Environment-specific configuration through environment variables.

### Containerization Strategy
The application is designed for containerized deployment:

**Docker Support**: Multi-stage Docker builds for optimized image sizes.

**Kubernetes Ready**: Stateless design compatible with Kubernetes deployment.

**Health Checks**: Container health checks for orchestration platforms.

**Resource Management**: Configurable resource limits and requests.

### CI/CD Integration
The application supports continuous integration and deployment:

**Automated Testing**: Unit tests, integration tests, and end-to-end tests.

**Code Quality**: Linting, formatting, and security scanning.

**Automated Deployment**: Pipeline-based deployment to different environments.

**Rollback Capability**: Quick rollback mechanisms for failed deployments.

---

## Data Architecture

### Data Modeling Principles
The database design follows several key principles:

**Normalization**: Proper database normalization to reduce redundancy and maintain data integrity.

**Referential Integrity**: Foreign key constraints to maintain data relationships.

**Indexing Strategy**: Strategic indexing for optimal query performance.

**Data Types**: Appropriate data types for different kinds of information.

### Data Migration Strategy
The application implements a robust migration system:

**Versioned Migrations**: Sequential migration files with version numbers.

**Rollback Support**: Ability to rollback migrations if needed.

**Data Preservation**: Migrations designed to preserve existing data.

**Testing**: Migration testing in staging environments before production.

### Backup and Recovery
The application implements comprehensive backup strategies:

**Database Backups**: Regular automated database backups.

**File Backups**: S3 versioning and cross-region replication.

**Configuration Backups**: Version-controlled configuration management.

**Recovery Procedures**: Documented recovery procedures for different failure scenarios.

---

## Conclusion

The ASA Job Portal Backend represents a modern, scalable, and maintainable Go application with:

- **Clean Architecture**: Well-organized code structure with clear separation of concerns
- **Modern Technologies**: Latest Go version with contemporary libraries and frameworks
- **Comprehensive Features**: Full-featured job portal with advanced capabilities
- **Security-First Design**: Multiple layers of security and validation
- **Scalable Infrastructure**: Designed for horizontal scaling and high availability
- **Production Ready**: Comprehensive monitoring, logging, and deployment support
- **Maintainable Codebase**: Consistent patterns and comprehensive documentation

The architecture demonstrates best practices in modern web application development, providing a solid foundation for a production-grade job portal system that can scale to meet growing user demands while maintaining security, performance, and reliability standards.