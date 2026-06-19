-- Author: Haru
-- Date: 2026-06-18
-- Task/Jira ID: BT-4
-- Purpose: Create posts, comments, likes, and saves tables in tutora_app schema

-- Posts Table
CREATE TABLE IF NOT EXISTS tutora_app.posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    subject VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_posts_user_id ON tutora_app.posts(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_subject ON tutora_app.posts(subject);

-- Comments Table
CREATE TABLE IF NOT EXISTS tutora_app.comments (
    id SERIAL PRIMARY KEY,
    post_id INTEGER NOT NULL REFERENCES tutora_app.posts(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_comments_post_id ON tutora_app.comments(post_id);

-- Likes Table
CREATE TABLE IF NOT EXISTS tutora_app.post_likes (
    post_id INTEGER NOT NULL REFERENCES tutora_app.posts(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (post_id, user_id)
);

-- Saves Table
CREATE TABLE IF NOT EXISTS tutora_app.post_saves (
    post_id INTEGER NOT NULL REFERENCES tutora_app.posts(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (post_id, user_id)
);
