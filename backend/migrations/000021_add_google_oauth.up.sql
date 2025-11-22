-- Add Google OAuth support to users table
ALTER TABLE users
ADD COLUMN IF NOT EXISTS google_id VARCHAR(255) UNIQUE,
ADD COLUMN IF NOT EXISTS auth_provider VARCHAR(50) DEFAULT 'email';

-- Index for faster Google ID lookups
CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id) WHERE google_id IS NOT NULL;

-- Comment for documentation
COMMENT ON COLUMN users.google_id IS 'Google OAuth unique identifier (sub claim from Google token)';
COMMENT ON COLUMN users.auth_provider IS 'Primary authentication method: email, google';
