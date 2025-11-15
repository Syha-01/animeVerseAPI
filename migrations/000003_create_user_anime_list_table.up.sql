-- First, create a custom ENUM type for the user's tracking status.
-- This enforces data integrity and is more efficient than using a VARCHAR.
CREATE TYPE user_anime_status AS ENUM (
    'Watching',
    'Completed',
    'Dropped',
    'Watch later'
);

-- Create the user_anime_list table
-- This is the core "join" table that tracks each user's relationship with each anime on their list.
CREATE TABLE user_anime_list (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign key linking to the users table.
    -- ON DELETE CASCADE means if a user is deleted, all their list entries are also deleted.
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Foreign key linking to the anime table.
    -- ON DELETE CASCADE means if an anime is removed from our DB, list entries are also removed.
    anime_id INTEGER NOT NULL REFERENCES anime(id) ON DELETE CASCADE,

    -- The user's progress status, using our custom ENUM type.
    status user_anime_status NOT NULL,

    -- The number of episodes the user has watched. Cannot be negative.
    current_episode INTEGER NOT NULL DEFAULT 0 CHECK (current_episode >= 0),

    -- The user's personal score, constrained to a 1-10 scale.
    score INTEGER CHECK (score >= 1 AND score <= 10),

    -- Date the user started watching the anime.
    started_watching_date DATE,

    -- Date the user finished the anime.
    finished_watching_date DATE,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    -- This is a critical constraint! It prevents a user from adding the same anime to their list more than once.
    UNIQUE (user_id, anime_id)
);

-- Apply the trigger to automatically update the `updated_at` timestamp on the user_anime_list table.
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON user_anime_list
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Create indexes for fast retrieval of a user's list or to see who has a specific anime.
CREATE INDEX idx_user_anime_list_user_id ON user_anime_list(user_id);
CREATE INDEX idx_user_anime_list_anime_id ON user_anime_list(anime_id);
