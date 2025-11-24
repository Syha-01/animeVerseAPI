package main

import (
	"errors"
	"net/http"

	"github.com/Syha-01/animeVerseAPI/internal/data"
	"github.com/Syha-01/animeVerseAPI/internal/validator"
)

func (a *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var incomingData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Username:  incomingData.Username,
		Email:     incomingData.Email,
		Activated: false,
	}

	err = user.Password.Set(incomingData.Password)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateUser(v, user)

	if !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			a.failedValidationResponse(w, r, v.Errors)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// Send a 201 Created status code
	err = a.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
