package main

import (
	"fmt"
	"net/http"

	"github.com/Syha-01/animeVerseAPI/internal/data"
	"github.com/Syha-01/animeVerseAPI/internal/validator"
)

func (a *application) createAnimeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title         string   `json:"title"`
		TotalEpisodes int32    `json:"total_episodes"`
		Score         float32  `json:"score"`
		Genres        []string `json:"genres"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	anime := &data.Anime{
		Title:         input.Title,
		TotalEpisodes: input.TotalEpisodes,
		Score:         input.Score,
		Genres:        input.Genres,
	}

	v := validator.New()

	if data.ValidateAnime(v, anime); !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}
