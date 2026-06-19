ALTER TABLE tutora_app.users ADD COLUMN roles TEXT[] NOT NULL DEFAULT '{student}';

UPDATE tutora_app.users SET roles = ARRAY[role]::TEXT[];

ALTER TABLE tutora_app.users DROP COLUMN role;
