-- Migration: Remove obsolete file/binary columns from student_profiles and employerapplication tables
-- Up
ALTER TABLE student_profiles
  DROP COLUMN IF EXISTS profile_photo,
  DROP COLUMN IF EXISTS profile_photo_name,
  DROP COLUMN IF EXISTS profile_photo_type,
  DROP COLUMN IF EXISTS profile_photo_size,
  DROP COLUMN IF EXISTS resume,
  DROP COLUMN IF EXISTS resume_name,
  DROP COLUMN IF EXISTS resume_type,
  DROP COLUMN IF EXISTS resume_size;

ALTER TABLE employer_applications
  DROP COLUMN IF EXISTS student_resume_file,
  DROP COLUMN IF EXISTS avatar;

-- Down (optional: add columns back, but types may need to be adjusted)
ALTER TABLE student_profiles
  ADD COLUMN profile_photo bytea,
  ADD COLUMN profile_photo_name text,
  ADD COLUMN profile_photo_type text,
  ADD COLUMN profile_photo_size bigint,
  ADD COLUMN resume bytea,
  ADD COLUMN resume_name text,
  ADD COLUMN resume_type text,
  ADD COLUMN resume_size bigint;

ALTER TABLE employer_applications
  ADD COLUMN student_resume_file bytea,
  ADD COLUMN avatar bytea; 