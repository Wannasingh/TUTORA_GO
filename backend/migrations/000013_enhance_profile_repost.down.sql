DROP TABLE IF EXISTS tutora_app.reposts;
DROP INDEX IF EXISTS tutora_app.idx_posts_original;
ALTER TABLE tutora_app.posts DROP COLUMN IF EXISTS original_post_id;
ALTER TABLE tutora_app.users DROP COLUMN IF EXISTS birthdate;
ALTER TABLE tutora_app.users DROP COLUMN IF EXISTS school;
ALTER TABLE tutora_app.users DROP COLUMN IF EXISTS phone;
ALTER TABLE tutora_app.users DROP COLUMN IF EXISTS cover_url;
ALTER TABLE tutora_app.users DROP COLUMN IF EXISTS bio;
