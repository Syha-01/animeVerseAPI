package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Syha-01/animeVerseAPI/internal/data"
	"github.com/Syha-01/animeVerseAPI/internal/validator"
)

func (a *application) createUserAnimeListHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID               string `json:"user_id"`
		AnimeID              int64  `json:"anime_id"`
		Status               string `json:"status"`
		CurrentEpisode       int32  `json:"current_episode"`
		Score                int32  `json:"score"`
		StartedWatchingDate  string `json:"started_watching_date"`
		FinishedWatchingDate string `json:"finished_watching_date"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	list := &data.UserAnimeList{
		UserID:         input.UserID,
		AnimeID:        input.AnimeID,
		Status:         input.Status,
		CurrentEpisode: input.CurrentEpisode,
		Score:          input.Score,
	}

	if input.StartedWatchingDate != "" {
		startedWatchingDate, err := time.Parse("2006-01-02", input.StartedWatchingDate)
		if err != nil {
			a.badRequestResponse(w, r, err)
			return
		}
		list.StartedWatchingDate = startedWatchingDate
	}

	if input.FinishedWatchingDate != "" {
		finishedWatchingDate, err := time.Parse("2006-01-02", input.FinishedWatchingDate)
		if err != nil {
			a.badRequestResponse(w, r, err)
			return
		}
		list.FinishedWatchingDate = finishedWatchingDate
	}

	v := validator.New()

	if data.ValidateUserAnimeList(v, list); !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.models.UserAnimeList.InsertUserAnimeList(list)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/user_anime_list/%s", list.ID))

	err = a.writeJSON(w, http.StatusCreated, envelope{"user_anime_list": list}, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *application) displayUserAnimeListHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readUUIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	list, err := a.models.UserAnimeList.GetUserAnimeList(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"user_anime_list": list}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *application) updateUserAnimeListHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readUUIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	list, err := a.models.UserAnimeList.GetUserAnimeList(id)
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
		Status               *string `json:"status"`
		CurrentEpisode       *int32  `json:"current_episode"`
		Score                *int32  `json:"score"`
		StartedWatchingDate  *string `json:"started_watching_date"`
		FinishedWatchingDate *string `json:"finished_watching_date"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if input.Status != nil {
		list.Status = *input.Status
	}
	if input.CurrentEpisode != nil {
		list.CurrentEpisode = *input.CurrentEpisode
	}
	if input.Score != nil {
		list.Score = *input.Score
	}
	if input.StartedWatchingDate != nil {
		startedWatchingDate, err := time.Parse("2006-01-02", *input.StartedWatchingDate)
		if err != nil {
			a.badRequestResponse(w, r, err)
			return
		}
		list.StartedWatchingDate = startedWatchingDate
	}
	if input.FinishedWatchingDate != nil {
		finishedWatchingDate, err := time.Parse("2006-01-02", *input.FinishedWatchingDate)
		if err != nil {
			a.badRequestResponse(w, r, err)
			return
		}
		list.FinishedWatchingDate = finishedWatchingDate
	}

	v := validator.New()

	if data.ValidateUserAnimeList(v, list); !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.models.UserAnimeList.UpdateUserAnimeList(list)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"user_anime_list": list}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *application) deleteUserAnimeListHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readUUIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.models.UserAnimeList.DeleteUserAnimeList(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"message": "user anime list successfully deleted"}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *application) listUserAnimeListsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID string
		Status string
		data.Filters
	}

	v := validator.New()
	queryParameters := r.URL.Query()

	input.UserID = a.getSingleQueryParameter(queryParameters, "user_id", "")
	input.Status = a.getSingleQueryParameter(queryParameters, "status", "")
	input.Filters.Page = a.getSingleIntegerParameter(queryParameters, "page", 1, v)
	input.Filters.PageSize = a.getSingleIntegerParameter(queryParameters, "page_size", 20, v)
	input.Filters.Sort = a.getSingleQueryParameter(queryParameters, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "user_id", "anime_id", "status", "score", "created_at", "-id", "-user_id", "-anime_id", "-status", "-score", "-created_at"}

	data.ValidateFilters(v, input.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	lists, metadata, err := a.models.UserAnimeList.GetAllUserAnimeLists(input.UserID, input.Status, input.Filters)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{
		"user_anime_lists": lists,
		"@metadata":        metadata,
	}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
