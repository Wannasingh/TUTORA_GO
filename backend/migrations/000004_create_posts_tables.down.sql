-- Author: Haru
-- Date: 2026-06-18
-- Task/Jira ID: BT-4
-- Purpose: Drop posts, comments, likes, and saves tables in tutora_app schema

DROP TABLE IF EXISTS tutora_app.post_saves;
DROP TABLE IF EXISTS tutora_app.post_likes;
DROP TABLE IF EXISTS tutora_app.comments;
DROP TABLE IF EXISTS tutora_app.posts;
