-- Author: Antigravity
-- Date: 2026-06-19
-- Task/Jira ID: BT-7
-- Purpose: Create store marketplace, coin ledger, and payout tables for Apple IAP compliance

-- Coin Transactions (Double-Entry Ledger)
CREATE TABLE IF NOT EXISTS tutora_app.coin_transactions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    amount INTEGER NOT NULL, -- positive for credit, negative for debit
    transaction_type VARCHAR(50) NOT NULL, -- 'iap_purchase', 'marketplace_buy', 'marketplace_sale', 'withdrawal'
    reference_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_coin_transactions_user_id ON tutora_app.coin_transactions(user_id);

-- Store Items Table
CREATE TABLE IF NOT EXISTS tutora_app.store_items (
    id SERIAL PRIMARY KEY,
    seller_id INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    category VARCHAR(100) NOT NULL, -- 'notes', 'practice_exams', 'video_course'
    subject VARCHAR(255) NOT NULL,
    price_in_coins INTEGER NOT NULL,
    file_url VARCHAR(512) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active', -- 'active', 'suspended', 'deleted'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_store_items_seller_id ON tutora_app.store_items(seller_id);
CREATE INDEX IF NOT EXISTS idx_store_items_subject ON tutora_app.store_items(subject);

-- Item Purchases Table
CREATE TABLE IF NOT EXISTS tutora_app.item_purchases (
    id SERIAL PRIMARY KEY,
    buyer_id INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    item_id INTEGER NOT NULL REFERENCES tutora_app.store_items(id) ON DELETE CASCADE,
    coins_spent INTEGER NOT NULL,
    coins_platform_fee INTEGER NOT NULL,
    coins_seller_amount INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_item_purchases_buyer_id ON tutora_app.item_purchases(buyer_id);
CREATE INDEX IF NOT EXISTS idx_item_purchases_item_id ON tutora_app.item_purchases(item_id);

-- Payouts Table
CREATE TABLE IF NOT EXISTS tutora_app.payouts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES tutora_app.users(id) ON DELETE CASCADE,
    coins_debited INTEGER NOT NULL,
    cash_amount_thb NUMERIC(10, 2) NOT NULL,
    bank_account_details TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'requested', -- 'requested', 'processed', 'rejected'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_payouts_user_id ON tutora_app.payouts(user_id);
