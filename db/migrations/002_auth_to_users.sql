Alter TABLE users
ADD COLUMN email VARCHAR(255) unique NOT NULL,
ADD COLUMN password_hash TEXT NOT NULL,
ADD COLUMN role VARCHAR(255) DEFAULT 'users' NOT NULL,
ADD COLUMN created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
ADD COLUMN updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL;

CREATE OR REPLACE FUNCTION update_updated_at_column()
Returns TRIGGER AS $$
BEGIN 
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_modtime
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column()