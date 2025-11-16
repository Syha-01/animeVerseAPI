package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Syha-01/animeVerseAPI/internal/validator"
)

type Anime struct {
	ID                   int64     `json:"id"`
	Title                string    `json:"title"`
	Synopsis             string    `json:"synopsis,omitempty"`
	CoverImageURL        string    `json:"cover_image_url,omitempty"`
	TotalEpisodes        int32     `json:"total_episodes,omitempty"`
	Status               string    `json:"status,omitempty"`
	ReleaseDate          time.Time `json:"release_date,omitempty"`
	Rating               string    `json:"rating,omitempty"`
	Score                float32   `json:"score,omitempty"`
	Genres               []string  `json:"genres,omitempty"`
	Studios              []string  `json:"studios,omitempty"`
	BroadcastInformation string    `json:"broadcast_information,omitempty"`
	JikanLastSyncedAt    time.Time `json:"-"`
}

func ValidateAnime(v *validator.Validator, anime *Anime) {
	v.Check(anime.Title != "", "title", "must be provided")
	v.Check(len(anime.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(anime.TotalEpisodes >= 0, "total_episodes", "must be a positive integer")
	v.Check(anime.Score >= 0, "score", "must be a positive number")
	v.Check(anime.Score <= 10, "score", "must not be greater than 10")
}

type AnimeModel struct {
	DB *sql.DB
}

func (m AnimeModel) InsertAnime(anime *Anime) error {
	query := `
		INSERT INTO anime (id, title, synopsis, cover_image_url, total_episodes, status, release_date, rating, score, genres, studios, broadcast_information)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING jikan_last_synced_at`

	genresJSON, err := json.Marshal(anime.Genres)
	if err != nil {
		return err
	}

	studiosJSON, err := json.Marshal(anime.Studios)
	if err != nil {
		return err
	}

	args := []any{
		anime.ID,
		anime.Title,
		anime.Synopsis,
		anime.CoverImageURL,
		anime.TotalEpisodes,
		anime.Status,
		anime.ReleaseDate,
		anime.Rating,
		anime.Score,
		genresJSON,
		studiosJSON,
		anime.BroadcastInformation,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&anime.JikanLastSyncedAt)
}

func (m AnimeModel) GetAnime(id int64) (*Anime, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, title, synopsis, cover_image_url, total_episodes, status, release_date, rating, score, genres, studios, broadcast_information, jikan_last_synced_at
		FROM anime
		WHERE id = $1`

	var anime Anime
	var genres, studios []byte

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&anime.ID,
		&anime.Title,
		&anime.Synopsis,
		&anime.CoverImageURL,
		&anime.TotalEpisodes,
		&anime.Status,
		&anime.ReleaseDate,
		&anime.Rating,
		&anime.Score,
		&genres,
		&studios,
		&anime.BroadcastInformation,
		&anime.JikanLastSyncedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	if err := json.Unmarshal(genres, &anime.Genres); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(studios, &anime.Studios); err != nil {
		return nil, err
	}

	return &anime, nil
}

func (m AnimeModel) UpdateAnime(anime *Anime) error {
	query := `
		UPDATE anime
		SET title = $1, synopsis = $2, cover_image_url = $3, total_episodes = $4, status = $5, release_date = $6, rating = $7, score = $8, genres = $9, studios = $10, broadcast_information = $11
		WHERE id = $12
		RETURNING jikan_last_synced_at`

	genresJSON, err := json.Marshal(anime.Genres)
	if err != nil {
		return err
	}

	studiosJSON, err := json.Marshal(anime.Studios)
	if err != nil {
		return err
	}

	args := []any{
		anime.Title,
		anime.Synopsis,
		anime.CoverImageURL,
		anime.TotalEpisodes,
		anime.Status,
		anime.ReleaseDate,
		anime.Rating,
		anime.Score,
		genresJSON,
		studiosJSON,
		anime.BroadcastInformation,
		anime.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&anime.JikanLastSyncedAt)
}

func (m AnimeModel) DeleteAnime(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM anime
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// GetAllAnimes returns a slice of animes, filtered by title and genres.
func (m AnimeModel) GetAllAnimes(title string, genres []string) ([]*Anime, error) {
	// Start with the base query. The `WHERE 1=1` is a useful trick to make it easy
	// to append additional `AND` clauses dynamically.
	query := `
		SELECT id, title, synopsis, cover_image_url, total_episodes, status, release_date, rating, score, genres, studios, broadcast_information, jikan_last_synced_at
		FROM anime
		WHERE 1=1`

	// Create a slice to hold the query arguments.
	args := []any{}
	argCount := 1

	// If a title is provided, add a full-text search clause.
	// We use PostgreSQL's full-text search operators `to_tsvector` and `plainto_tsquery`.
	if title != "" {
		query += fmt.Sprintf(" AND (to_tsvector('simple', title) @@ plainto_tsquery('simple', $%d))", argCount)
		args = append(args, title)
		argCount++
	}

	// If genres are provided, add a clause to check if the `genres` jsonb array
	// contains all the provided genres. The `@>` operator means "contains".
	if len(genres) > 0 {
		genresJSON, err := json.Marshal(genres)
		if err != nil {
			return nil, err
		}
		query += fmt.Sprintf(" AND (genres @> $%d)", argCount)
		args = append(args, genresJSON)
		argCount++
	}

	// Always order by id for consistent results.
	query += " ORDER BY id"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Pass the dynamically built query and arguments to QueryContext.
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// The rest of the scanning logic remains the same.
	animes := []*Anime{}

	for rows.Next() {
		var anime Anime
		var synopsis sql.NullString
		var coverImageURL sql.NullString
		var totalEpisodes sql.NullInt32
		var status sql.NullString
		var releaseDate sql.NullTime
		var rating sql.NullString
		var score sql.NullFloat64
		var broadcastInformation sql.NullString
		var genresBytes, studiosBytes []byte // Renamed to avoid shadowing

		err := rows.Scan(
			&anime.ID,
			&anime.Title,
			&synopsis,
			&coverImageURL,
			&totalEpisodes,
			&status,
			&releaseDate,
			&rating,
			&score,
			&genresBytes,
			&studiosBytes,
			&broadcastInformation,
			&anime.JikanLastSyncedAt,
		)
		if err != nil {
			return nil, err
		}

		if synopsis.Valid {
			anime.Synopsis = synopsis.String
		}
		if coverImageURL.Valid {
			anime.CoverImageURL = coverImageURL.String
		}
		if totalEpisodes.Valid {
			anime.TotalEpisodes = totalEpisodes.Int32
		}
		if status.Valid {
			anime.Status = status.String
		}
		if releaseDate.Valid {
			anime.ReleaseDate = releaseDate.Time
		}
		if rating.Valid {
			anime.Rating = rating.String
		}
		if score.Valid {
			anime.Score = float32(score.Float64)
		}
		if broadcastInformation.Valid {
			anime.BroadcastInformation = broadcastInformation.String
		}

		if err := json.Unmarshal(genresBytes, &anime.Genres); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(studiosBytes, &anime.Studios); err != nil {
			return nil, err
		}

		animes = append(animes, &anime)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return animes, nil
}
