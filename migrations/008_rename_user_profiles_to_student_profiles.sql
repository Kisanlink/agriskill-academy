-- Migration: Rename user_profiles table to student_profiles
-- This migration renames the user_profiles table to student_profiles
-- and updates the foreign key in certificates table

-- Step 1: Rename the table
ALTER TABLE user_profiles RENAME TO student_profiles;

-- Step 2: Update the foreign key column in certificates table
ALTER TABLE certificates RENAME COLUMN user_profile_id TO student_profile_id;

-- Step 3: Update the foreign key constraint
ALTER TABLE certificates DROP CONSTRAINT IF EXISTS certificates_user_profile_id_fkey;
ALTER TABLE certificates ADD CONSTRAINT certificates_student_profile_id_fkey 
    FOREIGN KEY (student_profile_id) REFERENCES student_profiles(id) ON DELETE CASCADE;

-- Step 4: Update any indexes that reference the old column name
-- (PostgreSQL will automatically update indexes when we rename the column)

-- Step 5: Update any sequences or triggers if they exist
-- (In this case, we don't have any custom sequences or triggers)

-- Verify the changes
-- SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'student_profiles';
-- SELECT column_name FROM information_schema.columns WHERE table_name = 'certificates' AND column_name = 'student_profile_id'; 