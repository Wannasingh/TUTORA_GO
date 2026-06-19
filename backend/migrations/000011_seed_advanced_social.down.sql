-- Author: Antigravity
-- Date: 2026-06-19
-- Task/Jira ID: BT-11
-- Purpose: Clean up advanced social seeded records

TRUNCATE TABLE tutora_app.reports RESTART IDENTITY CASCADE;
TRUNCATE TABLE tutora_app.comment_likes RESTART IDENTITY CASCADE;
DELETE FROM tutora_app.comments WHERE post_id IN (SELECT id FROM tutora_app.posts WHERE subject IN ('Calculus', 'Python', 'Moderation Test'));
DELETE FROM tutora_app.post_likes WHERE post_id IN (SELECT id FROM tutora_app.posts WHERE subject IN ('Calculus', 'Python', 'Moderation Test'));
DELETE FROM tutora_app.posts WHERE subject IN ('Calculus', 'Python', 'Moderation Test');
