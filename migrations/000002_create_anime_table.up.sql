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

-- Add an index on the title for fast text-based searching.
-- A GIN index is also excellent for searching within the JSONB columns.
CREATE INDEX idx_anime_title ON anime(title);
CREATE INDEX idx_anime_genres ON anime USING GIN(genres);
CREATE INDEX idx_anime_studios ON anime USING GIN(studios);
