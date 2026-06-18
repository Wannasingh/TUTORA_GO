-- Author: Haru
-- Date: 2026-06-18
-- Task/Jira ID: BT-1
-- Purpose: Tear down tutora_app schema and tables

DROP TABLE IF EXISTS tutora_app.tutors;
DROP TABLE IF EXISTS tutora_app.users;
DROP SCHEMA IF EXISTS tutora_app CASCADE;
