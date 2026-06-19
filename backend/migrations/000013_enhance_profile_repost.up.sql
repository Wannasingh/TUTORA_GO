-- Migration 000013: Enhanced Profile & Repost/Quote

-- Add profile fields to users
ALTER TABLE tutora_app.users ADD COLUMN bio TEXT;
ALTER TABLE tutora_app.users ADD COLUMN cover_url VARCHAR(512);
ALTER TABLE tutora_app.users ADD COLUMN phone VARCHAR(50);
ALTER TABLE tutora_app.users ADD COLUMN school VARCHAR(255);
ALTER TABLE tutora_app.users ADD COLUMN birthdate DATE;

-- Add quote post reference to posts
ALTER TABLE tutora_app.posts ADD COLUMN original_post_id INTEGER
    REFERENCES tutora_app.posts(id) ON DELETE SET NULL;
CREATE INDEX idx_posts_original ON tutora_app.posts(original_post_id)
    WHERE original_post_id IS NOT NULL;

-- Simple repost (share toggle)
CREATE TABLE tutora_app.reposts (
    user_id    INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    post_id    INTEGER NOT NULL REFERENCES tutora_app.posts(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, post_id)
);
CREATE INDEX idx_reposts_post ON tutora_app.reposts(post_id);
