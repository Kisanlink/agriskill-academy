# ASA Backend

A robust Go-based backend API for the ASA (Agricultural Student Association) job portal platform. This application provides comprehensive job posting, application management, and user profile services for agricultural students and employers.

## 🚀 Features

### Core Functionality
- **User Authentication & Authorization** - JWT-based authentication with role-based access control
- **Student Profiles** - Complete student profile management with resume and certificate uploads
- **Employer Profiles** - Employer profile management and company information
- **Job Postings** - Create, manage, and search job postings
- **Applications** - Job application system for students and application management for employers
- **Bookmarks** - Save and manage favorite job postings
- **File Storage** - Secure file upload system for resumes, certificates, and documents
- **Notifications** - Real-time notification system
- **Admin Panel** - Administrative tools for platform management

### Technical Features
- **RESTful API** - Clean, well-documented REST endpoints
- **Database Integration** - PostgreSQL with GORM ORM
- **File Upload** - Multi-format file upload with validation
- **CORS Support** - Cross-origin resource sharing enabled (configurable via CORS_ALLOW_ORIGINS in .env)
- **Hot Reloading** - Development mode with Air for hot reloading
- **Comprehensive Testing** - Built-in testing framework

## 📋 Prerequisites

- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **PostgreSQL 12+** - [Download PostgreSQL](https://www.postgresql.org/download/)
- **Git** - [Download Git](https://git-scm.com/downloads)

## 🛠️ Installation

### 1. Clone the Repository
```bash
git clone <repository-url>
cd agrijobs
```

### 2. Environment Setup
Create a `.env` file in the root directory:
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
POSTGRES_USER=your_username
POSTGRES_PASS=your_password
DB_NAME=asa_db
DB_SSLMODE=disable

# AAA Service Configuration (Optional - for external auth)
AAA_SERVICE_URL=http://localhost:8080
SECRET_KEY=your_aaa_secret

# Server Configuration
PORT=3333

# CORS Configuration
CORS_ALLOW_ORIGINS=*
```

### 3. Install Dependencies
```bash
make setup
```

### 4. Database Setup
```bash
# Run all migrations
make migrate

# Or reset database and run migrations
make migrate-reset
```

## 🚀 Quick Start

### Development Mode
```bash
# Start with hot reloading (recommended for development)
make air

# Or start normally
make run
```

### Production Mode
```bash
# Build the application
make build

# Run the built application
make run-build
```

## 📚 Available Commands

### Development
```bash
make run          # Run the application
make air          # Run with hot reloading
make dev          # Alias for air
```

### Database Management
```bash
make migrate      # Run all migrations
make migrate-reset # Reset database and run migrations
make migrate-schema # Apply only schema migration
make migrate-messages # Apply only messages fix
make migrate-profiles # Apply only profiles rename
```

### Building
```bash
make build        # Build for current platform
make build-linux  # Build for Linux
make build-windows # Build for Windows
make build-mac    # Build for macOS
```

### Testing & Quality
```bash
make test         # Run tests
make test-verbose # Run tests with verbose output
make test-coverage # Run tests with coverage
make fmt          # Format code
make vet          # Run go vet
make tidy         # Tidy Go modules
```

### Maintenance
```bash
make clean        # Clean build artifacts
make clean-uploads # Clean uploaded files
make deps         # Install dependencies
make setup        # Setup development environment
```

## 🏗️ Project Structure

```
agrijobs/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── config/
│   └── config.go            # Configuration management
├── internal/
│   ├── admin/               # Admin functionality
│   ├── application/         # Job applications
│   ├── auth/               # Authentication & authorization
│   ├── bookmark/           # Bookmark management
│   ├── employerapplication/ # Employer application views
│   ├── employerprofile/    # Employer profiles
│   ├── jobpost/           # Job posting management
│   ├── middleware/         # HTTP middleware
│   ├── notification/       # Notification system
│   ├── storage/           # File upload & storage
│   ├── studentprofile/    # Student profiles
│   └── worker/            # Background job processing
├── migrations/             # Database migrations
├── pkg/
│   ├── authz/             # Authorization utilities
│   ├── jwtutil/           # JWT utilities
│   └── mailutil/          # Email utilities
├── scripts/               # Database scripts
├── uploads/               # File upload storage
│   ├── certificates/      # Certificate files
│   └── resumes/          # Resume files
├── Makefile               # Build and management commands
├── go.mod                 # Go module definition
└── README.md             # This file
```

## 🔐 Authentication & Authorization

The application uses JWT-based authentication with role-based authorization:

### Roles
- **student** - Can apply for jobs, manage profile, upload documents
- **employer** - Can post jobs, manage applications, manage company profile

### Authorization Flow
1. User signs up/logs in via `/api/auth/signup` or `/api/auth/login`
2. JWT token is generated with user role and ID
3. Token is included in subsequent requests via `Authorization: Bearer <token>`
4. Middleware validates token and extracts user information
5. Role-based permissions are checked for each protected endpoint

## 🗄️ Database Schema

### Core Tables
- **users** - User accounts and authentication
- **student_profiles** - Student profile information
- **employer_profiles** - Employer/company profiles
- **job_posts** - Job postings
- **applications** - Job applications
- **certificates** - Student certificates
- **bookmarks** - User bookmarks
- **notifications** - User notifications

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage
make test-coverage
```

## 🚀 Deployment

### Local Development
```bash
make setup    # Setup development environment
make migrate  # Setup database
make air      # Start with hot reloading
```

### Production
```bash
make build-linux  # Build for Linux
./build/asa-linux # Run the binary
```

## 🔧 Configuration

### Environment Variables
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `POSTGRES_USER` - Database username
- `POSTGRES_PASS` - Database password
- `DB_NAME` - Database name
- `DB_SSLMODE` - Database SSL mode
- `AAA_SERVICE_URL` - AAA Service URL
- `SECRET_KEY` - JWT secret key

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Contact the development team
- Check the API documentation

## 🔄 Changelog

### Version 1.0.0
- Initial release
- Complete job portal functionality
- User authentication and authorization
- File upload system
- Student and employer profiles
- Job posting and application system
---