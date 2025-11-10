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

-- First, create a custom ENUM type for the user's tracking status.
-- This enforces data integrity and is more efficient than using a VARCHAR.
CREATE TYPE user_anime_status AS ENUM (
    'Watching',
    'Completed',
    'Dropped',
    'Watch later'
);

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

-- Create the anime table
-- This table serves as a local cache for all anime data fetched from the Jikan API.
CREATE TABLE anime (
    -- The primary key is the MyAnimeList ID, which simplifies syncing with the Jikan API.
    id INTEGER PRIMARY KEY,

    -- The main title of the anime.
    title VARCHAR(255) NOT NULL,

    -- A detailed plot summary.
    synopsis TEXT,

    -- URL for the cover art.
    cover_image_url VARCHAR(255),

    -- Total episodes. Can be NULL for shows that are still airing.
    total_episodes INTEGER,

    -- Airing status like 'Finished Airing', 'Currently Airing'.
    status VARCHAR(50),

    -- The initial air date of the anime.
    release_date DATE,

    -- Age rating, e.g., 'PG-13'.
    rating VARCHAR(50),

    -- The public score from MyAnimeList. Using DECIMAL for precision.
    score DECIMAL(4, 2),

    -- Using JSONB is highly efficient for storing and querying semi-structured data like genres.
    genres JSONB,

    -- Storing studio names in a JSONB array.
    studios JSONB,

    -- Broadcast information for currently airing shows. e.g., "Mondays at 01:00 (JST)".
    broadcast_information VARCHAR(255),

    -- A timestamp to track when we last synced this record with the Jikan API.
    -- This is crucial for keeping our local data fresh.
    jikan_last_synced_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
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

-- Create the reviews table
-- This table stores user reviews for anime.
CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign key linking to the users table.
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Foreign key linking to the anime table.
    anime_id INTEGER NOT NULL REFERENCES anime(id) ON DELETE CASCADE,

    -- The user's star rating, constrained to a 1-5 scale.
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),

    -- The user's review comment.
    comment TEXT,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    -- This prevents a user from reviewing the same anime more than once.
    UNIQUE (user_id, anime_id)
);

-- Apply the trigger to automatically update the `updated_at` timestamp on the users table.
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Apply the trigger to automatically update the `updated_at` timestamp on the user_anime_list table.
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON user_anime_list
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Apply the trigger to automatically update the `updated_at` timestamp on the reviews table.
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON reviews
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Create indexes for faster lookups during login and registration.
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

-- Add an index on the title for fast text-based searching.
-- A GIN index is also excellent for searching within the JSONB columns.
CREATE INDEX idx_anime_title ON anime(title);
CREATE INDEX idx_anime_genres ON anime USING GIN(genres);
CREATE INDEX idx_anime_studios ON anime USING GIN(studios);

-- Create indexes for fast retrieval of a user's list or to see who has a specific anime.
CREATE INDEX idx_user_anime_list_user_id ON user_anime_list(user_id);
CREATE INDEX idx_user_anime_list_anime_id ON user_anime_list(anime_id);

-- Create indexes for fast retrieval of reviews.
CREATE INDEX idx_reviews_user_id ON reviews(user_id);
CREATE INDEX idx_reviews_anime_id ON reviews(anime_id);
