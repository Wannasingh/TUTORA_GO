-- Rollback seed data for new features
DELETE FROM tutora_app.notifications;
DELETE FROM tutora_app.messages;
DELETE FROM tutora_app.conversation_members;
DELETE FROM tutora_app.conversations;
DELETE FROM tutora_app.user_badges;
DELETE FROM tutora_app.badges;
DELETE FROM tutora_app.certifications;
DELETE FROM tutora_app.user_exams;
DELETE FROM tutora_app.user_courses;
DELETE FROM tutora_app.flashcards;
DELETE FROM tutora_app.flashcard_decks;
DELETE FROM tutora_app.user_notes;
DELETE FROM tutora_app.reposts;
DELETE FROM tutora_app.tutor_reviews;
DELETE FROM tutora_app.follows;
UPDATE tutora_app.users SET bio = NULL, cover_url = NULL, phone = NULL, school = NULL, birthdate = NULL;
