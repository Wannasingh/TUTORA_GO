-- Author: Antigravity
-- Date: 2026-06-19
-- Task/Jira ID: BT-10
-- Purpose: Create comment replies, comment likes, reports, and post video fields

-- Upgrades to Comments for Nested Replies
ALTER TABLE tutora_app.comments ADD COLUMN IF NOT EXISTS parent_id INTEGER REFERENCES tutora_app.comments(id) ON DELETE CASCADE;
ALTER TABLE tutora_app.comments ADD COLUMN IF NOT EXISTS status VARCHAR(50) NOT NULL DEFAULT 'active';

-- Upgrades to Posts for Video Support
ALTER TABLE tutora_app.posts ADD COLUMN IF NOT EXISTS video_url VARCHAR(512) NULL;

-- Comment Likes Table
CREATE TABLE IF NOT EXISTS tutora_app.comment_likes (
    comment_id INTEGER NOT NULL REFERENCES tutora_app.comments(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (comment_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_comment_likes_comment_id ON tutora_app.comment_likes(comment_id);

-- Content Reports Table
CREATE TABLE IF NOT EXISTS tutora_app.reports (
    id SERIAL PRIMARY KEY,
    reporter_id INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    target_type VARCHAR(50) NOT NULL, -- 'post' | 'comment'
    target_id INTEGER NOT NULL,
    reason TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_reports_reporter_id ON tutora_app.reports(reporter_id);
