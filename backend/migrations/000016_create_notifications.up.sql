-- Migration 000016: Notifications

CREATE TABLE tutora_app.notifications (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    type       VARCHAR(50) NOT NULL,
    title      VARCHAR(255) NOT NULL,
    body       TEXT,
    data_json  JSONB,
    is_read    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_notifications_user ON tutora_app.notifications(user_id, created_at DESC);
CREATE INDEX idx_notifications_unread ON tutora_app.notifications(user_id) WHERE is_read = FALSE;
