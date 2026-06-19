-- Migration 000014: Study Tools

-- Personal notes
CREATE TABLE tutora_app.user_notes (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    title      VARCHAR(255) NOT NULL,
    body       TEXT,
    subject    VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_user_notes_user ON tutora_app.user_notes(user_id);

-- Flashcard decks
CREATE TABLE tutora_app.flashcard_decks (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    title      VARCHAR(255) NOT NULL,
    subject    VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_flashcard_decks_user ON tutora_app.flashcard_decks(user_id);

-- Flashcards
CREATE TABLE tutora_app.flashcards (
    id         SERIAL PRIMARY KEY,
    deck_id    INTEGER NOT NULL REFERENCES tutora_app.flashcard_decks(id) ON DELETE CASCADE,
    front_text TEXT NOT NULL,
    back_text  TEXT NOT NULL,
    image_url  VARCHAR(512),
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_flashcards_deck ON tutora_app.flashcards(deck_id);

-- Course tracking
CREATE TABLE tutora_app.user_courses (
    id           SERIAL PRIMARY KEY,
    user_id      INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    title        VARCHAR(255) NOT NULL,
    institution  VARCHAR(255),
    status       VARCHAR(50) NOT NULL DEFAULT 'in_progress',
    started_at   DATE,
    completed_at DATE,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_user_courses_user ON tutora_app.user_courses(user_id);

-- Exam records
CREATE TABLE tutora_app.user_exams (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    title      VARCHAR(255) NOT NULL,
    subject    VARCHAR(255),
    score      VARCHAR(50),
    max_score  VARCHAR(50),
    exam_date  DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_user_exams_user ON tutora_app.user_exams(user_id);

-- Certifications
CREATE TABLE tutora_app.certifications (
    id             SERIAL PRIMARY KEY,
    user_id        INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    title          VARCHAR(255) NOT NULL,
    issuer         VARCHAR(255) NOT NULL,
    date_earned    DATE,
    expiry_date    DATE,
    image_url      VARCHAR(512),
    credential_url VARCHAR(512),
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_certifications_user ON tutora_app.certifications(user_id);

-- Badge definitions
CREATE TABLE tutora_app.badges (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    icon_url    VARCHAR(512),
    criteria    TEXT NOT NULL
);

-- User badge awards
CREATE TABLE tutora_app.user_badges (
    user_id   INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    badge_id  INTEGER NOT NULL REFERENCES tutora_app.badges(id) ON DELETE CASCADE,
    earned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, badge_id)
);
