package data

import (
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
