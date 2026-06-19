-- Migration 000015: Messaging System

-- Conversations (direct or group)
CREATE TABLE tutora_app.conversations (
    id         SERIAL PRIMARY KEY,
    type       VARCHAR(20) NOT NULL DEFAULT 'direct',
    title      VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Conversation members
CREATE TABLE tutora_app.conversation_members (
    conversation_id INTEGER NOT NULL REFERENCES tutora_app.conversations(id) ON DELETE CASCADE,
    user_id         INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    last_read_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    joined_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (conversation_id, user_id)
);
CREATE INDEX idx_conv_members_user ON tutora_app.conversation_members(user_id);

-- Messages
CREATE TABLE tutora_app.messages (
    id              SERIAL PRIMARY KEY,
    conversation_id INTEGER NOT NULL REFERENCES tutora_app.conversations(id) ON DELETE CASCADE,
    sender_id       INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    body            TEXT,
    image_url       VARCHAR(512),
    message_type    VARCHAR(50) NOT NULL DEFAULT 'text',
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_messages_conv ON tutora_app.messages(conversation_id, created_at DESC);
CREATE INDEX idx_messages_sender ON tutora_app.messages(sender_id);
