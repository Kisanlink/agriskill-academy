-- Complete Database Schema Migration
-- This file contains all tables and structures matching the actual database schema

-- Enable UUID extension (Postgres)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- USERS TABLE
CREATE TABLE IF NOT EXISTS public.users (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name text COLLATE pg_catalog."default",
    email text COLLATE pg_catalog."default" NOT NULL,
    password text COLLATE pg_catalog."default" NOT NULL,
    role text COLLATE pg_catalog."default" NOT NULL,
    avatar text COLLATE pg_catalog."default",
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_email_key UNIQUE (email)
);

-- EMPLOYER PROFILES TABLE
CREATE TABLE IF NOT EXISTS public.employer_profiles (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    company_name text COLLATE pg_catalog."default",
    logo text COLLATE pg_catalog."default",
    website_url text COLLATE pg_catalog."default",
    industry text COLLATE pg_catalog."default",
    company_size text COLLATE pg_catalog."default",
    company_description text COLLATE pg_catalog."default",
    recruiter_name text COLLATE pg_catalog."default",
    designation text COLLATE pg_catalog."default",
    official_email text COLLATE pg_catalog."default",
    phone_number text COLLATE pg_catalog."default",
    linkedin_profile text COLLATE pg_catalog."default",
    job_categories text[] COLLATE pg_catalog."default",
    hiring_locations text[] COLLATE pg_catalog."default",
    hiring_types text[] COLLATE pg_catalog."default",
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    gstin_number text COLLATE pg_catalog."default",
    company_address text COLLATE pg_catalog."default",
    city text COLLATE pg_catalog."default",
    state text COLLATE pg_catalog."default",
    pincode text COLLATE pg_catalog."default",
    CONSTRAINT employer_profiles_pkey PRIMARY KEY (id),
    CONSTRAINT employer_profiles_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

-- JOB POSTS TABLE
CREATE TABLE IF NOT EXISTS public.job_posts (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    employer_id uuid,
    title text COLLATE pg_catalog."default",
    role_overview text COLLATE pg_catalog."default",
    requirements text COLLATE pg_catalog."default",
    location text COLLATE pg_catalog."default",
    required_skills text[] COLLATE pg_catalog."default",
    employer_name text COLLATE pg_catalog."default",
    employer_email text COLLATE pg_catalog."default",
    status text COLLATE pg_catalog."default",
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    application_deadline date,
    job_type text COLLATE pg_catalog."default",
    experience text COLLATE pg_catalog."default",
    salary_min numeric,
    salary_max numeric,
    salary_currency text COLLATE pg_catalog."default",
    completed_at timestamp with time zone,
    hired_candidate_name text COLLATE pg_catalog."default",
    benefits text[] COLLATE pg_catalog."default",
    is_remote boolean DEFAULT false,
    applications_count integer DEFAULT 0,
    CONSTRAINT job_posts_pkey PRIMARY KEY (id),
    CONSTRAINT job_posts_employer_id_fkey FOREIGN KEY (employer_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

-- APPLICATIONS TABLE
CREATE TABLE IF NOT EXISTS public.applications (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    job_id uuid,
    student_id uuid,
    applied_at timestamp with time zone DEFAULT now(),
    status text COLLATE pg_catalog."default",
    cover_letter text COLLATE pg_catalog."default",
    resume_file text COLLATE pg_catalog."default",
    job_title text COLLATE pg_catalog."default",
    company text COLLATE pg_catalog."default",
    location text COLLATE pg_catalog."default",
    job_type text COLLATE pg_catalog."default",
    experience text COLLATE pg_catalog."default",
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT applications_pkey PRIMARY KEY (id),
    CONSTRAINT applications_job_id_fkey FOREIGN KEY (job_id)
        REFERENCES public.job_posts (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT applications_student_id_fkey FOREIGN KEY (student_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

-- USER PROFILES TABLE
CREATE TABLE IF NOT EXISTS public.user_profiles (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid,
    name text COLLATE pg_catalog."default",
    email text COLLATE pg_catalog."default",
    location text COLLATE pg_catalog."default",
    profile_photo text COLLATE pg_catalog."default",
    resume text COLLATE pg_catalog."default",
    skills text[] COLLATE pg_catalog."default",
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    experience double precision,
    education text COLLATE pg_catalog."default",
    portfolio text COLLATE pg_catalog."default",
    linkedin text COLLATE pg_catalog."default",
    github text COLLATE pg_catalog."default",
    CONSTRAINT user_profiles_pkey PRIMARY KEY (id),
    CONSTRAINT user_profiles_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

-- BOOKMARKS TABLE
CREATE TABLE IF NOT EXISTS public.bookmarks (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid,
    job_id uuid,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT bookmarks_pkey PRIMARY KEY (id),
    CONSTRAINT bookmarks_job_id_fkey FOREIGN KEY (job_id)
        REFERENCES public.job_posts (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT bookmarks_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

-- CERTIFICATES TABLE
CREATE TABLE IF NOT EXISTS public.certificates (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_profile_id uuid,
    name text COLLATE pg_catalog."default",
    file text COLLATE pg_catalog."default",
    issue_date date,
    CONSTRAINT certificates_pkey PRIMARY KEY (id),
    CONSTRAINT certificates_user_profile_id_fkey FOREIGN KEY (user_profile_id)
        REFERENCES public.user_profiles (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

-- MESSAGES TABLE
CREATE TABLE IF NOT EXISTS public.messages (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    application_id uuid,
    sender_id uuid,
    message text COLLATE pg_catalog."default",
    sent_at timestamp with time zone DEFAULT now(),
    CONSTRAINT messages_pkey PRIMARY KEY (id),
    CONSTRAINT messages_application_id_fkey FOREIGN KEY (application_id)
        REFERENCES public.applications (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT messages_sender_id_fkey FOREIGN KEY (sender_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

-- NOTIFICATION PREFERENCES TABLE
CREATE TABLE IF NOT EXISTS public.notification_preferences (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    email_notifications boolean DEFAULT true,
    push_notifications boolean DEFAULT true,
    job_alerts boolean DEFAULT true,
    application_updates boolean DEFAULT true,
    company_news boolean DEFAULT false,
    marketing_emails boolean DEFAULT false,
    weekly_digest boolean DEFAULT true,
    daily_job_matches boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT notification_preferences_pkey PRIMARY KEY (id),
    CONSTRAINT notification_preferences_user_id_key UNIQUE (user_id),
    CONSTRAINT notification_preferences_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

-- JOB ALERTS TABLE
CREATE TABLE IF NOT EXISTS public.job_alerts (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    keywords text[] COLLATE pg_catalog."default",
    location character varying(255) COLLATE pg_catalog."default",
    job_type text[] COLLATE pg_catalog."default",
    experience text[] COLLATE pg_catalog."default",
    skills text[] COLLATE pg_catalog."default",
    salary_min numeric(10,2),
    salary_max numeric(10,2),
    salary_currency character varying(10) COLLATE pg_catalog."default" DEFAULT 'USD'::character varying,
    is_remote boolean,
    frequency character varying(20) COLLATE pg_catalog."default" NOT NULL DEFAULT 'weekly'::character varying,
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT job_alerts_pkey PRIMARY KEY (id),
    CONSTRAINT job_alerts_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT job_alerts_frequency_check CHECK (frequency::text = ANY (ARRAY['daily'::character varying, 'weekly'::character varying, 'immediate'::character varying]::text[]))
);

-- CREATE INDEXES
CREATE INDEX IF NOT EXISTS idx_job_posts_employer_id ON job_posts(employer_id);
CREATE INDEX IF NOT EXISTS idx_applications_job_id ON applications(job_id);
CREATE INDEX IF NOT EXISTS idx_applications_student_id ON applications(student_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_user_id ON bookmarks(user_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_job_id ON bookmarks(job_id);
CREATE INDEX IF NOT EXISTS idx_notification_preferences_user_id ON notification_preferences(user_id);
CREATE INDEX IF NOT EXISTS idx_job_alerts_user_id ON job_alerts(user_id);
CREATE INDEX IF NOT EXISTS idx_job_alerts_is_active ON job_alerts(is_active);
CREATE INDEX IF NOT EXISTS idx_job_alerts_frequency ON job_alerts(frequency);
CREATE INDEX IF NOT EXISTS idx_job_posts_applications_count ON job_posts(applications_count);
CREATE INDEX IF NOT EXISTS idx_job_posts_is_remote ON job_posts(is_remote);
CREATE INDEX IF NOT EXISTS idx_job_posts_salary_range ON job_posts(salary_min, salary_max);

-- CREATE TRIGGERS AND FUNCTIONS
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE OR REPLACE FUNCTION update_job_alerts_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- CREATE TRIGGERS
CREATE TRIGGER update_applications_updated_at
    BEFORE UPDATE ON applications
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_update_job_alerts_updated_at
    BEFORE UPDATE ON job_alerts
    FOR EACH ROW
    EXECUTE FUNCTION update_job_alerts_updated_at();

-- ADD COMMENTS
COMMENT ON COLUMN applications.updated_at IS 'Timestamp when the application was last updated';
COMMENT ON COLUMN job_posts.salary_min IS 'Minimum salary for the job position';
COMMENT ON COLUMN job_posts.salary_max IS 'Maximum salary for the job position';
COMMENT ON COLUMN job_posts.salary_currency IS 'Currency code for salary (USD, EUR, etc.)';
COMMENT ON COLUMN job_posts.benefits IS 'Array of benefits offered with the job';
COMMENT ON COLUMN job_posts.is_remote IS 'Whether the job allows remote work';
COMMENT ON COLUMN job_posts.applications_count IS 'Number of applications received for this job';
COMMENT ON TABLE job_alerts IS 'Job alerts for users to get notified about new jobs matching their criteria'; 