package validator

import (
	"regexp"
	"testing"
)

func TestNew(t *testing.T) {
	v := New()
	
	if v == nil {
		t.Fatal("expected non-nil Validator")
	}
	
	if v.Errors == nil {
		t.Error("expected Errors map to be initialized")
	}
	
	if len(v.Errors) != 0 {
		t.Errorf("expected empty Errors map, got %d entries", len(v.Errors))
	}
}

func TestAddError(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*Validator)
		key      string
		message  string
		expected map[string]string
	}{
		{
			name:     "add first error",
			setup:    func(v *Validator) {},
			key:      "email",
			message:  "must be a valid email",
			expected: map[string]string{"email": "must be a valid email"},
		},
		{
			name: "don't overwrite existing error",
			setup: func(v *Validator) {
				v.AddError("email", "first error")
			},
			key:      "email",
			message:  "second error",
			expected: map[string]string{"email": "first error"},
		},
		{
			name: "add multiple different errors",
			setup: func(v *Validator) {
				v.AddError("email", "invalid email")
			},
			key:      "password",
			message:  "too short",
			expected: map[string]string{"email": "invalid email", "password": "too short"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			tt.setup(v)
			v.AddError(tt.key, tt.message)
			
			if len(v.Errors) != len(tt.expected) {
				t.Errorf("expected %d errors, got %d", len(tt.expected), len(v.Errors))
			}
			
			for key, expectedMsg := range tt.expected {
				if actualMsg, exists := v.Errors[key]; !exists {
					t.Errorf("expected error with key %q to exist", key)
				} else if actualMsg != expectedMsg {
					t.Errorf("expected message %q for key %q, got %q", expectedMsg, key, actualMsg)
				}
			}
		})
	}
}

func TestCheck(t *testing.T) {
	tests := []struct {
		name          string
		ok            bool
		key           string
		message       string
		expectError   bool
		expectedMsg   string
	}{
		{
			name:        "validation passes - no error added",
			ok:          true,
			key:         "email",
			message:     "must be valid",
			expectError: false,
		},
		{
			name:        "validation fails - error added",
			ok:          false,
			key:         "email",
			message:     "must be valid",
			expectError: true,
			expectedMsg: "must be valid",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			v.Check(tt.ok, tt.key, tt.message)
			
			msg, exists := v.Errors[tt.key]
			
			if tt.expectError {
				if !exists {
					t.Error("expected error to be added but none found")
				}
				if msg != tt.expectedMsg {
					t.Errorf("expected message %q, got %q", tt.expectedMsg, msg)
				}
			} else {
				if exists {
					t.Errorf("expected no error, but found: %q", msg)
				}
			}
		})
	}
}

func TestValid(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*Validator)
		expected bool
	}{
		{
			name:     "valid when no errors",
			setup:    func(v *Validator) {},
			expected: true,
		},
		{
			name: "invalid when errors exist",
			setup: func(v *Validator) {
				v.AddError("email", "invalid")
			},
			expected: false,
		},
		{
			name: "invalid with multiple errors",
			setup: func(v *Validator) {
				v.AddError("email", "invalid")
				v.AddError("password", "too short")
			},
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			tt.setup(v)
			
			if got := v.Valid(); got != tt.expected {
				t.Errorf("Valid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*Validator)
		expected bool
	}{
		{
			name:     "empty when no errors",
			setup:    func(v *Validator) {},
			expected: true,
		},
		{
			name: "not empty when errors exist",
			setup: func(v *Validator) {
				v.AddError("email", "invalid")
			},
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			tt.setup(v)
			
			if got := v.IsEmpty(); got != tt.expected {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMatches(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		pattern  *regexp.Regexp
		expected bool
	}{
		{
			name:     "valid email",
			value:    "test@example.com",
			pattern:  EmailRX,
			expected: true,
		},
		{
			name:     "invalid email - no @",
			value:    "testexample.com",
			pattern:  EmailRX,
			expected: false,
		},
		{
			name:     "invalid email - no domain",
			value:    "test@",
			pattern:  EmailRX,
			expected: false,
		},
		{
			name:     "valid email with subdomain",
			value:    "user@mail.example.com",
			pattern:  EmailRX,
			expected: true,
		},
		{
			name:     "valid email with special characters",
			value:    "user+tag@example.com",
			pattern:  EmailRX,
			expected: true,
		},
		{
			name:     "empty string",
			value:    "",
			pattern:  EmailRX,
			expected: false,
		},
		{
			name:     "custom pattern - digits only",
			value:    "12345",
			pattern:  regexp.MustCompile("^[0-9]+$"),
			expected: true,
		},
		{
			name:     "custom pattern - digits only fail",
			value:    "123a45",
			pattern:  regexp.MustCompile("^[0-9]+$"),
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Matches(tt.value, tt.pattern); got != tt.expected {
				t.Errorf("Matches(%q, pattern) = %v, want %v", tt.value, got, tt.expected)
			}
		})
	}
}

func TestPermittedValue(t *testing.T) {
	tests := []struct {
		name            string
		value           string
		permittedValues []string
		expected        bool
	}{
		{
			name:            "value in list",
			value:           "admin",
			permittedValues: []string{"admin", "user", "guest"},
			expected:        true,
		},
		{
			name:            "value not in list",
			value:           "superuser",
			permittedValues: []string{"admin", "user", "guest"},
			expected:        false,
		},
		{
			name:            "empty list",
			value:           "admin",
			permittedValues: []string{},
			expected:        false,
		},
		{
			name:            "single value match",
			value:           "only",
			permittedValues: []string{"only"},
			expected:        true,
		},
		{
			name:            "case sensitive - no match",
			value:           "Admin",
			permittedValues: []string{"admin"},
			expected:        false,
		},
		{
			name:            "empty string in list",
			value:           "",
			permittedValues: []string{"", "value"},
			expected:        true,
		},
		{
			name:            "numeric strings",
			value:           "1",
			permittedValues: []string{"1", "2", "3"},
			expected:        true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PermittedValue(tt.value, tt.permittedValues...); got != tt.expected {
				t.Errorf("PermittedValue(%q, %v) = %v, want %v", tt.value, tt.permittedValues, got, tt.expected)
			}
		})
	}
}
