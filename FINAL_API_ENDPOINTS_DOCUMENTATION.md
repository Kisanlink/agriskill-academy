# AgriJobs FINAL API Endpoints Documentation

## Base URL
```
http://localhost:3000/api
```

---

# 🔐 Authentication Endpoints

### 1. POST /api/auth/signup
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
**Error Response:**
```json
{
  "success": false,
  "message": "Invalid request"
}
```

### 2. POST /api/auth/login
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
**Error Response:**
```json
{
  "success": false,
  "message": "Invalid request"
}
```

### 3. PUT /api/auth/profile
**Description:** Update basic user information (name only; email cannot be changed)

**Request Body:**
```json
{
  "name": "Updated Name"
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
**Error Response:**
```json
{
  "success": false,
  "message": "Invalid request"
}
```

### 4. GET /api/auth/verify
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
**Error Response:**
```json
{
  "success": false,
  "message": "Missing or invalid token"
}
```

### 5. POST /api/auth/forgot-password
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
**Error Response:**
```json
{
  "success": false,
  "message": "Invalid request"
}
```

### 6. POST /api/auth/reset-password
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
**Error Response:**
```json
{
  "success": false,
  "message": "Invalid request"
}
```

// ... The rest of the 79 endpoints will be added here in the same detailed format ... 