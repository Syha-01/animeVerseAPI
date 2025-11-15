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

-- Apply the trigger to automatically update the `updated_at` timestamp on the reviews table.
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON reviews
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Create indexes for fast retrieval of reviews.
CREATE INDEX idx_reviews_user_id ON reviews(user_id);
CREATE INDEX idx_reviews_anime_id ON reviews(anime_id);
