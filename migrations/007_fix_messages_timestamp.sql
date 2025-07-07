-- Fix messages table sent_at timestamp issue
-- This migration ensures the sent_at column has proper default value and constraints

-- Drop the existing messages table and recreate it with proper timestamp handling
DROP TABLE IF EXISTS public.messages CASCADE;

-- Recreate messages table with proper timestamp handling
CREATE TABLE IF NOT EXISTS public.messages (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    application_id uuid,
    sender_id uuid,
    message text COLLATE pg_catalog."default",
    sent_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
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

-- Add index for better performance
CREATE INDEX IF NOT EXISTS idx_messages_application_id ON messages(application_id);
CREATE INDEX IF NOT EXISTS idx_messages_sent_at ON messages(sent_at);

-- Add comment
COMMENT ON COLUMN messages.sent_at IS 'Timestamp when the message was sent (auto-set by database)'; 