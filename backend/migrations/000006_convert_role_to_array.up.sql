-- Author: Antigravity
-- Date: 2026-06-19
-- Task/Jira ID: BT-6
-- Purpose: Convert user role column from VARCHAR to roles TEXT[] array

ALTER TABLE tutora_app.users ADD COLUMN roles TEXT[] NOT NULL DEFAULT '{student}';

UPDATE tutora_app.users SET roles = ARRAY[role]::TEXT[];

ALTER TABLE tutora_app.users DROP COLUMN role;
