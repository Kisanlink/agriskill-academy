-- Migration 012: Create job_hires table for tracking multiple hires per job

-- Create job_hires table if not exists
CREATE TABLE IF NOT EXISTS job_hires (
    id VARCHAR(255) PRIMARY KEY,
    job_id VARCHAR(255) NOT NULL,
    application_id VARCHAR(255) NOT NULL,
    candidate_name VARCHAR(255) NOT NULL,
    candidate_email VARCHAR(255) NOT NULL,
    student_id VARCHAR(255) NOT NULL,
    hired_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    deleted_by VARCHAR(255)
);

-- Create indexes for efficient queries
CREATE INDEX IF NOT EXISTS idx_job_hires_job ON job_hires(job_id);
CREATE INDEX IF NOT EXISTS idx_job_hires_student ON job_hires(student_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_job_hires_application ON job_hires(application_id);

-- Add foreign key constraints
ALTER TABLE job_hires
    ADD CONSTRAINT IF NOT EXISTS fk_job_hires_job
    FOREIGN KEY (job_id) REFERENCES job_posts(id) ON DELETE CASCADE;

ALTER TABLE job_hires
    ADD CONSTRAINT IF NOT EXISTS fk_job_hires_application
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE;

-- Add comment for documentation
COMMENT ON TABLE job_hires IS 'Tracks all hired candidates for job posts, supporting multiple hires per job';
COMMENT ON COLUMN job_hires.job_id IS 'Reference to the job post';
COMMENT ON COLUMN job_hires.application_id IS 'Reference to the application (unique constraint ensures one hire per application)';
COMMENT ON COLUMN job_hires.candidate_name IS 'Name of the hired candidate';
COMMENT ON COLUMN job_hires.candidate_email IS 'Email of the hired candidate';
COMMENT ON COLUMN job_hires.student_id IS 'Reference to the student user';
COMMENT ON COLUMN job_hires.hired_at IS 'Timestamp when the candidate was hired';
