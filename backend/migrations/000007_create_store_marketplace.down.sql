-- Author: Antigravity
-- Date: 2026-06-19
-- Task/Jira ID: BT-7
-- Purpose: Drop store marketplace, coin ledger, and payout tables

DROP TABLE IF EXISTS tutora_app.payouts;
DROP TABLE IF EXISTS tutora_app.item_purchases;
DROP TABLE IF EXISTS tutora_app.store_items;
DROP TABLE IF EXISTS tutora_app.coin_transactions;
