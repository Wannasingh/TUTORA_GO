-- Author: Haru
-- Date: 2026-06-18
-- Task/Jira ID: BT-3
-- Purpose: Rollback auth credentials and oauth identity tracking columns from users table

ALTER TABLE tutora_app.users
DROP COLUMN IF EXISTS password_hash,
DROP COLUMN IF EXISTS google_id,
DROP COLUMN IF EXISTS apple_id;
