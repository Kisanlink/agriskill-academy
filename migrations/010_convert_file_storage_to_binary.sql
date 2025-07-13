-- Migration: Convert file storage from file paths to binary storage
-- Date: 2025-01-XX

-- Convert text columns to BYTEA with explicit casting
ALTER TABLE users ALTER COLUMN avatar TYPE BYTEA USING avatar::bytea;
ALTER TABLE employer_profiles ALTER COLUMN logo TYPE BYTEA USING logo::bytea;
ALTER TABLE student_profiles ALTER COLUMN profile_photo TYPE BYTEA USING profile_photo::bytea;
ALTER TABLE student_profiles ALTER COLUMN resume TYPE BYTEA USING resume::bytea;
ALTER TABLE applications ALTER COLUMN resume_file TYPE BYTEA USING resume_file::bytea;
ALTER TABLE certificates ALTER COLUMN file TYPE BYTEA USING file::bytea;

-- Add metadata columns for file information
ALTER TABLE users ADD COLUMN avatar_name VARCHAR(255);
ALTER TABLE users ADD COLUMN avatar_type VARCHAR(100);
ALTER TABLE users ADD COLUMN avatar_size BIGINT;

ALTER TABLE employer_profiles ADD COLUMN logo_name VARCHAR(255);
ALTER TABLE employer_profiles ADD COLUMN logo_type VARCHAR(100);
ALTER TABLE employer_profiles ADD COLUMN logo_size BIGINT;

ALTER TABLE student_profiles ADD COLUMN profile_photo_name VARCHAR(255);
ALTER TABLE student_profiles ADD COLUMN profile_photo_type VARCHAR(100);
ALTER TABLE student_profiles ADD COLUMN profile_photo_size BIGINT;

ALTER TABLE student_profiles ADD COLUMN resume_name VARCHAR(255);
ALTER TABLE student_profiles ADD COLUMN resume_type VARCHAR(100);
ALTER TABLE student_profiles ADD COLUMN resume_size BIGINT;

ALTER TABLE applications ADD COLUMN resume_file_name VARCHAR(255);
ALTER TABLE applications ADD COLUMN resume_file_type VARCHAR(100);
ALTER TABLE applications ADD COLUMN resume_file_size BIGINT;

ALTER TABLE certificates ADD COLUMN file_name VARCHAR(255);
ALTER TABLE certificates ADD COLUMN file_type VARCHAR(100);
ALTER TABLE certificates ADD COLUMN file_size BIGINT;

-- Add indexes for better performance on file metadata
CREATE INDEX IF NOT EXISTS idx_users_avatar_name ON users(avatar_name);
CREATE INDEX IF NOT EXISTS idx_employer_profiles_logo_name ON employer_profiles(logo_name);
CREATE INDEX IF NOT EXISTS idx_student_profiles_resume_name ON student_profiles(resume_name);
CREATE INDEX IF NOT EXISTS idx_certificates_file_name ON certificates(file_name);

-- Add comments to document the changes
COMMENT ON COLUMN users.avatar IS 'Binary file data for user avatar/profile photo';
COMMENT ON COLUMN users.avatar_name IS 'Original filename of the avatar';
COMMENT ON COLUMN users.avatar_type IS 'MIME type of the avatar file';
COMMENT ON COLUMN users.avatar_size IS 'File size in bytes';

COMMENT ON COLUMN employer_profiles.logo IS 'Binary file data for company logo';
COMMENT ON COLUMN employer_profiles.logo_name IS 'Original filename of the logo';
COMMENT ON COLUMN employer_profiles.logo_type IS 'MIME type of the logo file';
COMMENT ON COLUMN employer_profiles.logo_size IS 'File size in bytes';

COMMENT ON COLUMN student_profiles.profile_photo IS 'Binary file data for student profile photo';
COMMENT ON COLUMN student_profiles.profile_photo_name IS 'Original filename of the profile photo';
COMMENT ON COLUMN student_profiles.profile_photo_type IS 'MIME type of the profile photo file';
COMMENT ON COLUMN student_profiles.profile_photo_size IS 'File size in bytes';

COMMENT ON COLUMN student_profiles.resume IS 'Binary file data for student resume';
COMMENT ON COLUMN student_profiles.resume_name IS 'Original filename of the resume';
COMMENT ON COLUMN student_profiles.resume_type IS 'MIME type of the resume file';
COMMENT ON COLUMN student_profiles.resume_size IS 'File size in bytes';

COMMENT ON COLUMN applications.resume_file IS 'Binary file data for application resume';
COMMENT ON COLUMN applications.resume_file_name IS 'Original filename of the resume';
COMMENT ON COLUMN applications.resume_file_type IS 'MIME type of the resume file';
COMMENT ON COLUMN applications.resume_file_size IS 'File size in bytes';

COMMENT ON COLUMN certificates.file IS 'Binary file data for certificate';
COMMENT ON COLUMN certificates.file_name IS 'Original filename of the certificate';
COMMENT ON COLUMN certificates.file_type IS 'MIME type of the certificate file';
COMMENT ON COLUMN certificates.file_size IS 'File size in bytes'; 