-- Author: Antigravity
-- Date: 2026-06-19
-- Task/Jira ID: BT-10
-- Purpose: Revert reports, comment likes, video url, and parent comments

DROP TABLE IF EXISTS tutora_app.reports;
DROP TABLE IF EXISTS tutora_app.comment_likes;
ALTER TABLE tutora_app.posts DROP COLUMN IF EXISTS video_url;
ALTER TABLE tutora_app.comments DROP COLUMN IF EXISTS status;
ALTER TABLE tutora_app.comments DROP COLUMN IF EXISTS parent_id;
