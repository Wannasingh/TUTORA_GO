-- Author: Antigravity
-- Date: 2026-06-19
-- Task/Jira ID: BT-9
-- Purpose: Clean up store seeded records

TRUNCATE TABLE tutora_app.coin_transactions RESTART IDENTITY CASCADE;
TRUNCATE TABLE tutora_app.payouts RESTART IDENTITY CASCADE;
TRUNCATE TABLE tutora_app.item_purchases RESTART IDENTITY CASCADE;
TRUNCATE TABLE tutora_app.item_qa_replies RESTART IDENTITY CASCADE;
TRUNCATE TABLE tutora_app.item_qa_pins RESTART IDENTITY CASCADE;
TRUNCATE TABLE tutora_app.store_items RESTART IDENTITY CASCADE;
