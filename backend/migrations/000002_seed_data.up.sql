-- Author: Haru
-- Date: 2026-06-18
-- Task/Jira ID: BT-2
-- Purpose: Seed initial data matching the iOS App UI MockData for users and tutors in tutora_app schema

-- Clean up existing seeds
TRUNCATE TABLE tutora_app.tutors RESTART IDENTITY CASCADE;
DELETE FROM tutora_app.users;

-- Insert Users
INSERT INTO tutora_app.users (name, email, role) VALUES
('Haru Learner', 'haru@tutora.com', 'student'),
('Mina Park', 'mina@tutora.com', 'tutor'),
('Jay Chen', 'jay@tutora.com', 'tutor'),
('Ava Santos', 'ava@tutora.com', 'tutor')
ON CONFLICT (email) DO NOTHING;

-- Insert Tutors
-- Mina Park
INSERT INTO tutora_app.tutors (user_id, subject, bio, price_per_hour, rating)
SELECT id, 'Calculus, Physics', 'Ex-Olympiad coach turning hard math into sticky mental models.', 24.00, 4.90
FROM tutora_app.users
WHERE email = 'mina@tutora.com'
ON CONFLICT DO NOTHING;

-- Jay Chen
INSERT INTO tutora_app.tutors (user_id, subject, bio, price_per_hour, rating)
SELECT id, 'Python, AI', 'Python, data, and exam prep with tiny projects that compound.', 19.00, 4.80
FROM tutora_app.users
WHERE email = 'jay@tutora.com'
ON CONFLICT DO NOTHING;

-- Ava Santos
INSERT INTO tutora_app.tutors (user_id, subject, bio, price_per_hour, rating)
SELECT id, 'Biology, Chemistry', 'Memory coach for bio, chem, and spaced repetition systems.', 29.00, 4.95
FROM tutora_app.users
WHERE email = 'ava@tutora.com'
ON CONFLICT DO NOTHING;
