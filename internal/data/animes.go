package data

import (
	"time"
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
