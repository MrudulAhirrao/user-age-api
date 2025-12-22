-- Add activation fields to users table
ALTER TABLE users 
ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN activation_token VARCHAR(255),
ADD COLUMN activation_sent_at TIMESTAMPTZ;

-- Optional: Create an index on token for faster lookup
CREATE INDEX idx_users_activation_token ON users(activation_token);