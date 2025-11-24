package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFoundResponse(w, r)
	})
	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.methodNotAllowedResponse(w, r)
	})

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/animes", app.listAnimesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/animes", app.createAnimeHandler)
	router.HandlerFunc(http.MethodGet, "/v1/animes/:id", app.displayAnimeHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/animes/:id", app.updateAnimeHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/animes/:id", app.deleteAnimeHandler)

	router.HandlerFunc(http.MethodGet, "/v1/user_anime_list", app.listUserAnimeListsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/user_anime_list", app.createUserAnimeListHandler)
	router.HandlerFunc(http.MethodGet, "/v1/user_anime_list/:id", app.displayUserAnimeListHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/user_anime_list/:id", app.updateUserAnimeListHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/user_anime_list/:id", app.deleteUserAnimeListHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(router)))
}
