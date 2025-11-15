package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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

func (m AnimeModel) Insert(anime *Anime) error {
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

func (m AnimeModel) Get(id int64) (*Anime, error) {
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
