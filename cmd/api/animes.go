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

	err = a.models.Animes.InsertAnime(anime)
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

	anime, err := a.models.Animes.GetAnime(id)
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

func (a *application) updateAnimeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	anime, err := a.models.Animes.GetAnime(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title                *string  `json:"title"`
		Synopsis             *string  `json:"synopsis"`
		CoverImageURL        *string  `json:"cover_image_url"`
		TotalEpisodes        *int32   `json:"total_episodes"`
		Status               *string  `json:"status"`
		ReleaseDate          *string  `json:"release_date"`
		Rating               *string  `json:"rating"`
		Score                *float32 `json:"score"`
		Genres               []string `json:"genres"`
		Studios              []string `json:"studios"`
		BroadcastInformation *string  `json:"broadcast_information"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		anime.Title = *input.Title
	}
	if input.Synopsis != nil {
		anime.Synopsis = *input.Synopsis
	}
	if input.CoverImageURL != nil {
		anime.CoverImageURL = *input.CoverImageURL
	}
	if input.TotalEpisodes != nil {
		anime.TotalEpisodes = *input.TotalEpisodes
	}
	if input.Status != nil {
		anime.Status = *input.Status
	}
	if input.ReleaseDate != nil {
		releaseDate, err := time.Parse("2006-01-02", *input.ReleaseDate)
		if err != nil {
			a.badRequestResponse(w, r, err)
			return
		}
		anime.ReleaseDate = releaseDate
	}
	if input.Rating != nil {
		anime.Rating = *input.Rating
	}
	if input.Score != nil {
		anime.Score = *input.Score
	}
	if input.Genres != nil {
		anime.Genres = input.Genres
	}
	if input.Studios != nil {
		anime.Studios = input.Studios
	}
	if input.BroadcastInformation != nil {
		anime.BroadcastInformation = *input.BroadcastInformation
	}

	v := validator.New()

	if data.ValidateAnime(v, anime); !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.models.Animes.UpdateAnime(anime)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"anime": anime}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *application) deleteAnimeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.models.Animes.DeleteAnime(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"message": "anime successfully deleted"}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *application) listAnimesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}

	v := validator.New()
	queryParameters := r.URL.Query()

	input.Title = a.getSingleQueryParameter(queryParameters, "title", "")
	input.Genres = a.getMultipleQueryParameters(queryParameters, "genres", []string{})
	input.Filters.Page = a.getSingleIntegerParameter(queryParameters, "page", 1, v)
	input.Filters.PageSize = a.getSingleIntegerParameter(queryParameters, "page_size", 20, v)

	data.ValidateFilters(v, input.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Update the call to receive the metadata.
	animes, metadata, err := a.models.Animes.GetAllAnimes(input.Title, input.Genres, input.Filters)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Include metadata in the response envelope.
	err = a.writeJSON(w, http.StatusOK, envelope{
		"animes":    animes,
		"@metadata": metadata,
	}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
