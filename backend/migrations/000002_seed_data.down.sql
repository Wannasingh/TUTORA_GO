-- Author: Haru
-- Date: 2026-06-18
-- Task/Jira ID: BT-2
-- Purpose: Tear down seed data from tutora_app schema

TRUNCATE TABLE tutora_app.tutors RESTART IDENTITY CASCADE;
TRUNCATE TABLE tutora_app.users RESTART IDENTITY CASCADE;
