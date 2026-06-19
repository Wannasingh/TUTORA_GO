-- Migration 000012: Social Graph (Follows + Tutor Reviews)

-- Follow system
CREATE TABLE tutora_app.follows (
    follower_id  INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    following_id INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (follower_id, following_id),
    CHECK (follower_id != following_id)
);
CREATE INDEX idx_follows_follower ON tutora_app.follows(follower_id);
CREATE INDEX idx_follows_following ON tutora_app.follows(following_id);

-- Tutor review system
CREATE TABLE tutora_app.tutor_reviews (
    id          SERIAL PRIMARY KEY,
    reviewer_id INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    tutor_id    INTEGER NOT NULL REFERENCES tutora_app.tutors(id) ON DELETE CASCADE,
    rating      NUMERIC(2,1) NOT NULL CHECK (rating >= 1 AND rating <= 5),
    body        TEXT,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (reviewer_id, tutor_id)
);
CREATE INDEX idx_tutor_reviews_tutor ON tutora_app.tutor_reviews(tutor_id);
