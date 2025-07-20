-- Migration: Convert binary storage to S3 key storage
-- This migration converts BYTEA columns to TEXT columns for S3 keys

-- Convert users.avatar from BYTEA to TEXT for S3 key storage
ALTER TABLE users 
DROP COLUMN IF EXISTS avatar,
ADD COLUMN avatar_key TEXT;

-- Convert employer_profiles.logo from BYTEA to TEXT for S3 key storage
ALTER TABLE employer_profiles 
DROP COLUMN IF EXISTS logo,
ADD COLUMN logo_key TEXT;

-- Convert student_profiles.profile_photo from BYTEA to TEXT for S3 key storage
-- Note: profile_photo_key already exists, just ensure it's properly set
ALTER TABLE student_profiles 
DROP COLUMN IF EXISTS profile_photo;

-- Convert student_profiles.resume from BYTEA to TEXT for S3 key storage
-- Note: resume_key already exists, just ensure it's properly set
ALTER TABLE student_profiles 
DROP COLUMN IF EXISTS resume;

-- Convert applications.resume_file from BYTEA to TEXT for S3 key storage
ALTER TABLE applications 
DROP COLUMN IF EXISTS resume_file,
ADD COLUMN resume_key TEXT;

-- Convert certificates.file from BYTEA to TEXT for S3 key storage
ALTER TABLE certificates 
DROP COLUMN IF EXISTS file,
ADD COLUMN file_key TEXT;

-- Add comments to document the new S3 key columns
COMMENT ON COLUMN users.avatar_key IS 'S3 key for user avatar/profile photo';
COMMENT ON COLUMN employer_profiles.logo_key IS 'S3 key for company logo';
COMMENT ON COLUMN student_profiles.profile_photo_key IS 'S3 key for student profile photo';
COMMENT ON COLUMN student_profiles.resume_key IS 'S3 key for student resume';
COMMENT ON COLUMN applications.resume_key IS 'S3 key for application resume';
COMMENT ON COLUMN certificates.file_key IS 'S3 key for certificate file'; 