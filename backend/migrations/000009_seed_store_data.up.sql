-- Author: Antigravity
-- Date: 2026-06-19
-- Task/Jira ID: BT-9
-- Purpose: Seed store items, coin transactions, purchases, and interactive Q&A pins

-- Clean up existing store seeds
TRUNCATE TABLE tutora_app.coin_transactions RESTART IDENTITY CASCADE;
TRUNCATE TABLE tutora_app.payouts RESTART IDENTITY CASCADE;
TRUNCATE TABLE tutora_app.item_purchases RESTART IDENTITY CASCADE;
TRUNCATE TABLE tutora_app.item_qa_replies RESTART IDENTITY CASCADE;
TRUNCATE TABLE tutora_app.item_qa_pins RESTART IDENTITY CASCADE;
TRUNCATE TABLE tutora_app.store_items RESTART IDENTITY CASCADE;

-- Insert Store Items
-- Mina Park (Calculus, Physics)
INSERT INTO tutora_app.store_items (seller_id, title, description, category, subject, price_in_coins, file_url, status)
SELECT id, 'Ultimate Calculus Cheat Sheet', '10 pages of essential derivative and integral shortcuts with visual proofs.', 'notes', 'Calculus', 120, 'https://tutora-endpoint-axbabq3egzii.private.compat.objectstorage.ap-osaka-1.oci.customer-oci.com/TUTORA/calculus-cheat-sheet.pdf', 'active'
FROM tutora_app.users WHERE email = 'mina@tutora.com';

INSERT INTO tutora_app.store_items (seller_id, title, description, category, subject, price_in_coins, file_url, status)
SELECT id, 'Classical Mechanics Formula Pack', 'Equations sheet covering Kinematics, Dynamics, and Rotational motion.', 'notes', 'Physics', 80, 'https://tutora-endpoint-axbabq3egzii.private.compat.objectstorage.ap-osaka-1.oci.customer-oci.com/TUTORA/mechanics-formulas.pdf', 'active'
FROM tutora_app.users WHERE email = 'mina@tutora.com';

-- Jay Chen (Python, AI)
INSERT INTO tutora_app.store_items (seller_id, title, description, category, subject, price_in_coins, file_url, status)
SELECT id, 'Python OOP Practical Exercises', 'Hands-on practice workbook for Classes, Inheritance, and Interfaces.', 'practice_exams', 'Python', 90, 'https://tutora-endpoint-axbabq3egzii.private.compat.objectstorage.ap-osaka-1.oci.customer-oci.com/TUTORA/python-oop-exercises.pdf', 'active'
FROM tutora_app.users WHERE email = 'jay@tutora.com';

-- Ava Santos (Biology, Chemistry)
INSERT INTO tutora_app.store_items (seller_id, title, description, category, subject, price_in_coins, file_url, status)
SELECT id, 'Organic Chemistry Reaction Roadmap', 'Full dynamic mind map summarizing all fundamental organic synthesis pathways.', 'notes', 'Chemistry', 150, 'https://tutora-endpoint-axbabq3egzii.private.compat.objectstorage.ap-osaka-1.oci.customer-oci.com/TUTORA/organic-chem-roadmap.pdf', 'active'
FROM tutora_app.users WHERE email = 'ava@tutora.com';

-- Seed Coin Balance for Haru Learner (1000 Coins via Apple IAP)
INSERT INTO tutora_app.coin_transactions (user_id, amount, transaction_type, reference_id)
SELECT id, 1000, 'iap_purchase', 'apple_ref_mock_seed_1000'
FROM tutora_app.users WHERE email = 'haru@tutora.com';

-- Seed purchases: Haru Learner buys 'Ultimate Calculus Cheat Sheet' (120 coins)
-- We insert the purchase and transactions manually for reproducibility
INSERT INTO tutora_app.item_purchases (buyer_id, item_id, coins_spent, coins_platform_fee, coins_seller_amount)
SELECT 
    u.id AS buyer_id, 
    s.id AS item_id, 
    120 AS coins_spent, 
    18 AS coins_platform_fee, 
    102 AS coins_seller_amount
FROM tutora_app.users u
CROSS JOIN tutora_app.store_items s
WHERE u.email = 'haru@tutora.com' AND s.title = 'Ultimate Calculus Cheat Sheet';

-- Deduct from Haru
INSERT INTO tutora_app.coin_transactions (user_id, amount, transaction_type, reference_id)
SELECT u.id, -120, 'marketplace_buy', 'purchase_' || ip.id
FROM tutora_app.users u
CROSS JOIN tutora_app.item_purchases ip
JOIN tutora_app.store_items s ON ip.item_id = s.id
WHERE u.email = 'haru@tutora.com' AND s.title = 'Ultimate Calculus Cheat Sheet';

-- Credit to Mina Park
INSERT INTO tutora_app.coin_transactions (user_id, amount, transaction_type, reference_id)
SELECT u.id, 102, 'marketplace_sale', 'purchase_' || ip.id
FROM tutora_app.users u
CROSS JOIN tutora_app.item_purchases ip
JOIN tutora_app.store_items s ON ip.item_id = s.id
WHERE u.email = 'mina@tutora.com' AND s.title = 'Ultimate Calculus Cheat Sheet';

-- Seed Q&A pin on the Calculus Sheet by Haru
INSERT INTO tutora_app.item_qa_pins (item_id, user_id, page_number, coordinate_x, coordinate_y, question_text)
SELECT s.id, u.id, 2, 45.5, 72.8, 'How did we derive the chain rule shortcut shown on line 12?'
FROM tutora_app.store_items s
CROSS JOIN tutora_app.users u
WHERE s.title = 'Ultimate Calculus Cheat Sheet' AND u.email = 'haru@tutora.com';

-- Seed Q&A reply by Mina Park
INSERT INTO tutora_app.item_qa_replies (pin_id, user_id, reply_text)
SELECT p.id, u.id, 'It uses the limit definition of derivative composite functions. I added a video walkthrough link on the next page.'
FROM tutora_app.item_qa_pins p
CROSS JOIN tutora_app.users u
WHERE p.question_text = 'How did we derive the chain rule shortcut shown on line 12?' AND u.email = 'mina@tutora.com';
