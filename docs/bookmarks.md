# Bookmarks Integration Guide

## Overview
The bookmark (save job) functionality is currently stubbed out in the frontend. This document provides instructions to enable the feature by uncommenting and updating the API calls.

## Backend API Endpoints

The backend has these endpoints ready:

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/api/bookmarks` | Get user's saved jobs | Yes (student) |
| `POST` | `/api/bookmarks/:jobId` | Save a job | Yes (student) |
| `DELETE` | `/api/bookmarks/:jobId` | Remove a saved job | Yes (student) |

## Frontend Changes Required

### File: `src/services/jobService.ts`

### 1. Save Job (line ~774)

**Current (stubbed):**
```typescript
async saveJob(jobId: string): Promise<{ success: boolean; message: string }> {
  try {
    // TODO: Replace with actual API call
    // const response = await fetch(`${this.baseURL}/api/jobs/${jobId}/save`, {
    //   method: 'POST',
    //   headers: this.getAuthHeaders(),
    // });
    // return await response.json();

    return { success: true, message: 'Job saved successfully!' };
  } catch (error) {
    console.error('Save job error:', error);
    return { success: false, message: 'Failed to save job.' };
  }
}
```

**Update to:**
```typescript
async saveJob(jobId: string): Promise<{ success: boolean; message: string }> {
  try {
    const response = await fetch(`${this.baseURL}/api/bookmarks/${jobId}`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
    });

    if (!response.ok) {
      const error = await response.json();
      return { success: false, message: error.message || 'Failed to save job.' };
    }

    return { success: true, message: 'Job saved successfully!' };
  } catch (error) {
    console.error('Save job error:', error);
    return { success: false, message: 'Failed to save job.' };
  }
}
```

---

### 2. Remove Saved Job (line ~791)

**Current (stubbed):**
```typescript
async removeSavedJob(jobId: string): Promise<{ success: boolean; message: string }> {
  try {
    // TODO: Replace with actual API call
    // const response = await fetch(`${this.baseURL}/api/jobs/${jobId}/unsave`, {
    //   method: 'DELETE',
    //   headers: this.getAuthHeaders(),
    // });
    // return await response.json();

    return { success: true, message: 'Job removed from saved jobs!' };
  } catch (error) {
    console.error('Remove saved job error:', error);
    return { success: false, message: 'Failed to remove saved job.' };
  }
}
```

**Update to:**
```typescript
async removeSavedJob(jobId: string): Promise<{ success: boolean; message: string }> {
  try {
    const response = await fetch(`${this.baseURL}/api/bookmarks/${jobId}`, {
      method: 'DELETE',
      headers: this.getAuthHeaders(),
    });

    if (!response.ok) {
      const error = await response.json();
      return { success: false, message: error.message || 'Failed to remove saved job.' };
    }

    return { success: true, message: 'Job removed from saved jobs!' };
  } catch (error) {
    console.error('Remove saved job error:', error);
    return { success: false, message: 'Failed to remove saved job.' };
  }
}
```

---

### 3. Get Saved Jobs (line ~808)

**Current (stubbed):**
```typescript
async getSavedJobs(): Promise<{ success: boolean; jobs: Job[] }> {
  try {
    // TODO: Replace with actual API call
    // const response = await fetch(`${this.baseURL}/api/jobs/saved`, {
    //   method: 'GET',
    //   headers: this.getAuthHeaders(),
    // });
    // return await response.json();

    // Mock implementation - return empty array for now
    return { success: true, jobs: [] };
  } catch (error) {
    console.error('Get saved jobs error:', error);
    return { success: false, jobs: [] };
  }
}
```

**Update to:**
```typescript
async getSavedJobs(): Promise<{ success: boolean; jobs: Job[] }> {
  try {
    const response = await fetch(`${this.baseURL}/api/bookmarks`, {
      method: 'GET',
      headers: this.getAuthHeaders(),
    });

    if (!response.ok) {
      return { success: false, jobs: [] };
    }

    const data = await response.json();
    return { success: true, jobs: data.jobs || data.data || [] };
  } catch (error) {
    console.error('Get saved jobs error:', error);
    return { success: false, jobs: [] };
  }
}
```

---

## URL Mapping Summary

| Frontend (commented) | Backend (correct) |
|---------------------|-------------------|
| `/api/jobs/${jobId}/save` | `/api/bookmarks/${jobId}` |
| `/api/jobs/${jobId}/unsave` | `/api/bookmarks/${jobId}` |
| `/api/jobs/saved` | `/api/bookmarks` |

## Notes

- Bookmarks are **student-only** functionality (requires student role)
- The backend returns the full job details when fetching bookmarks
- Make sure the user is authenticated before calling these endpoints
