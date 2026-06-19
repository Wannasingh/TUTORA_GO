ALTER TABLE tutora_app.users ADD COLUMN role VARCHAR(50) NOT NULL DEFAULT 'student';

UPDATE tutora_app.users SET role = COALESCE(roles[1], 'student');

ALTER TABLE tutora_app.users DROP COLUMN roles;
