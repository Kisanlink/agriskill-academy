-- Add Contact Requests Table Migration
-- This migration creates a table to store contact form submissions

CREATE TABLE IF NOT EXISTS public.contact_requests (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    first_name text COLLATE pg_catalog."default" NOT NULL,
    last_name text COLLATE pg_catalog."default" NOT NULL,
    email text COLLATE pg_catalog."default" NOT NULL,
    phone text COLLATE pg_catalog."default" NOT NULL,
    subject text COLLATE pg_catalog."default" NOT NULL,
    message text COLLATE pg_catalog."default" NOT NULL,
    status text COLLATE pg_catalog."default" DEFAULT 'new',
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT contact_requests_pkey PRIMARY KEY (id)
);

-- Add index for better query performance on admin panel
CREATE INDEX IF NOT EXISTS idx_contact_requests_status ON public.contact_requests(status);
CREATE INDEX IF NOT EXISTS idx_contact_requests_created_at ON public.contact_requests(created_at DESC); 