-- Migration: Add username column to users table
-- This migration adds a username field to separate username from email

-- Add username column to users table
ALTER TABLE public.users 
ADD COLUMN IF NOT EXISTS username text COLLATE pg_catalog."default";

-- Add unique constraint on username
ALTER TABLE public.users 
ADD CONSTRAINT IF NOT EXISTS users_username_key UNIQUE (username);

-- Add index on username for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_username ON public.users (username);

-- Update existing users to have a username based on their email (temporary)
-- This will be handled by the application logic during login/signup
UPDATE public.users 
SET username = email 
WHERE username IS NULL;

-- Make username NOT NULL after setting default values
ALTER TABLE public.users 
ALTER COLUMN username SET NOT NULL;
