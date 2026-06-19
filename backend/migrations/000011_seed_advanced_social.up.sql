-- Author: Antigravity
-- Date: 2026-06-19
-- Task/Jira ID: BT-11
-- Purpose: Seed advanced posts, videos, nested replies, comment likes, and reports

-- Clean up existing social seeds
TRUNCATE TABLE tutora_app.reports RESTART IDENTITY CASCADE;
TRUNCATE TABLE tutora_app.comment_likes RESTART IDENTITY CASCADE;
DELETE FROM tutora_app.comments WHERE post_id IN (SELECT id FROM tutora_app.posts WHERE subject IN ('Calculus', 'Python', 'Moderation Test'));
DELETE FROM tutora_app.post_likes WHERE post_id IN (SELECT id FROM tutora_app.posts WHERE subject IN ('Calculus', 'Python', 'Moderation Test'));
DELETE FROM tutora_app.posts WHERE subject IN ('Calculus', 'Python', 'Moderation Test');

-- 1. Insert Video Post by Mina Park (Tutor)
INSERT INTO tutora_app.posts (user_id, subject, title, body, video_url)
SELECT id, 'Calculus', 'Visualizing Limits Dynamically', 'Heres a quick video explaining how to approach limits conceptually without memorizing formulas.', 'https://tutora-endpoint-axbabq3egzii.private.compat.objectstorage.ap-osaka-1.oci.customer-oci.com/TUTORA/limits-video.mp4'
FROM tutora_app.users WHERE email = 'mina@tutora.com';

-- 2. Insert Image/GIF Post by Jay Chen (Tutor)
INSERT INTO tutora_app.posts (user_id, subject, title, body, image_url)
SELECT id, 'Python', 'OOP Cheat Sheet & GIF Guide', 'A summary diagram explaining encapsulation, inheritance, and polymorphism in Python.', 'https://tutora-endpoint-axbabq3egzii.private.compat.objectstorage.ap-osaka-1.oci.customer-oci.com/TUTORA/python-oop.gif'
FROM tutora_app.users WHERE email = 'jay@tutora.com';

-- 3. Add Root Comment by Haru Learner (Student) on Minas Video Post
INSERT INTO tutora_app.comments (post_id, user_id, body)
SELECT p.id, u.id, 'Wow, this makes calculus so simple! Thanks Mina!'
FROM tutora_app.posts p
CROSS JOIN tutora_app.users u
WHERE p.title = 'Visualizing Limits Dynamically' AND u.email = 'haru@tutora.com';

-- 4. Add Nested Reply by Ava Santos (Tutor) to Harus Comment
INSERT INTO tutora_app.comments (post_id, user_id, body, parent_id)
SELECT 
    p.id AS post_id, 
    u.id AS user_id, 
    'Mina is always the best at visualizing calculus concepts.' AS body, 
    c.id AS parent_id
FROM tutora_app.posts p
CROSS JOIN tutora_app.users u
CROSS JOIN tutora_app.comments c
WHERE p.title = 'Visualizing Limits Dynamically' 
  AND u.email = 'ava@tutora.com' 
  AND c.body = 'Wow, this makes calculus so simple! Thanks Mina!';

-- 5. Add Nested Reply by Haru Learner back to Ava
INSERT INTO tutora_app.comments (post_id, user_id, body, parent_id)
SELECT 
    p.id AS post_id, 
    u.id AS user_id, 
    'Agreed, I finally got my homework done!' AS body, 
    c.id AS parent_id
FROM tutora_app.posts p
CROSS JOIN tutora_app.users u
CROSS JOIN tutora_app.comments c
WHERE p.title = 'Visualizing Limits Dynamically' 
  AND u.email = 'haru@tutora.com' 
  AND c.body = 'Mina is always the best at visualizing calculus concepts.';

-- 6. Seed Comment Likes
-- Haru likes Avas comment
INSERT INTO tutora_app.comment_likes (comment_id, user_id)
SELECT c.id, u.id
FROM tutora_app.comments c
CROSS JOIN tutora_app.users u
WHERE c.body = 'Mina is always the best at visualizing calculus concepts.' AND u.email = 'haru@tutora.com';

-- Ava likes Harus comment (root)
INSERT INTO tutora_app.comment_likes (comment_id, user_id)
SELECT c.id, u.id
FROM tutora_app.comments c
CROSS JOIN tutora_app.users u
WHERE c.body = 'Wow, this makes calculus so simple! Thanks Mina!' AND u.email = 'ava@tutora.com';

-- 7. Seed Post Likes
-- Haru likes Minas video post
INSERT INTO tutora_app.post_likes (post_id, user_id)
SELECT p.id, u.id
FROM tutora_app.posts p
CROSS JOIN tutora_app.users u
WHERE p.title = 'Visualizing Limits Dynamically' AND u.email = 'haru@tutora.com';

-- 8. Seed Content Reports
-- Haru reports a mock comment (just for database seeding purposes)
INSERT INTO tutora_app.reports (reporter_id, target_type, target_id, reason)
SELECT u.id, 'comment', c.id, 'Self-reported testing comment.'
FROM tutora_app.users u
CROSS JOIN tutora_app.comments c
WHERE u.email = 'haru@tutora.com' AND c.body = 'Agreed, I finally got my homework done!';
