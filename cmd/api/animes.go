package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Syha-01/animeVerseAPI/internal/data"
)

func (a *application) createAnimeHandler(w http.ResponseWriter, r *http.Request) {
	var input data.Anime

	// perform the decoding
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		a.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// for now display the result
	fmt.Fprintf(w, "%+v\n", input)
}
