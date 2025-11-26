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

	router.HandlerFunc(http.MethodGet, "/v1/animes", app.requirePermission("animes:read", app.listAnimesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/animes", app.requirePermission("animes:write", app.createAnimeHandler))
	router.HandlerFunc(http.MethodGet, "/v1/animes/:id", app.requirePermission("animes:read", app.displayAnimeHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/animes/:id", app.requirePermission("animes:write", app.updateAnimeHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/animes/:id", app.requirePermission("animes:write", app.deleteAnimeHandler))

	router.HandlerFunc(http.MethodGet, "/v1/user_anime_list", app.requireActivatedUser(app.listUserAnimeListsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/user_anime_list", app.requireActivatedUser(app.createUserAnimeListHandler))
	router.HandlerFunc(http.MethodGet, "/v1/user_anime_list/:id", app.requireActivatedUser(app.displayUserAnimeListHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/user_anime_list/:id", app.requireActivatedUser(app.updateUserAnimeListHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/user_anime_list/:id", app.requireActivatedUser(app.deleteUserAnimeListHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
}
