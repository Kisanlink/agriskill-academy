-- Migration: Add phone_number column to student_profiles table
-- Date: 2025-07-10

-- Add phone_number column to student_profiles table
ALTER TABLE student_profiles 
ADD COLUMN phone_number VARCHAR(20);

-- Add comment to document the column
COMMENT ON COLUMN student_profiles.phone_number IS 'Phone number from AAA service registration'; 