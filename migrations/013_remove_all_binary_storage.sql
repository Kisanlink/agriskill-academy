-- Migration: Remove all binary storage and ensure S3 key storage only
-- This migration ensures all file storage uses S3 keys instead of binary data

-- Remove any remaining binary columns and ensure S3 key columns exist

-- Users table - ensure avatar_key exists and remove any binary avatar
ALTER TABLE users 
DROP COLUMN IF EXISTS avatar,
DROP COLUMN IF EXISTS avatar_name,
DROP COLUMN IF EXISTS avatar_type,
DROP COLUMN IF EXISTS avatar_size,
ADD COLUMN IF NOT EXISTS avatar_key TEXT;

-- Employer profiles table - ensure logo_key exists and remove any binary logo
ALTER TABLE employer_profiles 
DROP COLUMN IF EXISTS logo,
DROP COLUMN IF EXISTS logo_name,
DROP COLUMN IF EXISTS logo_type,
DROP COLUMN IF EXISTS logo_size,
ADD COLUMN IF NOT EXISTS logo_key TEXT;

-- Student profiles table - ensure S3 key columns exist and remove any binary data
ALTER TABLE student_profiles 
DROP COLUMN IF EXISTS profile_photo,
DROP COLUMN IF EXISTS profile_photo_name,
DROP COLUMN IF EXISTS profile_photo_type,
DROP COLUMN IF EXISTS profile_photo_size,
DROP COLUMN IF EXISTS resume,
DROP COLUMN IF EXISTS resume_name,
DROP COLUMN IF EXISTS resume_type,
DROP COLUMN IF EXISTS resume_size,
ADD COLUMN IF NOT EXISTS profile_photo_key TEXT,
ADD COLUMN IF NOT EXISTS resume_key TEXT;

-- Applications table - ensure resume_key exists and remove any binary resume_file
ALTER TABLE applications 
DROP COLUMN IF EXISTS resume_file,
DROP COLUMN IF EXISTS resume_file_name,
DROP COLUMN IF EXISTS resume_file_type,
DROP COLUMN IF EXISTS resume_file_size,
ADD COLUMN IF NOT EXISTS resume_key TEXT;

-- Certificates table - ensure file_key exists and remove any binary file
ALTER TABLE certificates 
DROP COLUMN IF EXISTS file,
DROP COLUMN IF EXISTS file_name,
DROP COLUMN IF EXISTS file_type,
DROP COLUMN IF EXISTS file_size,
ADD COLUMN IF NOT EXISTS file_key TEXT;

-- User profiles table (legacy) - remove any binary columns
ALTER TABLE user_profiles 
DROP COLUMN IF EXISTS profile_photo,
DROP COLUMN IF EXISTS resume;

-- Add comments to document the S3 key columns
COMMENT ON COLUMN users.avatar_key IS 'S3 key for user avatar/profile photo';
COMMENT ON COLUMN employer_profiles.logo_key IS 'S3 key for company logo';
COMMENT ON COLUMN student_profiles.profile_photo_key IS 'S3 key for student profile photo';
COMMENT ON COLUMN student_profiles.resume_key IS 'S3 key for student resume';
COMMENT ON COLUMN applications.resume_key IS 'S3 key for application resume';
COMMENT ON COLUMN certificates.file_key IS 'S3 key for certificate file';

-- Create indexes for S3 key columns for better performance
CREATE INDEX IF NOT EXISTS idx_users_avatar_key ON users(avatar_key);
CREATE INDEX IF NOT EXISTS idx_employer_profiles_logo_key ON employer_profiles(logo_key);
CREATE INDEX IF NOT EXISTS idx_student_profiles_profile_photo_key ON student_profiles(profile_photo_key);
CREATE INDEX IF NOT EXISTS idx_student_profiles_resume_key ON student_profiles(resume_key);
CREATE INDEX IF NOT EXISTS idx_applications_resume_key ON applications(resume_key);
CREATE INDEX IF NOT EXISTS idx_certificates_file_key ON certificates(file_key); 