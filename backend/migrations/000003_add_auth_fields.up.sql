-- Author: Haru
-- Date: 2026-06-18
-- Task/Jira ID: BT-3
-- Purpose: Add auth credentials and oauth identity tracking to users table

ALTER TABLE tutora_app.users 
ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255) NULL,
ADD COLUMN IF NOT EXISTS google_id VARCHAR(255) UNIQUE NULL,
ADD COLUMN IF NOT EXISTS apple_id VARCHAR(255) UNIQUE NULL;
