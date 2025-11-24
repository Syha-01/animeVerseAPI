package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Syha-01/animeVerseAPI/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, destination any) error {
	// what is the max size of the request body (250KB seems reasonable)
	maxBytes := 256_000
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	// our decoder will check for unknown fields
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	// let start the decoding
	err := dec.Decode(destination)
	if err != nil {
		// check for the different errors
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("the body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		// Decode can also send back an io error message
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("the body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("the body contains the incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("the body contains the incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("the body must not be empty")
		// check for unknown field error
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		// does the body exceed our limit of 250KB?
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("the body must not be larger than %d bytes", maxBytesError.Limit)
		// the programmer messed up
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		// some other type of error
		default:
			return err
		}
	}
	// almost done. Let's lastly check if there is any data after
	// the valid JSON data. Maybe the person is trying to send
	// multiple request bodies during one request
	// We call decode once more to see if it gives us back anything
	// we use a throw away struct 'struct{}{}' to hold the result
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		// there is more data present
		return errors.New("the body must only contain a single JSON value")
	}
	return nil
}

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) readUUIDParam(r *http.Request) (string, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id := params.ByName("id")
	if id == "" {
		return "", errors.New("invalid id parameter")
	}

	return id, nil
}

// getSingleQueryParameter reads a string value from the query string.
// If the key doesn't exist, it returns the provided default value.
func (app *application) getSingleQueryParameter(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

// getMultipleQueryParameters reads a comma-separated list of strings from the query string.
// If the key doesn't exist, it returns the provided default value.
func (app *application) getMultipleQueryParameters(qs url.Values, key string, defaultValue []string) []string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return strings.Split(s, ",")
}

// getSingleIntegerParameter reads an integer value from the query string.
// If the key doesn't exist or is invalid, it adds an error to the validator and returns the default value.
func (app *application) getSingleIntegerParameter(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	return i
}

func (app *application) background(fn func()) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				app.logger.Error(fmt.Sprintf("%v", err))
			}
		}()

		fn()
	}()
}
