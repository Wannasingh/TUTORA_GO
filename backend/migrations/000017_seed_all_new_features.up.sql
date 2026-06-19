-- Migration 000017: Seed Data — focused on user ID 11 as primary test user
-- Real user IDs: 4(Haru), 5(Mina/Tutor1), 6(Jay/Tutor2), 7(Ava/Tutor3), 9(Google), 11(Wannasingh)

-- ============ FOLLOWS ============
INSERT INTO tutora_app.follows (follower_id, following_id) VALUES
    -- User 11 follows users 4, 5, 6, 7, 9
    (11, 4), (11, 5), (11, 6), (11, 7), (11, 9),
    -- Users follow user 11
    (4, 11), (5, 11), (6, 11), (7, 11), (9, 11),
    -- Other follow relationships
    (4, 5), (4, 6), (5, 4), (6, 7), (9, 5)
ON CONFLICT DO NOTHING;

-- ============ TUTOR REVIEWS ============
-- User 11 reviews tutors (tutor IDs 1,2,3 correspond to user IDs 5,6,7)
INSERT INTO tutora_app.tutor_reviews (reviewer_id, tutor_id, rating, body) VALUES
    (11, 1, 5.0, 'สอนคณิตศาสตร์ดีมากครับ อธิบายเข้าใจง่าย แนะนำเลย!'),
    (11, 2, 4.5, 'สอน Python สนุก ทำให้ชอบวิชานี้ขึ้นมา'),
    (11, 3, 4.0, 'สอน Biology ดี เนื้อหาเยอะมาก'),
    -- Others review tutors
    (4, 1, 4.5, 'อาจารย์สอนได้ดีมาก'),
    (6, 1, 5.0, 'สุดยอดครับ เข้าใจ Calculus ทันทีเลย'),
    (4, 2, 4.5, 'สอน AI สนุกมาก'),
    (9, 3, 3.5, 'โอเคครับ เนื้อหาค่อนข้างเร็ว')
ON CONFLICT (reviewer_id, tutor_id) DO NOTHING;

-- Update tutor ratings from reviews
UPDATE tutora_app.tutors SET rating = sub.avg_rating
FROM (
    SELECT tutor_id, ROUND(AVG(rating)::numeric, 2)::float8 as avg_rating
    FROM tutora_app.tutor_reviews
    GROUP BY tutor_id
) sub
WHERE tutora_app.tutors.id = sub.tutor_id;

-- ============ USER PROFILE DATA ============
UPDATE tutora_app.users SET 
    bio = 'นักศึกษาวิศวกรรมซอฟต์แวร์ ชอบ Coding, Math, Physics | BytesTutor Power User 🚀',
    school = 'จุฬาลงกรณ์มหาวิทยาลัย',
    phone = '098-765-4321'
WHERE id = 11;

UPDATE tutora_app.users SET bio = 'ชอบเรียนรู้สิ่งใหม่ๆ ทุกวัน', school = 'มหาวิทยาลัยธรรมศาสตร์' WHERE id = 4;
UPDATE tutora_app.users SET bio = 'ติวเตอร์คณิตศาสตร์และฟิสิกส์ ประสบการณ์ 5 ปี', school = 'มหาวิทยาลัยมหิดล' WHERE id = 5;
UPDATE tutora_app.users SET bio = 'Python Developer & AI Teacher', school = 'สถาบันเทคโนโลยีพระจอมเกล้าฯ ลาดกระบัง' WHERE id = 6;
UPDATE tutora_app.users SET bio = 'อาจารย์ชีววิทยาและเคมี สอนสนุก', school = 'มหาวิทยาลัยเกษตรศาสตร์' WHERE id = 7;

-- ============ REPOSTS ============
-- Need existing posts; check what we have
INSERT INTO tutora_app.reposts (user_id, post_id)
SELECT 11, p.id FROM tutora_app.posts p ORDER BY p.id LIMIT 3
ON CONFLICT DO NOTHING;

INSERT INTO tutora_app.reposts (user_id, post_id)
SELECT 5, p.id FROM tutora_app.posts p ORDER BY p.id LIMIT 1
ON CONFLICT DO NOTHING;

-- ============ USER NOTES (focused on user 11) ============
INSERT INTO tutora_app.user_notes (user_id, title, body, subject) VALUES
    (11, 'สรุป Calculus บทที่ 1 - Limits', 'Limits คือการหาค่าที่ฟังก์ชันเข้าใกล้เมื่อ x เข้าใกล้จุดใดจุดหนึ่ง

สูตรสำคัญ:
- lim(x→a) f(x) = L
- lim(x→a) [f(x) + g(x)] = lim f(x) + lim g(x)
- lim(x→0) sin(x)/x = 1

เทคนิค: ใช้ L''Hôpital''s Rule เมื่อ 0/0 หรือ ∞/∞', 'Mathematics'),
    (11, 'สรุป Newton Laws', 'กฎข้อที่ 1: วัตถุจะอยู่นิ่งหรือเคลื่อนที่ด้วยความเร็วคงที่ จนกว่าจะมีแรงภายนอกมากระทำ
กฎข้อที่ 2: F = ma (แรง = มวล × ความเร่ง)
กฎข้อที่ 3: แรงกิริยา = แรงปฏิกิริยา (ทิศตรงข้าม)', 'Physics'),
    (11, 'Python Data Structures Cheatsheet', '# List
my_list = [1, 2, 3]
my_list.append(4)

# Dictionary
my_dict = {"key": "value"}

# Set
my_set = {1, 2, 3}

# Tuple (immutable)
my_tuple = (1, 2, 3)', 'Programming'),
    (11, 'เทคนิคจำสูตรเคมี', 'PV = nRT → กฎของแก๊สอุดมคติ
pH = -log[H+] → การหา pH
E = hf → พลังงานโฟตอน', 'Chemistry'),
    (11, 'English Grammar Notes', 'Present Perfect: have/has + V3
- I have studied for 3 hours
- She has lived here since 2020

Past Simple: V2
- I went to school yesterday', 'English'),
    (6, 'Python Tips & Tricks', 'List comprehension: [x**2 for x in range(10)]', 'Programming');

-- ============ FLASHCARD DECKS (focused on user 11) ============
INSERT INTO tutora_app.flashcard_decks (user_id, title, subject) VALUES
    (11, 'Calculus สูตรที่ต้องรู้', 'Mathematics'),
    (11, 'Physics Constants', 'Physics'),
    (11, 'Python Built-in Functions', 'Programming'),
    (11, 'คำศัพท์ IELTS', 'English'),
    (6, 'Quantum Physics', 'Physics');

-- ============ FLASHCARDS ============
-- Get the deck IDs dynamically
INSERT INTO tutora_app.flashcards (deck_id, front_text, back_text, sort_order)
SELECT d.id, v.front, v.back, v.ord
FROM tutora_app.flashcard_decks d
CROSS JOIN (VALUES
    ('d/dx [x^n] = ?', 'n × x^(n-1)  (Power Rule)', 1),
    ('∫ x^n dx = ?', 'x^(n+1) / (n+1) + C  (n ≠ -1)', 2),
    ('d/dx [sin(x)] = ?', 'cos(x)', 3),
    ('d/dx [cos(x)] = ?', '-sin(x)', 4),
    ('d/dx [e^x] = ?', 'e^x', 5),
    ('d/dx [ln(x)] = ?', '1/x', 6),
    ('∫ 1/x dx = ?', 'ln|x| + C', 7)
) AS v(front, back, ord)
WHERE d.user_id = 11 AND d.title = 'Calculus สูตรที่ต้องรู้';

INSERT INTO tutora_app.flashcards (deck_id, front_text, back_text, sort_order)
SELECT d.id, v.front, v.back, v.ord
FROM tutora_app.flashcard_decks d
CROSS JOIN (VALUES
    ('Speed of light (c)', '3 × 10⁸ m/s', 1),
    ('Gravitational constant (G)', '6.674 × 10⁻¹¹ N⋅m²/kg²', 2),
    ('Planck constant (h)', '6.626 × 10⁻³⁴ J⋅s', 3),
    ('Boltzmann constant (k)', '1.381 × 10⁻²³ J/K', 4),
    ('Avogadro number (Nₐ)', '6.022 × 10²³ /mol', 5)
) AS v(front, back, ord)
WHERE d.user_id = 11 AND d.title = 'Physics Constants';

INSERT INTO tutora_app.flashcards (deck_id, front_text, back_text, sort_order)
SELECT d.id, v.front, v.back, v.ord
FROM tutora_app.flashcard_decks d
CROSS JOIN (VALUES
    ('len()', 'Returns the length/size of an object', 1),
    ('map(func, iterable)', 'Applies function to every item, returns map object', 2),
    ('enumerate(iterable)', 'Returns index and value pairs as tuples', 3),
    ('zip(iter1, iter2)', 'Combines iterables element-wise into tuples', 4),
    ('sorted(iterable)', 'Returns new sorted list from iterable', 5)
) AS v(front, back, ord)
WHERE d.user_id = 11 AND d.title = 'Python Built-in Functions';

INSERT INTO tutora_app.flashcards (deck_id, front_text, back_text, sort_order)
SELECT d.id, v.front, v.back, v.ord
FROM tutora_app.flashcard_decks d
CROSS JOIN (VALUES
    ('ubiquitous', 'พบได้ทุกที่ (adj.) — found everywhere', 1),
    ('pragmatic', 'เน้นปฏิบัติ (adj.) — practical, realistic', 2),
    ('meticulous', 'ละเอียดรอบคอบ (adj.) — very careful, precise', 3),
    ('ephemeral', 'ชั่วคราว (adj.) — lasting for a very short time', 4)
) AS v(front, back, ord)
WHERE d.user_id = 11 AND d.title = 'คำศัพท์ IELTS';

-- ============ USER COURSES (focused on user 11) ============
INSERT INTO tutora_app.user_courses (user_id, title, institution, status, started_at, completed_at) VALUES
    (11, 'Software Engineering', 'จุฬาลงกรณ์มหาวิทยาลัย', 'in_progress', '2024-06-01', NULL),
    (11, 'Data Structures & Algorithms', 'Coursera', 'completed', '2024-01-15', '2024-04-30'),
    (11, 'Machine Learning Specialization', 'Stanford / Coursera', 'in_progress', '2025-09-01', NULL),
    (11, 'Calculus I-II', 'จุฬาลงกรณ์มหาวิทยาลัย', 'completed', '2024-06-01', '2025-05-30'),
    (6, 'Quantum Mechanics I', 'สถาบันเทคโนโลยีพระจอมเกล้าฯ', 'completed', '2023-08-01', '2024-05-30');

-- ============ USER EXAMS (focused on user 11) ============
INSERT INTO tutora_app.user_exams (user_id, title, subject, score, max_score, exam_date) VALUES
    (11, 'Midterm Calculus I', 'Mathematics', '92', '100', '2025-10-15'),
    (11, 'Final Calculus I', 'Mathematics', '88', '100', '2025-12-20'),
    (11, 'Midterm Physics I', 'Physics', '85', '100', '2025-10-18'),
    (11, 'Final Physics I', 'Physics', '90', '100', '2025-12-22'),
    (11, 'Data Structures Midterm', 'Computer Science', '95', '100', '2025-10-20'),
    (11, 'IELTS Academic', 'English', '7.0', '9.0', '2026-01-15'),
    (6, 'Quantum Mechanics Final', 'Physics', '92', '100', '2024-12-15');

-- ============ CERTIFICATIONS (focused on user 11) ============
INSERT INTO tutora_app.certifications (user_id, title, issuer, date_earned) VALUES
    (11, 'Google IT Support Professional', 'Google / Coursera', '2025-06-01'),
    (11, 'AWS Certified Cloud Practitioner', 'Amazon Web Services', '2025-09-15'),
    (11, 'Meta Front-End Developer', 'Meta / Coursera', '2026-01-20'),
    (7, 'Python Institute PCAP', 'Python Institute', '2024-08-20'),
    (6, 'Stanford Machine Learning', 'Stanford / Coursera', '2025-03-10');

-- ============ BADGES ============
INSERT INTO tutora_app.badges (name, description, icon_url, criteria) VALUES
    ('first_post', 'โพสต์แรก', NULL, 'สร้างโพสต์แรกในชุมชน'),
    ('helpful_answer', 'คำตอบที่มีประโยชน์', NULL, 'ได้รับไลค์ 5 ครั้งในคอมเมนต์เดียว'),
    ('popular_post', 'โพสต์ยอดนิยม', NULL, 'โพสต์ได้รับไลค์ 10 ครั้ง'),
    ('top_tutor', 'ติวเตอร์ยอดเยี่ยม', NULL, 'ได้รับเรตติ้งเฉลี่ย 4.5 ขึ้นไป'),
    ('bookworm', 'นักอ่านตัวยง', NULL, 'ซื้อสื่อการเรียนรู้ 5 รายการ'),
    ('note_taker', 'จดโน้ตขยัน', NULL, 'สร้างโน้ตส่วนตัว 10 รายการ'),
    ('flashcard_master', 'แฟลชการ์ดมาสเตอร์', NULL, 'สร้างชุดแฟลชการ์ด 5 ชุด'),
    ('social_butterfly', 'สายสังคม', NULL, 'มีผู้ติดตาม 10 คน'),
    ('certified', 'มีใบรับรอง', NULL, 'เพิ่มใบรับรองแรก'),
    ('five_star', 'รีวิว 5 ดาว', NULL, 'ได้รับรีวิว 5 ดาวจากนักเรียน')
ON CONFLICT DO NOTHING;

-- ============ USER BADGES (focused on user 11) ============
INSERT INTO tutora_app.user_badges (user_id, badge_id)
SELECT 11, b.id FROM tutora_app.badges b WHERE b.name IN ('first_post', 'bookworm', 'note_taker', 'social_butterfly', 'certified')
ON CONFLICT DO NOTHING;

INSERT INTO tutora_app.user_badges (user_id, badge_id)
SELECT 5, b.id FROM tutora_app.badges b WHERE b.name IN ('first_post', 'top_tutor')
ON CONFLICT DO NOTHING;

INSERT INTO tutora_app.user_badges (user_id, badge_id)
SELECT 6, b.id FROM tutora_app.badges b WHERE b.name = 'first_post'
ON CONFLICT DO NOTHING;

-- ============ CONVERSATIONS & MESSAGES ============
-- Direct conversation: user 11 <-> user 5 (Mina/Tutor)
INSERT INTO tutora_app.conversations (type, title) VALUES ('direct', NULL);
INSERT INTO tutora_app.conversation_members (conversation_id, user_id)
SELECT currval('tutora_app.conversations_id_seq'), unnest(ARRAY[11, 5]);
INSERT INTO tutora_app.messages (conversation_id, sender_id, body, message_type)
SELECT currval('tutora_app.conversations_id_seq'), v.sender, v.body, 'text'
FROM (VALUES
    (11, 'สวัสดีครับอาจารย์ อยากสอบถามเรื่อง Calculus หน่อยครับ'),
    (5, 'ได้เลยค่ะ สงสัยเรื่องอะไร?'),
    (11, 'เรื่อง Integration by Parts ครับ ยังไม่ค่อยเข้าใจ'),
    (5, 'ใช้สูตร ∫u dv = uv - ∫v du ค่ะ ลองเลือก u ให้เป็นฟังก์ชันที่ differentiate แล้วง่ายขึ้น'),
    (11, 'เข้าใจแล้วครับ ขอบคุณมากครับ! 🙏')
) AS v(sender, body);

-- Direct conversation: user 11 <-> user 6 (Jay)
INSERT INTO tutora_app.conversations (type, title) VALUES ('direct', NULL);
INSERT INTO tutora_app.conversation_members (conversation_id, user_id)
SELECT currval('tutora_app.conversations_id_seq'), unnest(ARRAY[11, 6]);
INSERT INTO tutora_app.messages (conversation_id, sender_id, body, message_type)
SELECT currval('tutora_app.conversations_id_seq'), v.sender, v.body, 'text'
FROM (VALUES
    (6, 'เห็นโพสต์เรื่อง Python ลองอ่านหนังสือเล่มนี้ดูนะ'),
    (11, 'ขอบคุณครับ! จะลองหาอ่านดู'),
    (11, 'อ่านแล้วเข้าใจขึ้นเยอะเลยครับ')
) AS v(sender, body);

-- Direct conversation: user 11 <-> user 9
INSERT INTO tutora_app.conversations (type, title) VALUES ('direct', NULL);
INSERT INTO tutora_app.conversation_members (conversation_id, user_id)
SELECT currval('tutora_app.conversations_id_seq'), unnest(ARRAY[11, 9]);
INSERT INTO tutora_app.messages (conversation_id, sender_id, body, message_type)
SELECT currval('tutora_app.conversations_id_seq'), v.sender, v.body, 'text'
FROM (VALUES
    (11, 'เปิดกลุ่มติว Calculus ด้วยกันไหมครับ?'),
    (9, 'ดีเลยครับ สนใจมากเลย!')
) AS v(sender, body);

-- Group conversation: study group
INSERT INTO tutora_app.conversations (type, title) VALUES ('group', 'กลุ่มติว Calculus & Physics 📚');
INSERT INTO tutora_app.conversation_members (conversation_id, user_id)
SELECT currval('tutora_app.conversations_id_seq'), unnest(ARRAY[11, 5, 9, 6]);
INSERT INTO tutora_app.messages (conversation_id, sender_id, body, message_type)
SELECT currval('tutora_app.conversations_id_seq'), v.sender, v.body, 'text'
FROM (VALUES
    (11, 'สวัสดีทุกคน! กลุ่มนี้ไว้ถาม-ตอบเรื่อง Calculus & Physics นะครับ'),
    (5, 'ยินดีช่วยตอบทุกคำถามค่ะ 😊'),
    (9, 'ดีเลยครับ มีหลายเรื่องอยากถาม'),
    (6, 'Python ถามได้เลยนะ!'),
    (11, 'ใครพอจะอธิบาย Taylor Series ได้บ้างครับ?')
) AS v(sender, body);

-- ============ NOTIFICATIONS (focused on user 11) ============
INSERT INTO tutora_app.notifications (user_id, type, title, body, data_json) VALUES
    (11, 'system', 'ยินดีต้อนรับสู่ BytesTutor! 🎉', 'เริ่มต้นการเรียนรู้ของคุณได้เลย', '{"action": "welcome"}'),
    (11, 'system', 'มีคนติดตามคุณ', 'Haru Learner เริ่มติดตามคุณ', '{"follower_id": 4}'),
    (11, 'system', 'มีคนติดตามคุณ', 'Mina Park เริ่มติดตามคุณ', '{"follower_id": 5}'),
    (11, 'system', 'มีคนถูกใจโพสต์ของคุณ', 'มีคนถูกใจโพสต์ของคุณ 3 ครั้ง', '{"post_id": 1}'),
    (11, 'message', 'ข้อความใหม่', 'คุณมีข้อความใหม่จาก Mina Park', '{"conversation_id": 1}'),
    (11, 'remind', 'ถึงเวลาทบทวน! 📖', 'ถึงเวลาทบทวน Flashcard Calculus สูตรที่ต้องรู้ แล้ว!', '{"action": "review_flashcard"}'),
    (11, 'remind', 'อย่าลืมทำ Practice! 📝', 'คุณยังไม่ได้ทบทวน Physics Constants มา 3 วันแล้ว', '{"action": "review_flashcard"}'),
    (11, 'promotion', 'สื่อการเรียนใหม่! 🛒', 'มีสื่อการเรียน Calculus ใหม่ในร้านค้า ลดราคา 20%', '{"category": "notes", "subject": "Mathematics"}'),
    (11, 'promotion', 'โปรโมชั่นพิเศษ! 💰', 'เติมเหรียญวันนี้ รับโบนัสเพิ่ม 10%', '{"action": "top_up"}'),
    (11, 'critical', 'อัพเดทข้อกำหนดการใช้งาน', 'กรุณาอ่านข้อกำหนดการใช้งานใหม่ของเรา', '{"action": "tos_update"}'),
    (5, 'message', 'ข้อความใหม่', 'คุณมีข้อความใหม่', '{"conversation_id": 1}'),
    (9, 'remind', 'ถึงเวลาทบทวน!', 'ถึงเวลาทบทวน Chemistry แล้ว!', '{"action": "review_flashcard"}');
