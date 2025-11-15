package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Syha-01/animeVerseAPI/internal/data"
	"github.com/Syha-01/animeVerseAPI/internal/validator"
)

func (a *application) createAnimeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID                   int64    `json:"id"`
		Title                string   `json:"title"`
		Synopsis             string   `json:"synopsis"`
		CoverImageURL        string   `json:"cover_image_url"`
		TotalEpisodes        int32    `json:"total_episodes"`
		Status               string   `json:"status"`
		ReleaseDate          string   `json:"release_date"`
		Rating               string   `json:"rating"`
		Score                float32  `json:"score"`
		Genres               []string `json:"genres"`
		Studios              []string `json:"studios"`
		BroadcastInformation string   `json:"broadcast_information"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	releaseDate, err := time.Parse("2006-01-02", input.ReleaseDate)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	anime := &data.Anime{
		ID:                   input.ID,
		Title:                input.Title,
		Synopsis:             input.Synopsis,
		CoverImageURL:        input.CoverImageURL,
		TotalEpisodes:        input.TotalEpisodes,
		Status:               input.Status,
		ReleaseDate:          releaseDate,
		Rating:               input.Rating,
		Score:                input.Score,
		Genres:               input.Genres,
		Studios:              input.Studios,
		BroadcastInformation: input.BroadcastInformation,
	}

	v := validator.New()

	if data.ValidateAnime(v, anime); !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.models.Animes.Insert(anime)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/animes/%d", anime.ID))

	err = a.writeJSON(w, http.StatusCreated, envelope{"anime": anime}, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *application) displayAnimeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	anime, err := a.models.Animes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"anime": anime}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
