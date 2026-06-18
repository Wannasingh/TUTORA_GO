-- Author: Haru
-- Date: 2026-06-18
-- Task/Jira ID: BT-2
-- Purpose: Seed initial data for users and tutors in tutora_app schema

-- Insert Users
INSERT INTO tutora_app.users (name, email, role) VALUES
('Alice Johnson', 'alice@bytestutor.com', 'student'),
('Bob Smith', 'bob@bytestutor.com', 'tutor'),
('Charlie Brown', 'charlie@bytestutor.com', 'tutor')
ON CONFLICT (email) DO NOTHING;

-- Insert Tutors (matching user IDs)
-- Bob Smith (tutor)
INSERT INTO tutora_app.tutors (user_id, subject, bio, price_per_hour, rating)
SELECT id, 'Go programming language', 'Senior Go Engineer with 8 years of experience building high-throughput systems.', 45.00, 4.8
FROM tutora_app.users
WHERE email = 'bob@bytestutor.com'
ON CONFLICT DO NOTHING;

-- Charlie Brown (tutor)
INSERT INTO tutora_app.tutors (user_id, subject, bio, price_per_hour, rating)
SELECT id, 'iOS Native Development with SwiftUI', 'iOS Developer specializing in SwiftUI, Combine, and premium UI animations.', 50.00, 4.9
FROM tutora_app.users
WHERE email = 'charlie@bytestutor.com'
ON CONFLICT DO NOTHING;
