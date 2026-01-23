# Frontend API Catalog - AgriJobs Platform
**Last Updated:** January 8, 2026
**Backend Version:** v2.1.0
**Breaking Changes:** None - All changes are backward compatible

---

## 🚀 New Features (January 8, 2026)

### Multiple Hires Tracking System
Jobs can now have multiple hired candidates without data loss. The old `hired_candidate_name` field is still populated for backward compatibility.

### Manual Job Close/Reopen
Employers can now manually control job lifecycle without automatic closures.

---

## 📋 Table of Contents
1. [New Endpoints](#new-endpoints)
2. [Modified Behavior](#modified-behavior)
3. [Backward Compatibility](#backward-compatibility)
4. [Migration Guide](#migration-guide)

---

## 🆕 New Endpoints

### 1. Get Hired Candidates for a Job

**Endpoint:** `GET /api/jobs/:id/hires`

**Description:** Retrieves all hired candidates for a specific job post (Public endpoint)

**Authentication:** None required

**Response (200 OK):**
\`\`\`json
{
  "success": true,
  "message": "Hired candidates retrieved successfully",
  "hires": [
    {
      "id": "HIRE_1766581132397042100",
      "job_id": "JOBP_1234567890",
      "application_id": "APP_9876543210",
      "candidate_name": "John Doe",
      "candidate_email": "john.doe@example.com",
      "student_id": "STU_5555555555",
      "hired_at": "2026-01-05T14:30:00Z",
      "created_at": "2026-01-05T14:30:00Z"
    }
  ],
  "count": 1
}
\`\`\`

**Use Cases:**
- Display list of all hired candidates on job detail page
- Show hiring history for completed jobs
- Track multiple hires for the same position

---

### 2. Close Job Post (Employer Only)

**Endpoint:** `POST /api/jobs/:id/close`

**Description:** Manually close a job post. Only the employer who created the job can close it.

**Authentication:** Required (Bearer token)

**Authorization:** Employer role + Job ownership

**Response (200 OK):**
\`\`\`json
{
  "success": true,
  "message": "Job closed successfully",
  "job_id": "JOBP_1234567890"
}
\`\`\`

**Use Cases:**
- Close job after hiring sufficient candidates
- Close job if position is cancelled
- Close job due to budget constraints

**Database Changes:**
- Sets status to "completed"
- Sets completed_at timestamp

---

### 3. Reopen Job Post (Employer Only)

**Endpoint:** `POST /api/jobs/:id/reopen`

**Description:** Reopen a previously closed job post. Only the employer who created the job can reopen it.

**Authentication:** Required (Bearer token)

**Authorization:** Employer role + Job ownership

**Response (200 OK):**
\`\`\`json
{
  "success": true,
  "message": "Job reopened successfully",
  "job_id": "JOBP_1234567890"
}
\`\`\`

**Use Cases:**
- Reopen job if budget is approved
- Reopen job if selected candidate declines
- Reopen job if more candidates needed

**Database Changes:**
- Sets status to "published"
- Job becomes visible in active listings

---

## 🔄 Modified Behavior

### Application Acceptance - NO AUTO-CLOSE

**⚠️ IMPORTANT CHANGE:** Jobs no longer automatically close when an application is accepted.

**Previous Behavior:**
- Job automatically moved to "completed" status after first hire
- Only one candidate could be hired

**New Behavior:**
- Job remains "open" after hiring candidates
- Multiple candidates can be hired for the same job
- Employers must manually close jobs via POST /api/jobs/:id/close

**Example Flow:**
\`\`\`
1. Employer posts job → status: "published"
2. Employer accepts application → status: "published" (no auto-close)
3. Employer accepts another application → status: "published" (multiple hires tracked)
4. Employer manually closes job → status: "completed"
\`\`\`

---

## ✅ Backward Compatibility

### Guaranteed Compatibility

**No Breaking Changes:**
- All existing API endpoints work exactly as before
- Existing request/response structures unchanged
- Frontend code requires NO modifications

**Preserved Fields:**
- hired_candidate_name field STILL EXISTS
- hired_candidate_name is STILL POPULATED with first hired candidate

**Example Job Response:**
\`\`\`json
{
  "id": "JOBP_1234567890",
  "title": "Farm Manager",
  "hired_candidate_name": "John Doe",
  "status": "open"
}
\`\`\`

### What Frontend Can Continue Doing:
✅ Display hired_candidate_name field
✅ Check status field for job state
✅ Use all existing job endpoints

### What's New (Optional):
🆕 Display multiple hired candidates using /api/jobs/:id/hires
🆕 Add "Close Job" button calling POST /api/jobs/:id/close
🆕 Add "Reopen Job" button calling POST /api/jobs/:id/reopen

---

## ⚠️ Important Notes

### For Employers
- Jobs don't auto-close - you must manually close them
- Multiple candidates can be hired for the same job
- Closed jobs can be reopened anytime

---

## 📝 Summary

| Feature | Status | Action Required |
|---------|--------|----------------|
| hired_candidate_name field | ✅ Preserved | None |
| Multiple hires tracking | ✅ Implemented | Optional: Add UI |
| Manual job close | ✅ Implemented | Optional: Add button |
| Manual job reopen | ✅ Implemented | Optional: Add button |
| Existing APIs | ✅ Unchanged | None |

---

## 🔐 Admin Endpoints (January 8, 2026)

### Admin Job Viewing

Administrators (asa_admin role) can view all jobs in the system regardless of status. Admin users have read-only access and cannot create, edit, delete, or apply to jobs.

---

### 1. List All Jobs (Admin Only)

**Endpoint:** `GET /api/admin/jobs`

**Description:** Retrieve all jobs in the system regardless of status (draft, published, completed, closed)

**Authentication:** Required (Bearer token)

**Authorization:** asa_admin role only

**Query Parameters:**
- `status` (optional): Filter by status (draft, published, completed, closed)
- `employer_id` (optional): Filter by specific employer
- `page` (optional): Page number (default: 1)
- `limit` (optional): Results per page (default: 50)
- `sort_by` (optional): Sort by field
- `sort_order` (optional): Sort order (asc, desc)

**Response (200 OK):**
\`\`\`json
{
  "success": true,
  "message": "Jobs retrieved successfully",
  "data": {
    "jobs": [
      {
        "id": "JOBP_123",
        "title": "Farm Manager",
        "status": "draft",
        "employer_id": "USR_456",
        "employer_name": "Green Farms Ltd",
        "location": "Pune, Maharashtra",
        "job_type": "Full-time",
        "applications_count": 5,
        "created_at": "2026-01-01T00:00:00Z",
        "updated_at": "2026-01-05T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 50,
      "total": 25,
      "total_pages": 1
    }
  }
}
\`\`\`

**Use Cases:**
- Admin dashboard showing all jobs
- Monitor job postings across platform
- Support users with job-related issues
- View jobs in all states (including drafts)

---

### 2. Get Job Details (Admin Only)

**Endpoint:** `GET /api/admin/jobs/:id`

**Description:** Retrieve complete details for a specific job including employer information

**Authentication:** Required (Bearer token)

**Authorization:** asa_admin role only

**Response (200 OK):**
\`\`\`json
{
  "success": true,
  "message": "Job details retrieved successfully",
  "data": {
    "id": "JOBP_123",
    "title": "Farm Manager",
    "description": "We are looking for an experienced farm manager...",
    "status": "draft",
    "employer_id": "USR_456",
    "employer_name": "Green Farms Ltd",
    "employer_email": "employer@greenfarms.com",
    "location": "Pune, Maharashtra",
    "job_type": "Full-time",
    "salary": "₹25,000 - ₹35,000 per month",
    "requirements": "Bachelor's in Agriculture, 2+ years experience",
    "responsibilities": "Manage daily farm operations...",
    "benefits": "Health insurance, accommodation",
    "applications_count": 5,
    "hired_candidate_name": "John Doe",
    "completed_at": null,
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-05T10:00:00Z"
  }
}
\`\`\`

**Response (404 Not Found):**
\`\`\`json
{
  "success": false,
  "message": "Job not found"
}
\`\`\`

**Use Cases:**
- View complete job information for support
- Inspect job details including employer contact
- Monitor job status and completion

---

### 3. Get Job Statistics (Admin Only)

**Endpoint:** `GET /api/admin/jobs/statistics`

**Description:** Retrieve aggregated statistics about jobs across the platform

**Authentication:** Required (Bearer token)

**Authorization:** asa_admin role only

**Response (200 OK):**
\`\`\`json
{
  "success": true,
  "message": "Job statistics retrieved successfully",
  "data": {
    "total_jobs": 150,
    "draft_jobs": 20,
    "published_jobs": 100,
    "completed_jobs": 30,
    "total_applications": 500,
    "total_hires": 45
  }
}
\`\`\`

**Use Cases:**
- Admin dashboard metrics
- Platform health monitoring
- Business intelligence and reporting

---

### Admin Access Restrictions

**What Admins CAN Do:**
✅ View all jobs (any status)
✅ View job details and employer information
✅ Filter and search jobs
✅ View job statistics

**What Admins CANNOT Do:**
❌ Create new jobs
❌ Edit existing jobs
❌ Delete jobs
❌ Apply to jobs
❌ Close/reopen jobs
❌ Accept/reject applications

**Security:**
- All admin endpoints require valid JWT token
- All admin endpoints require asa_admin role
- Non-admin users receive 403 Forbidden
- Read-only access enforced at route level

---

**Document End - All changes are production-ready and backward compatible**
