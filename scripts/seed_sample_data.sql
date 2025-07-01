-- File: scripts/seed_sample_data.sql

INSERT INTO users (id, name, email, password, role)
VALUES (uuid_generate_v4(), 'Employer A', 'employerA@example.com', 'hashedpassword', 'employer'),
       (uuid_generate_v4(), 'Student B', 'studentB@example.com', 'hashedpassword', 'student');
