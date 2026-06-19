-- Author: Antigravity
-- Date: 2026-06-19
-- Task/Jira ID: BT-6
-- Purpose: Revert roles TEXT[] array column back to single role VARCHAR column

ALTER TABLE tutora_app.users ADD COLUMN role VARCHAR(50) NOT NULL DEFAULT 'student';

UPDATE tutora_app.users SET role = COALESCE(roles[1], 'student');

ALTER TABLE tutora_app.users DROP COLUMN roles;
