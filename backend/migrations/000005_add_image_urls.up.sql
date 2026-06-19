-- Author: Haru
-- Date: 2026-06-19
-- Task/Jira ID: BT-5
-- Purpose: Add avatar_url to users and image_url to posts and comments

ALTER TABLE tutora_app.users ADD COLUMN IF NOT EXISTS avatar_url VARCHAR(512) NULL;
ALTER TABLE tutora_app.posts ADD COLUMN IF NOT EXISTS image_url VARCHAR(512) NULL;
ALTER TABLE tutora_app.comments ADD COLUMN IF NOT EXISTS image_url VARCHAR(512) NULL;
