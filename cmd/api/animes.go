package main

import (
	"fmt"
	"net/http"

	"github.com/Syha-01/animeVerseAPI/internal/data"
)

func (a *application) createAnimeHandler(w http.ResponseWriter, r *http.Request) {
	var input data.Anime

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// for now display the result
	fmt.Fprintf(w, "%+v\n", input)
}
