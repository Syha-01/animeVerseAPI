package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Syha-01/animeVerseAPI/internal/validator"
	"github.com/julienschmidt/httprouter"
)

func TestReadIDParam(t *testing.T) {
	tests := []struct {
		name        string
		paramValue  string
		expectedID  int64
		expectError bool
	}{
		{
			name:        "valid positive ID",
			paramValue:  "123",
			expectedID:  123,
			expectError: false,
		},
		{
			name:        "valid large ID",
			paramValue:  "999999",
			expectedID:  999999,
			expectError: false,
		},
		{
			name:        "ID equal to 1",
			paramValue:  "1",
			expectedID:  1,
			expectError: false,
		},
		{
			name:        "zero ID - invalid",
			paramValue:  "0",
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "negative ID - invalid",
			paramValue:  "-5",
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "non-numeric string",
			paramValue:  "abc",
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "empty string",
			paramValue:  "",
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "decimal number",
			paramValue:  "12.5",
			expectedID:  0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock application
			app := &application{}

			// Create a request with httprouter params
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			params := httprouter.Params{
				httprouter.Param{Key: "id", Value: tt.paramValue},
			}
			ctx := context.WithValue(req.Context(), httprouter.ParamsKey, params)
			req = req.WithContext(ctx)

			id, err := app.readIDParam(req)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if id != tt.expectedID {
					t.Errorf("expected ID %d, got %d", tt.expectedID, id)
				}
			}
		})
	}
}

func TestReadUUIDParam(t *testing.T) {
	tests := []struct {
		name        string
		paramValue  string
		expectedID  string
		expectError bool
	}{
		{
			name:        "valid UUID",
			paramValue:  "550e8400-e29b-41d4-a716-446655440000",
			expectedID:  "550e8400-e29b-41d4-a716-446655440000",
			expectError: false,
		},
		{
			name:        "valid alphanumeric string",
			paramValue:  "abc123def456",
			expectedID:  "abc123def456",
			expectError: false,
		},
		{
			name:        "short string",
			paramValue:  "test",
			expectedID:  "test",
			expectError: false,
		},
		{
			name:        "empty string - invalid",
			paramValue:  "",
			expectedID:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock application
			app := &application{}

			// Create a request with httprouter params
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			params := httprouter.Params{
				httprouter.Param{Key: "id", Value: tt.paramValue},
			}
			ctx := context.WithValue(req.Context(), httprouter.ParamsKey, params)
			req = req.WithContext(ctx)

			id, err := app.readUUIDParam(req)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if id != tt.expectedID {
					t.Errorf("expected ID %q, got %q", tt.expectedID, id)
				}
			}
		})
	}
}

func TestGetSingleQueryParameter(t *testing.T) {
	tests := []struct {
		name         string
		queryString  string
		key          string
		defaultValue string
		expected     string
	}{
		{
			name:         "parameter exists",
			queryString:  "search=golang",
			key:          "search",
			defaultValue: "",
			expected:     "golang",
		},
		{
			name:         "parameter missing - use default",
			queryString:  "other=value",
			key:          "search",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "empty query string",
			queryString:  "",
			key:          "search",
			defaultValue: "fallback",
			expected:     "fallback",
		},
		{
			name:         "parameter with special characters",
			queryString:  "query=hello%20world",
			key:          "query",
			defaultValue: "",
			expected:     "hello world",
		},
		{
			name:         "multiple parameters - get specific one",
			queryString:  "page=1&limit=10&sort=name",
			key:          "sort",
			defaultValue: "",
			expected:     "name",
		},
		{
			name:         "empty parameter value",
			queryString:  "search=",
			key:          "search",
			defaultValue: "default",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &application{}

			qs, _ := url.ParseQuery(tt.queryString)
			got := app.getSingleQueryParameter(qs, tt.key, tt.defaultValue)

			if got != tt.expected {
				t.Errorf("getSingleQueryParameter() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestGetMultipleQueryParameters(t *testing.T) {
	tests := []struct {
		name         string
		queryString  string
		key          string
		defaultValue []string
		expected     []string
	}{
		{
			name:         "comma-separated values",
			queryString:  "genres=action,comedy,drama",
			key:          "genres",
			defaultValue: []string{},
			expected:     []string{"action", "comedy", "drama"},
		},
		{
			name:         "single value",
			queryString:  "genres=action",
			key:          "genres",
			defaultValue: []string{},
			expected:     []string{"action"},
		},
		{
			name:         "parameter missing - use default",
			queryString:  "other=value",
			key:          "genres",
			defaultValue: []string{"default1", "default2"},
			expected:     []string{"default1", "default2"},
		},
		{
			name:         "empty query string",
			queryString:  "",
			key:          "genres",
			defaultValue: []string{"fallback"},
			expected:     []string{"fallback"},
		},
		{
			name:         "values with spaces",
			queryString:  "tags=sci-fi,action%20adventure,thriller",
			key:          "tags",
			defaultValue: []string{},
			expected:     []string{"sci-fi", "action adventure", "thriller"},
		},
		{
			name:         "empty value uses default",
			queryString:  "genres=",
			key:          "genres",
			defaultValue: []string{"default"},
			expected:     []string{"default"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &application{}

			qs, _ := url.ParseQuery(tt.queryString)
			got := app.getMultipleQueryParameters(qs, tt.key, tt.defaultValue)

			if len(got) != len(tt.expected) {
				t.Fatalf("expected %d values, got %d", len(tt.expected), len(got))
			}

			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("index %d: got %q, want %q", i, got[i], tt.expected[i])
				}
			}
		})
	}
}

func TestGetSingleIntegerParameter(t *testing.T) {
	tests := []struct {
		name         string
		queryString  string
		key          string
		defaultValue int
		expected     int
		expectError  bool
	}{
		{
			name:         "valid positive integer",
			queryString:  "page=5",
			key:          "page",
			defaultValue: 1,
			expected:     5,
			expectError:  false,
		},
		{
			name:         "valid zero",
			queryString:  "page=0",
			key:          "page",
			defaultValue: 1,
			expected:     0,
			expectError:  false,
		},
		{
			name:         "valid negative integer",
			queryString:  "offset=-10",
			key:          "offset",
			defaultValue: 0,
			expected:     -10,
			expectError:  false,
		},
		{
			name:         "parameter missing - use default",
			queryString:  "other=value",
			key:          "page",
			defaultValue: 1,
			expected:     1,
			expectError:  false,
		},
		{
			name:         "invalid - not a number",
			queryString:  "page=abc",
			key:          "page",
			defaultValue: 1,
			expected:     1,
			expectError:  true,
		},
		{
			name:         "invalid - decimal number",
			queryString:  "page=3.14",
			key:          "page",
			defaultValue: 1,
			expected:     1,
			expectError:  true,
		},
		{
			name:         "invalid - empty string",
			queryString:  "page=",
			key:          "page",
			defaultValue: 10,
			expected:     10,
			expectError:  false,
		},
		{
			name:         "large integer",
			queryString:  "limit=999999",
			key:          "limit",
			defaultValue: 100,
			expected:     999999,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &application{}
			v := validator.New()

			qs, _ := url.ParseQuery(tt.queryString)
			got := app.getSingleIntegerParameter(qs, tt.key, tt.defaultValue, v)

			if got != tt.expected {
				t.Errorf("getSingleIntegerParameter() = %d, want %d", got, tt.expected)
			}

			if tt.expectError {
				if v.Valid() {
					t.Error("expected validation error but got none")
				}
				if _, exists := v.Errors[tt.key]; !exists {
					t.Errorf("expected error for key %q but not found", tt.key)
				}
			} else {
				if !v.Valid() {
					t.Errorf("unexpected validation errors: %v", v.Errors)
				}
			}
		})
	}
}
