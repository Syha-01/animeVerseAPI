CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- A trigger function to automatically update the `updated_at` timestamp whenever a row is modified.
-- This is a common and highly useful pattern in database design.
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the users table
-- This table stores all user profile and authentication information.
CREATE TABLE users (
    -- Using UUID as a primary key is great for security and scalability.
    -- gen_random_uuid() comes from the pgcrypto extension.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Unique constraint to prevent duplicate usernames.
    username VARCHAR(255) NOT NULL UNIQUE,

    -- Unique constraint for email, which is critical for password resets.
    email VARCHAR(255) NOT NULL UNIQUE,

    -- This will store the securely hashed password (e.g., using bcrypt).
    -- NEVER store plain-text passwords.
    password_hash VARCHAR(255) NOT NULL,

    -- URL for a user's avatar. Can be NULL if they don't set one.
    avatar_url VARCHAR(255),

    -- A short biography for the user's public profile.
    bio TEXT,

    -- Timestamps with time zone are crucial for applications with a global user base.
    -- `now()` automatically sets the timestamp on creation.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

-- Apply the trigger to automatically update the `updated_at` timestamp on the users table.
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Create indexes for faster lookups during login and registration.
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
