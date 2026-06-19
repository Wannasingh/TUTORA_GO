-- Author: Haru
-- Date: 2026-06-19
-- Task/Jira ID: BT-5
-- Purpose: Remove avatar_url from users and image_url from posts and comments

ALTER TABLE tutora_app.users DROP COLUMN IF EXISTS avatar_url;
ALTER TABLE tutora_app.posts DROP COLUMN IF EXISTS image_url;
ALTER TABLE tutora_app.comments DROP COLUMN IF EXISTS image_url;
