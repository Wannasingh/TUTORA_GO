-- Author: Antigravity
-- Date: 2026-06-19
-- Task/Jira ID: BT-8
-- Purpose: Create interactive notebook Q&A pins and replies tables

-- Item Q&A Pins
CREATE TABLE IF NOT EXISTS tutora_app.item_qa_pins (
    id SERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL REFERENCES tutora_app.store_items(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    page_number INTEGER NOT NULL,
    coordinate_x NUMERIC(6, 3) NOT NULL,
    coordinate_y NUMERIC(6, 3) NOT NULL,
    question_text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_item_qa_pins_item_id ON tutora_app.item_qa_pins(item_id);

-- Item Q&A Replies
CREATE TABLE IF NOT EXISTS tutora_app.item_qa_replies (
    id SERIAL PRIMARY KEY,
    pin_id INTEGER NOT NULL REFERENCES tutora_app.item_qa_pins(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    reply_text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_item_qa_replies_pin_id ON tutora_app.item_qa_replies(pin_id);
