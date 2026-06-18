-- Author: Haru
-- Date: 2026-06-18
-- Task/Jira ID: BT-1
-- Purpose: Initialize the tutora_app schema and basic users/tutors tables

-- Create Schema
CREATE SCHEMA IF NOT EXISTS tutora_app;

-- Create Users Table
CREATE TABLE IF NOT EXISTS tutora_app.users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'student',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create Tutors Table
CREATE TABLE IF NOT EXISTS tutora_app.tutors (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    subject VARCHAR(255) NOT NULL,
    bio TEXT,
    price_per_hour NUMERIC(10, 2) NOT NULL,
    rating NUMERIC(3, 2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
