-- Remove Google OAuth support from users table
DROP INDEX IF EXISTS idx_users_google_id;
ALTER TABLE users
DROP COLUMN IF EXISTS google_id,
DROP COLUMN IF EXISTS auth_provider;
