# AgriJobs API Documentation

## Base URL
```
http://localhost:3000/api
```

## Authentication
Most endpoints require authentication via JWT token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

---

## 🔐 Authentication Endpoints

### 1. **POST /api/auth/signup**
**Description:** Register a new user (student or employer)

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "confirmPassword": "password123",
  "role": "employer",
  "companyName": "AgriTech Solutions",
  "gstinNumber": "22AAAAA0000A1Z5",
  "companyAddress": "123 Farm Road",
  "city": "Hyderabad",
  "state": "Telangana",
  "pincode": "500001",
  "industryType": "AgriTech / Smart Farming",
  "companySize": "51-200 employees",
  "website": "https://agritech.com"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Signup successful",
  "user": {
    "id": "uuid-string",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "employer"
  },
  "token": "jwt_token_string"
}
```

---

### 2. **POST /api/auth/login**
**Description:** Login with email and password

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "password123",
  "role": "employer"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "user": {
    "id": "uuid-string",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "employer",
    "createdAt": "2024-01-15T10:00:00Z",
    "updatedAt": "2024-01-15T10:00:00Z"
  },
  "token": "jwt_token_string"
}
```

---

### 3. **PUT /api/auth/profile**
**Description:** Update basic user information (name, email)

**Request Body:**
```json
{
  "name": "Updated Name",
  "email": "updated@email.com"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Profile updated",
  "user": {
    "id": "uuid-string",
    "name": "Updated Name",
    "email": "updated@email.com",
    "role": "employer",
    "createdAt": "2024-01-15T10:00:00Z",
    "updatedAt": "2024-01-15T10:30:00Z"
  }
}
```

---

### 4. **GET /api/auth/verify**
**Description:** Verify JWT token validity

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Token is valid"
}
```

---

### 5. **POST /api/auth/forgot-password**
**Description:** Send password reset link

**Request Body:**
```json
{
  "email": "john@example.com"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Reset link sent"
}
```

---

### 6. **POST /api/auth/reset-password**
**Description:** Reset password using token

**Request Body:**
```json
{
  "token": "reset_token_string",
  "newPassword": "newpassword123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Password reset successful"
}
```

---

## 👔 Employer Profile Endpoints

### 7. **GET /api/employers/me/profile**
**Description:** Get current employer's profile

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": "uuid-string",
    "companyName": "AgriTech Solutions",
    "logo": "https://example.com/logo.png",
    "websiteUrl": "https://agritech.com",
    "industry": "AgriTech / Smart Farming",
    "companySize": "51-200 employees",
    "companyDescription": "Leading provider of smart farming solutions",
    "recruiterName": "John Doe",
    "designation": "HR Manager",
    "officialEmail": "hr@agritech.com",
    "phoneNumber": "+1-555-1234",
    "linkedinProfile": "https://linkedin.com/in/johndoe",
    "jobCategories": ["AgriTech Development", "Field Operations"],
    "hiringLocations": ["California, USA"],
    "hiringTypes": ["Full-time", "Internship"],
    "createdAt": "2024-01-15T10:00:00Z",
    "updatedAt": "2024-01-15T10:00:00Z"
  }
}
```

---

### 8. **PUT /api/employers/me/profile**
**Description:** Update employer profile

**Request Body:**
```json
{
  "companyName": "Updated Company Name",
  "websiteUrl": "https://updated-agritech.com",
  "industry": "Crop Production",
  "companySize": "201-500 employees",
  "companyDescription": "Updated company description",
  "recruiterName": "Updated Name",
  "designation": "Senior HR Manager",
  "officialEmail": "hr@updated-agritech.com",
  "phoneNumber": "+1-555-5678",
  "linkedinProfile": "https://linkedin.com/in/updated-profile",
  "jobCategories": ["AgriTech Development", "Research & Development"],
  "hiringLocations": ["California, USA", "Texas, USA"],
  "hiringTypes": ["Full-time", "Part-time", "Internship"]
}
```

**Response:**
```json
{
  "success": true,
  "message": "Employer profile updated successfully",
  "data": {
    // Updated profile data
  }
}
```

---

## 👤 Student Profile Endpoints

### 9. **GET /api/students/me/profile**
**Description:** Get current student's profile

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": "uuid-string",
    "name": "Student Name",
    "email": "student@email.com",
    "location": "Hyderabad, India",
    "profilePhoto": "https://example.com/photo.png",
    "resume": "https://example.com/resume.pdf",
    "certificates": [
      {
        "id": "cert1",
        "name": "Agri Certification",
        "file": "https://example.com/cert1.pdf",
        "issueDate": "2023-01-01"
      }
    ],
    "skills": ["Crop Management", "Data Analysis"],
    "createdAt": "2024-01-15T10:00:00Z",
    "updatedAt": "2024-01-15T10:00:00Z"
  }
}
```

---

### 10. **PUT /api/students/me/profile**
**Description:** Update student profile

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "string",
  "email": "string",
  "location": "string",
  "skills": ["string"],
  "certificates": [
    {
      "id": "string",
      "name": "string",
      "file": "string (valid file URL)",
      "issueDate": "string (ISO date)"
    }
  ],
  "profilePhoto": "string (valid file URL)",
  "resume": "string (valid file URL)"
}
```

**Example Request:**
```json
{
  "name": "studKA Updated",
  "email": "studka@kisanlink.com",
  "location": "Hyderabad, Telangana, India",
  "skills": ["IoT Sensors", "JavaScript", "React"],
  "certificates": [
    {
      "id": "1751490769347",
      "name": "Web Development Certificate",
      "file": "https://example.com/certificate.pdf",
      "issueDate": "2025-01-15T00:00:00.000Z"
    }
  ],
  "profilePhoto": "https://example.com/profile-photo.jpg",
  "resume": "https://example.com/resume.pdf"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "data": {
    // Updated profile data
  }
}
```

---

### 11. **GET /api/students/:studentId/profile**
**Description:** Get specific student's profile (for employers)

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Profile fetched",
  "data": {
    "id": "uuid-string",
    "name": "Student Name",
    "email": "student@email.com",
    "location": "Hyderabad, India",
    "profilePhoto": "https://example.com/photo.png",
    "resume": "https://example.com/resume.pdf",
    "certificates": [
      {
        "id": "cert1",
        "name": "Agri Certification",
        "file": "https://example.com/cert1.pdf",
        "issueDate": "2023-01-01"
      }
    ],
    "skills": ["Crop Management", "Data Analysis"],
    "createdAt": "2024-01-15T10:00:00Z",
    "updatedAt": "2024-01-15T10:00:00Z"
  }
}
```

---

### 12. **PUT /api/students/:studentId/profile**
**Description:** Update specific student's profile (admin only)

**Request Body:**
```json
{
  "name": "Updated Student Name",
  "email": "updated@email.com",
  "location": "Mumbai, India",
  "profilePhoto": "https://example.com/new-photo.png",
  "resume": "https://example.com/new-resume.pdf",
  "skills": ["Crop Management", "Data Analysis", "Sustainable Agriculture"]
}
```

**Response:**
```json
{
  "success": true,
  "message": "Profile updated",
  "data": {
    // Updated profile data
  }
}
```

---

### 13. **POST /api/students/me/certificates**
**Description:** Add certificate to current student's profile

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "name": "AgriTech Certification",
  "file": "https://example.com/certificate.pdf",
  "issueDate": "2024-01-01"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Certificate added successfully",
  "data": {
    "id": "cert-uuid",
    "name": "AgriTech Certification",
    "file": "https://example.com/certificate.pdf",
    "issueDate": "2024-01-01"
  }
}
```

---

### 14. **POST /api/students/:studentId/certificates**
**Description:** Add certificate to specific student's profile (admin only)

**Request Body:**
```json
{
  "name": "AgriTech Certification",
  "file": "https://example.com/certificate.pdf",
  "issueDate": "2024-01-01"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Certificate added",
  "data": {
    "id": "cert-uuid",
    "name": "AgriTech Certification",
    "file": "https://example.com/certificate.pdf",
    "issueDate": "2024-01-01"
  }
}
```

---

## 💼 Job Management Endpoints

### Basic Job Operations

#### Get All Published Jobs (For Students)
- **GET** `/api/jobs`
- **Description**: Get all published jobs with optional filtering and pagination
- **Auth**: Not required (Public endpoint for students)
- **Query Parameters**:
  - `page` (optional): Page number (default: 1)
  - `limit` (optional): Number of jobs per page (default: 20, max: 100)
  - `location` (optional): Filter by location
  - `jobType` (optional): Filter by job type (full-time, part-time, contract, internship)
  - `experience` (optional): Filter by experience level (entry, mid, senior)
  - `isRemote` (optional): Filter by remote work (true/false)
- **Response**:
```json
{
  "success": true,
  "jobs": [
    {
      "id": "job_123",
      "title": "Software Engineer",
      "company": "GreenTech Innovations",
      "location": "Remote",
      "jobType": "full-time",
      "experience": "mid",
      "description": "Work on agricultural software solutions.",
      "requirements": [
        "3+ years experience",
        "React, Node.js"
      ],
      "skills": ["React", "Node.js"],
      "postedAt": "2024-06-01T10:00:00.000Z",
      "applicationDeadline": "2024-07-01T23:59:59.000Z",
      "salary": {
        "min": 50000,
        "max": 80000,
        "currency": "USD"
      },
      "recruiter": {
        "name": "Jane Doe",
        "email": "jane@greentech.com",
        "company": "GreenTech Innovations",
        "avatar": null
      },
      "benefits": ["Health Insurance", "401k"],
      "isRemote": true,
      "applicationsCount": 12,
      "status": "active"
    }
  ]
}
```

#### Create Job Post
- **POST** `/api/jobs`
- **Description**: Create a new job post
- **Auth**: Required (Employer)
- **Request Body**:
```json
{
  "title": "Software Engineer",
  "roleOverview": "We are looking for a skilled software engineer...",
  "requirements": "Bachelor's degree in Computer Science...",
  "location": "New York, NY",
  "requiredSkills": ["JavaScript", "React", "Node.js"],
  "applicationDeadline": "2024-12-31T23:59:59Z",
  "jobType": "full-time",
  "experience": "mid",
  "salary": {
    "min": 80000,
    "max": 120000,
    "currency": "USD"
  },
  "benefits": ["Health Insurance", "401k", "Remote Work"],
  "isRemote": true
}
```
- **Response**:
```json
{
  "success": true,
  "message": "Job created successfully",
  "jobPost": {
    "id": "uuid",
    "title": "Software Engineer",
    "roleOverview": "We are looking for a skilled software engineer...",
    "requirements": "Bachelor's degree in Computer Science...",
    "location": "New York, NY",
    "requiredSkills": ["JavaScript", "React", "Node.js"],
    "employerId": "uuid",
    "employerName": "Tech Corp",
    "employerEmail": "hr@techcorp.com",
    "status": "draft",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z",
    "applicationDeadline": "2024-12-31T23:59:59Z",
    "jobType": "full-time",
    "experience": "mid",
    "salary": {
      "min": 80000,
      "max": 120000,
      "currency": "USD"
    },
    "benefits": ["Health Insurance", "401k", "Remote Work"],
    "isRemote": true,
    "applicationsCount": 0
  }
}
```

#### Update Job Post
- **PUT** `/api/jobs/:id`
- **Description**: Update an existing job post
- **Auth**: Required (Employer - owner only)
- **Request Body**: Same as Create Job Post (all fields optional)
- **Response**: Same as Create Job Post

#### Delete Job Post
- **DELETE** `/api/jobs/:id`
- **Description**: Delete a job post
- **Auth**: Required (Employer - owner only)
- **Response**:
```json
{
  "success": true,
  "message": "Job deleted successfully"
}
```

#### Get Job Post by ID
- **GET** `/api/jobs/:id`
- **Description**: Get a specific job post
- **Auth**: Not required
- **Response**: Same as Create Job Post

#### Get My Job Posts
- **GET** `/api/jobs/my-posts`
- **Description**: Get all job posts by the authenticated employer
- **Auth**: Required (Employer)
- **Response**:
```json
{
  "success": true,
  "message": "Jobs retrieved successfully",
  "jobPosts": [
    {
      "id": "uuid",
      "title": "Software Engineer",
      // ... other job fields
    }
  ]
}
```

### Enhanced Job Search & Discovery

#### Advanced Job Search
- **POST** `/api/jobs/advanced-search`
- **Description**: Advanced job search with multiple filters and sorting options
- **Auth**: Not required
- **Request Body**:
```json
{
  "keywords": "software engineer react",
  "location": "New York",
  "jobType": ["full-time", "part-time"],
  "experience": ["mid", "senior"],
  "skills": ["JavaScript", "React"],
  "industry": ["Technology", "Finance"],
  "companySize": ["50-200", "200-1000"],
  "salaryRange": {
    "min": 60000,
    "max": 150000,
    "currency": "USD"
  },
  "benefits": ["Health Insurance", "Remote Work"],
  "isRemote": true,
  "isHybrid": false,
  "isOnsite": false,
  "postedWithin": "30d",
  "urgent": true,
  "sortBy": "relevance",
  "sortOrder": "desc",
  "page": 1,
  "limit": 20
}
```
- **Response**:
```json
{
  "success": true,
  "message": "Search completed successfully",
  "jobs": [
    {
      "id": "uuid",
      "title": "Senior Software Engineer",
      // ... other job fields
    }
  ],
  "filters": {
    "availableLocations": ["New York", "San Francisco", "Remote"],
    "availableJobTypes": ["full-time", "part-time", "contract"],
    "availableExperience": ["entry", "mid", "senior"],
    "availableSkills": ["JavaScript", "React", "Node.js", "Python"],
    "availableIndustries": ["Technology", "Finance", "Healthcare"],
    "availableCompanySizes": ["1-50", "50-200", "200-1000", "1000+"],
    "salaryRanges": [
      {
        "min": 0,
        "max": 30000,
        "currency": "USD",
        "label": "$0-$30k"
      },
      {
        "min": 30000,
        "max": 50000,
        "currency": "USD",
        "label": "$30k-$50k"
      }
    ],
    "availableBenefits": ["Health Insurance", "401k", "Remote Work", "Flexible Hours"]
  },
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "totalPages": 8,
    "hasNext": true,
    "hasPrev": false
  }
}
```

#### Get Search Filters
- **GET** `/api/jobs/search-filters`
- **Description**: Get available search filters for the UI
- **Auth**: Not required
- **Response**: Same as the `filters` object in Advanced Search response

#### Get Featured Jobs
- **GET** `/api/jobs/featured?limit=10`
- **Description**: Get featured/popular jobs
- **Auth**: Not required
- **Response**:
```json
{
  "success": true,
  "message": "Featured jobs retrieved successfully",
  "jobPosts": [
    {
      "id": "uuid",
      "title": "Senior Software Engineer",
      // ... other job fields
    }
  ]
}
```

#### Get Recent Jobs
- **GET** `/api/jobs/recent?limit=20`
- **Description**: Get recently posted jobs
- **Auth**: Not required
- **Response**: Same as Featured Jobs

#### Get Trending Jobs
- **GET** `/api/jobs/trending?limit=10`
- **Description**: Get trending jobs (high applications, recent posts)
- **Auth**: Not required
- **Response**: Same as Featured Jobs

#### Get Similar Jobs
- **GET** `/api/jobs/:id/similar?maxResults=5`
- **Description**: Get jobs similar to a specific job
- **Auth**: Not required
- **Response**: Same as Featured Jobs

#### Get Job Recommendations
- **POST** `/api/jobs/recommendations`
- **Description**: Get personalized job recommendations
- **Auth**: Required (Student)
- **Request Body**:
```json
{
  "userId": "uuid",
  "userSkills": ["JavaScript", "React", "Node.js"],
  "userLocation": "New York",
  "userExperience": "mid",
  "preferredJobTypes": ["full-time", "remote"],
  "maxResults": 10
}
```
- **Response**:
```json
{
  "success": true,
  "message": "Recommendations generated successfully",
  "jobs": [
    {
      "id": "uuid",
      "title": "Frontend Developer",
      // ... other job fields
    }
  ],
  "reason": "Based on your skills: JavaScript and others"
}
```

### Job Alerts

#### Create Job Alert
- **POST** `/api/jobs/alerts`
- **Description**: Create a job alert for notifications
- **Auth**: Required (Student)
- **Request Body**:
```json
{
  "keywords": ["software engineer", "developer"],
  "location": "New York",
  "jobType": ["full-time", "remote"],
  "experience": ["mid", "senior"],
  "skills": ["JavaScript", "React"],
  "salaryRange": {
    "min": 60000,
    "max": 150000,
    "currency": "USD"
  },
  "isRemote": true,
  "frequency": "weekly",
  "isActive": true
}
```
- **Response**:
```json
{
  "success": true,
  "message": "Job alert created successfully",
  "alert": {
    "id": "uuid",
    "userId": "uuid",
    "keywords": ["software engineer", "developer"],
    "location": "New York",
    "jobType": ["full-time", "remote"],
    "experience": ["mid", "senior"],
    "skills": ["JavaScript", "React"],
    "salaryRange": {
      "min": 60000,
      "max": 150000,
      "currency": "USD"
    },
    "isRemote": true,
    "frequency": "weekly",
    "isActive": true,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

#### Get My Job Alerts
- **GET** `/api/jobs/alerts`
- **Description**: Get all job alerts for the authenticated user
- **Auth**: Required (Student)
- **Response**:
```json
{
  "success": true,
  "message": "Job alerts retrieved successfully",
  "alerts": [
    {
      "id": "uuid",
      "userId": "uuid",
      // ... other alert fields
    }
  ]
}
```

#### Get Job Alert by ID
- **GET** `/api/jobs/alerts/:id`
- **Description**: Get a specific job alert
- **Auth**: Required (Student - owner only)
- **Response**: Same as Create Job Alert

#### Update Job Alert
- **PUT** `/api/jobs/alerts/:id`
- **Description**: Update a job alert
- **Auth**: Required (Student - owner only)
- **Request Body**: Same as Create Job Alert (all fields optional)
- **Response**: Same as Create Job Alert

#### Delete Job Alert
- **DELETE** `/api/jobs/alerts/:id`
- **Description**: Delete a job alert
- **Auth**: Required (Student - owner only)
- **Response**:
```json
{
  "success": true,
  "message": "Job alert deleted successfully"
}
```

### Basic Job Search (Legacy)
- **POST** `/api/jobs/search`
- **Description**: Basic job search (legacy endpoint)
- **Auth**: Not required
- **Request Body**:
```json
{
  "location": "New York",
  "jobType": ["full-time"],
  "experience": ["mid"],
  "salaryRange": {
    "min": 60000,
    "max": 150000
  },
  "isRemote": true,
  "skills": ["JavaScript"],
  "postedWithin": "30d",
  "page": 1,
  "limit": 20
}
```
- **Response**: Same as Advanced Search (without filters and pagination)

---

## 📝 Application Endpoints

### 22. **POST /api/jobs/:jobId/apply**
**Description:** Apply for a job

**Request Body:**
```json
{
  "coverLetter": "I am very interested in this position...",
  "resumeFile": "file_upload_data"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Application submitted successfully"
}
```

---

### 23. **GET /api/applications/my**
**Description:** Get user's job applications

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Applications retrieved successfully",
  "applications": [
    {
      "id": "uuid-string",
      "jobId": "uuid-string",
      "studentId": "uuid-string",
      "appliedAt": "2024-01-15T10:00:00Z",
      "status": "pending",
      "coverLetter": "I am very interested...",
      "jobTitle": "Agricultural Engineer",
      "company": "AgriTech Solutions",
      "location": "California, USA",
      "jobType": "full-time",
      "experience": "senior"
    }
  ]
}
```

---

### 24. **DELETE /api/applications/:applicationId**
**Description:** Remove job application

**Response:**
```json
{
  "success": true,
  "message": "Application removed successfully"
}
```

---

### 25. **PUT /api/applications/:applicationId/status**
**Description:** Update application status

**Request Body:**
```json
{
  "status": "shortlisted"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Application status updated"
}
```

---

## 🏢 Employer Application Management

### 26. **GET /api/employer/jobs/:jobId/applications**
**Description:** Get applications for a specific job

**Response:**
```json
{
  "success": true,
  "message": "Applications retrieved successfully",
  "applications": [
    {
      "id": "uuid-string",
      "jobId": "uuid-string",
      "studentId": "uuid-string",
      "appliedAt": "2024-01-15T10:00:00Z",
      "status": "pending",
      "coverLetter": "I am very interested...",
      "jobTitle": "Agricultural Engineer",
      "company": "AgriTech Solutions",
      "location": "California, USA",
      "jobType": "full-time",
      "experience": "senior",
      "applicant": {
        "id": "uuid-string",
        "name": "Sarah Johnson",
        "email": "sarah@email.com",
        "phone": "+1-555-123-4567",
        "location": "California, USA",
        "skills": ["Agricultural Engineering", "Crop Management"],
        "experience": "senior",
        "education": "MS in Agricultural Engineering",
        "avatar": "https://example.com/avatar.jpg",
        "resumeUrl": "https://example.com/resume.pdf",
        "summary": "Experienced agricultural engineer..."
      }
    }
  ]
}
```

---

### 27. **PUT /api/employer/applications/:applicationId/status**
**Description:** Update application status (employer view)

**Request Body:**
```json
{
  "status": "shortlisted"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Applicant shortlisted successfully"
}
```

---

### 28. **GET /api/employer/applicants/:studentId/profile**
**Description:** Get applicant profile

**Response:**
```json
{
  "success": true,
  "message": "Applicant profile retrieved successfully",
  "profile": {
    "id": "uuid-string",
    "name": "Sarah Johnson",
    "email": "sarah@email.com",
    "phone": "+1-555-123-4567",
    "location": "California, USA",
    "skills": ["Agricultural Engineering", "Crop Management"],
    "experience": "senior",
    "education": "MS in Agricultural Engineering",
    "avatar": "https://example.com/avatar.jpg",
    "resumeUrl": "https://example.com/resume.pdf",
    "portfolio": "https://sarah-portfolio.com",
    "linkedIn": "https://linkedin.com/in/sarah-johnson",
    "summary": "Experienced agricultural engineer..."
  }
}
```

---

### 29. **POST /api/employer/applications/:applicationId/message**
**Description:** Send message to applicant

**Request Body:**
```json
{
  "message": "Thank you for your application. We would like to schedule an interview."
}
```

**Response:**
```json
{
  "success": true,
  "message": "Message sent to applicant successfully"
}
```

---

### 30. **GET /api/employer/applications/:applicationId/messages**
**Description:** Get messages for an application

**Response:**
```json
{
  "success": true,
  "message": "Messages retrieved successfully",
  "messages": [
    {
      "id": "uuid-string",
      "applicationId": "uuid-string",
      "senderId": "uuid-string",
      "senderName": "John Doe",
      "message": "Thank you for your application...",
      "sentAt": "2024-01-15T10:00:00Z"
    }
  ]
}
```

---

### 31. **GET /api/student/applications**
**Description:** Get applications by student (student view)

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Applications retrieved successfully",
  "applications": [
    {
      "id": "uuid-string",
      "jobId": "uuid-string",
      "studentId": "uuid-string",
      "appliedAt": "2024-01-15T10:00:00Z",
      "status": "pending",
      "coverLetter": "I am very interested...",
      "jobTitle": "Agricultural Engineer",
      "company": "AgriTech Solutions",
      "location": "California, USA",
      "jobType": "full-time",
      "experience": "senior"
    }
  ]
}
```

---

## 🔖 Bookmark Endpoints

### 32. **POST /api/bookmarks/:jobId**
**Description:** Save job to bookmarks

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Job saved successfully"
}
```

---

### 33. **DELETE /api/bookmarks/:jobId**
**Description:** Remove job from bookmarks

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Job removed from bookmarks"
}
```

---

## 📁 File Upload Endpoints

### 34. **POST /api/upload/:folder**
**Description:** Upload any file type (general upload)

**Request Body:** Multipart form data
```
file: <file_data>
```

**Response:**
```json
{
  "success": true,
  "message": "File uploaded successfully",
  "filePath": "resumes/1234567890_document.pdf",
  "fileName": "1234567890_document.pdf",
  "fileSize": 1024000,
  "fileType": "document",
  "fileUrl": "http://localhost:3000/api/files/resumes/1234567890_document.pdf"
}
```

---

### 35. **POST /api/upload/image/:folder**
**Description:** Upload image files (profile photos, logos, etc.)

**Allowed Formats:** JPG, JPEG, PNG, GIF, WebP
**Max Size:** 5MB

**Request Body:** Multipart form data
```
file: <image_file>
```

**Response:**
```json
{
  "success": true,
  "message": "Image uploaded successfully",
  "filePath": "images/1234567890_profile.jpg",
  "fileName": "1234567890_profile.jpg",
  "fileSize": 2048000,
  "fileType": "image",
  "fileUrl": "http://localhost:3000/api/files/images/1234567890_profile.jpg"
}
```

---

### 36. **POST /api/upload/document/:folder**
**Description:** Upload document files (PDFs, Word docs, etc.)

**Allowed Formats:** PDF, DOC, DOCX, TXT, RTF
**Max Size:** 10MB

**Request Body:** Multipart form data
```
file: <document_file>
```

**Response:**
```json
{
  "success": true,
  "message": "Document uploaded successfully",
  "filePath": "documents/1234567890_report.pdf",
  "fileName": "1234567890_report.pdf",
  "fileSize": 5120000,
  "fileType": "document",
  "fileUrl": "http://localhost:3000/api/files/documents/1234567890_report.pdf"
}
```

---

### 37. **POST /api/upload/resume/:folder**
**Description:** Upload resume files specifically

**Allowed Formats:** PDF, DOC, DOCX
**Max Size:** 10MB

**Request Body:** Multipart form data
```
file: <resume_file>
```

**Response:**
```json
{
  "success": true,
  "message": "Resume uploaded successfully",
  "filePath": "resumes/1234567890_resume.pdf",
  "fileName": "1234567890_resume.pdf",
  "fileSize": 1024000,
  "fileType": "document",
  "fileUrl": "http://localhost:3000/api/files/resumes/1234567890_resume.pdf"
}
```

---

### 38. **GET /api/files/:folder**
**Description:** List all files in a folder

**Response:**
```json
{
  "success": true,
  "message": "Files retrieved successfully",
  "files": [
    {
      "name": "1234567890_document.pdf",
      "size": 1024000,
      "type": "document",
      "path": "resumes/1234567890_document.pdf",
      "url": "http://localhost:3000/api/files/resumes/1234567890_document.pdf",
      "uploaded": "2024-01-15T10:00:00Z"
    }
  ]
}
```

---

### 39. **GET /api/files/info/*filePath**
**Description:** Get information about a specific file

**Response:**
```json
{
  "success": true,
  "message": "File info retrieved successfully",
  "file": {
    "name": "1234567890_document.pdf",
    "size": 1024000,
    "type": "document",
    "path": "resumes/1234567890_document.pdf",
    "url": "http://localhost:3000/api/files/resumes/1234567890_document.pdf",
    "uploaded": "2024-01-15T10:00:00Z"
  }
}
```

---

### 40. **DELETE /api/files/*filePath**
**Description:** Delete a file

**Response:**
```json
{
  "success": true,
  "message": "File deleted successfully"
}
```

---

## 📧 Notification Endpoints

### 35. **POST /api/notify/email**
**Description:** Send email notification

**Request Body:**
```json
{
  "to": "user@example.com",
  "subject": "Job Application Update",
  "body": "Your application has been reviewed..."
}
```

**Response:**
```json
{
  "success": true,
  "message": "Email sent successfully"
}
```

---

## 🔄 Background Job Endpoints

### 36. **POST /api/worker/job**
**Description:** Enqueue background job

**Request Body:**
```json
{
  "type": "email_notification",
  "payload": {
    "to": "user@example.com",
    "subject": "Welcome to AgriJobs",
    "body": "Thank you for joining our platform..."
  }
}
```

**Response:**
```json
{
  "success": true,
  "message": "Job enqueued"
}
```

---

## 👨‍💼 Admin Endpoints

### 41. **GET /api/admin/analytics/jobs**
**Description:** Get job analytics for admin dashboard

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Job analytics retrieved successfully",
  "data": {
    "totalJobs": 150,
    "publishedJobs": 120,
    "draftJobs": 20,
    "closedJobs": 10,
    "totalApplications": 450,
    "avgApplications": 3.0,
    "mostPopularJob": "Agricultural Engineer",
    "jobsThisMonth": 25,
    "jobsLastMonth": 20,
    "growthRate": 25.0,
    "jobsByLocation": [
      {
        "location": "California, USA",
        "count": 45,
        "percentage": 30.0
      }
    ],
    "jobsByType": [
      {
        "jobType": "full-time",
        "count": 90,
        "percentage": 60.0
      }
    ]
  }
}
```

---

### 42. **GET /api/admin/analytics/users**
**Description:** Get user analytics for admin dashboard

**Response:**
```json
{
  "success": true,
  "message": "User analytics retrieved successfully",
  "data": {
    "totalUsers": 500,
    "totalStudents": 350,
    "totalEmployers": 150,
    "activeUsers": 400,
    "newUsersThisMonth": 50,
    "newUsersLastMonth": 40,
    "userGrowthRate": 25.0,
    "topLocations": [
      {
        "location": "California, USA",
        "count": 120,
        "percentage": 24.0
      }
    ]
  }
}
```

---

### 43. **GET /api/admin/analytics/companies**
**Description:** Get company analytics for admin dashboard

**Response:**
```json
{
  "success": true,
  "message": "Company analytics retrieved successfully",
  "data": {
    "totalCompanies": 150,
    "activeCompanies": 120,
    "verifiedCompanies": 100,
    "newCompaniesThisMonth": 15,
    "newCompaniesLastMonth": 12,
    "companyGrowthRate": 25.0,
    "companiesByIndustry": [
      {
        "industry": "AgriTech",
        "count": 45,
        "percentage": 30.0
      }
    ],
    "companiesByLocation": [
      {
        "location": "California, USA",
        "count": 35,
        "percentage": 23.3
      }
    ],
    "companiesBySize": [
      {
        "size": "51-200 employees",
        "count": 60,
        "percentage": 40.0
      }
    ]
  }
}
```

---

### 44. **GET /api/admin/users**
**Description:** Get list of users with pagination and filtering

**Query Parameters:**
```
page=1&limit=10&role=employer&search=john&status=active&sortBy=created_at&sortOrder=desc
```

**Response:**
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": {
    "users": [
      {
        "id": "uuid-string",
        "name": "John Doe",
        "email": "john@example.com",
        "role": "employer",
        "status": "active",
        "createdAt": "2024-01-15T10:00:00Z",
        "updatedAt": "2024-01-15T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 150,
      "totalPages": 15
    }
  }
}
```

---

### 45. **GET /api/admin/companies**
**Description:** Get list of companies with pagination and filtering

**Query Parameters:**
```
page=1&limit=10&industry=AgriTech&search=tech&sortBy=created_at&sortOrder=desc
```

**Response:**
```json
{
  "success": true,
  "message": "Companies retrieved successfully",
  "data": {
    "companies": [
      {
        "id": "uuid-string",
        "companyName": "AgriTech Solutions",
        "industry": "AgriTech",
        "location": "California, USA",
        "companySize": "51-200 employees",
        "status": "active",
        "jobsCount": 15,
        "applicationsCount": 45,
        "createdAt": "2024-01-15T10:00:00Z",
        "updatedAt": "2024-01-15T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 150,
      "totalPages": 15
    }
  }
}
```

---

### 46. **GET /api/admin/companies/:id**
**Description:** Get detailed company information

**Response:**
```json
{
  "success": true,
  "message": "Company retrieved successfully",
  "data": {
    "id": "uuid-string",
    "companyName": "AgriTech Solutions",
    "industry": "AgriTech",
    "companySize": "51-200 employees",
    "companyDescription": "Leading provider of smart farming solutions",
    "location": "California, USA",
    "websiteUrl": "https://agritech.com",
    "recruiterName": "John Doe",
    "officialEmail": "hr@agritech.com",
    "phoneNumber": "+1-555-1234",
    "status": "active",
    "jobsCount": 15,
    "applicationsCount": 45,
    "createdAt": "2024-01-15T10:00:00Z",
    "updatedAt": "2024-01-15T10:00:00Z",
    "user": {
      "id": "uuid-string",
      "name": "John Doe",
      "email": "john@agritech.com",
      "role": "employer"
    }
  }
}
```

---

### 47. **PUT /api/admin/companies/:id**
**Description:** Update company information

**Request Body:**
```json
{
  "companyName": "Updated Company Name",
  "industry": "Crop Production",
  "companySize": "201-500 employees",
  "companyDescription": "Updated company description",
  "location": "Texas, USA",
  "websiteUrl": "https://updated-agritech.com",
  "recruiterName": "Updated Name",
  "officialEmail": "hr@updated-agritech.com",
  "phoneNumber": "+1-555-5678"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Company updated successfully"
}
```

---

### 48. **DELETE /api/admin/companies/:id**
**Description:** Delete company and all associated data

**Response:**
```json
{
  "success": true,
  "message": "Company deleted successfully"
}
```

---

## 📊 Error Responses

All endpoints return consistent error responses:

**400 Bad Request:**
```json
{
  "success": false,
  "message": "Invalid request"
}
```

**401 Unauthorized:**
```json
{
  "success": false,
  "message": "Missing or invalid token"
}
```

**404 Not Found:**
```json
{
  "success": false,
  "message": "Resource not found"
}
```

**500 Internal Server Error:**
```json
{
  "success": false,
  "message": "Internal server error"
}
```

---

## 🔧 Notes for Frontend Implementation

1. **Authentication:** Include JWT token in Authorization header for protected endpoints
2. **File Uploads:** Use multipart/form-data for file uploads
3. **Partial Updates:** For PUT endpoints, only send fields that need to be updated
4. **Error Handling:** Always check the `success` field in responses
5. **Pagination:** Some endpoints may need pagination for large datasets (not implemented yet)
6. **Real-time Updates:** Consider WebSocket implementation for real-time notifications

---

## 🚀 Testing

You can test these endpoints using tools like:
- Postman
- Insomnia
- curl
- Your frontend application

Make sure to:
1. Start the server: `go run ./cmd/server`
2. Use the correct base URL: `http://localhost:3000/api`
3. Include proper authentication headers where required 