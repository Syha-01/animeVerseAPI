package data

import (
	"testing"

	"github.com/Syha-01/animeVerseAPI/internal/validator"
)

func TestValidateFilters(t *testing.T) {
	tests := []struct {
		name           string
		filters        Filters
		expectValid    bool
		expectedErrors map[string]string
	}{
		{
			name: "valid filters",
			filters: Filters{
				Page:         1,
				PageSize:     10,
				Sort:         "id",
				SortSafeList: []string{"id", "title", "-created_at"},
			},
			expectValid: true,
		},
		{
			name: "page too low",
			filters: Filters{
				Page:         0,
				PageSize:     10,
				Sort:         "id",
				SortSafeList: []string{"id"},
			},
			expectValid:    false,
			expectedErrors: map[string]string{"page": "must be greater than zero"},
		},
		{
			name: "page too high",
			filters: Filters{
				Page:         1001,
				PageSize:     10,
				Sort:         "id",
				SortSafeList: []string{"id"},
			},
			expectValid:    false,
			expectedErrors: map[string]string{"page": "must be a maximum of 1000"},
		},
		{
			name: "page_size too low",
			filters: Filters{
				Page:         1,
				PageSize:     0,
				Sort:         "id",
				SortSafeList: []string{"id"},
			},
			expectValid:    false,
			expectedErrors: map[string]string{"page_size": "must be greater than zero"},
		},
		{
			name: "page_size too high",
			filters: Filters{
				Page:         1,
				PageSize:     101,
				Sort:         "id",
				SortSafeList: []string{"id"},
			},
			expectValid:    false,
			expectedErrors: map[string]string{"page_size": "must be a maximum of 100"},
		},
		{
			name: "invalid sort value",
			filters: Filters{
				Page:         1,
				PageSize:     10,
				Sort:         "malicious_field",
				SortSafeList: []string{"id", "title"},
			},
			expectValid:    false,
			expectedErrors: map[string]string{"sort": "invalid sort value"},
		},
		{
			name: "multiple validation errors",
			filters: Filters{
				Page:         0,
				PageSize:     101,
				Sort:         "invalid",
				SortSafeList: []string{"id"},
			},
			expectValid: false,
			expectedErrors: map[string]string{
				"page":      "must be greater than zero",
				"page_size": "must be a maximum of 100",
				"sort":      "invalid sort value",
			},
		},
		{
			name: "valid descending sort",
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "-created_at",
				SortSafeList: []string{"id", "-created_at"},
			},
			expectValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			ValidateFilters(v, tt.filters)

			if tt.expectValid {
				if !v.Valid() {
					t.Errorf("expected valid filters, but got errors: %v", v.Errors)
				}
			} else {
				if v.Valid() {
					t.Error("expected validation errors, but got none")
				}
				for key, expectedMsg := range tt.expectedErrors {
					if actualMsg, exists := v.Errors[key]; !exists {
						t.Errorf("expected error for key %q, but not found", key)
					} else if actualMsg != expectedMsg {
						t.Errorf("expected error message %q for key %q, got %q", expectedMsg, key, actualMsg)
					}
				}
			}
		})
	}
}

func TestFilters_sortColumn(t *testing.T) {
	tests := []struct {
		name        string
		filters     Filters
		expected    string
		expectPanic bool
	}{
		{
			name: "simple column name",
			filters: Filters{
				Sort:         "id",
				SortSafeList: []string{"id", "title"},
			},
			expected: "id",
		},
		{
			name: "descending sort - strip prefix",
			filters: Filters{
				Sort:         "-created_at",
				SortSafeList: []string{"id", "-created_at"},
			},
			expected: "created_at",
		},
		{
			name: "another descending sort",
			filters: Filters{
				Sort:         "-title",
				SortSafeList: []string{"id", "-title", "created_at"},
			},
			expected: "title",
		},
		{
			name: "unsafe sort parameter - should panic",
			filters: Filters{
				Sort:         "malicious_field",
				SortSafeList: []string{"id", "title"},
			},
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expected panic but didn't get one")
					}
				}()
				tt.filters.sortColumn()
			} else {
				got := tt.filters.sortColumn()
				if got != tt.expected {
					t.Errorf("sortColumn() = %q, want %q", got, tt.expected)
				}
			}
		})
	}
}

func TestFilters_sortDirection(t *testing.T) {
	tests := []struct {
		name     string
		sort     string
		expected string
	}{
		{
			name:     "ascending - no prefix",
			sort:     "id",
			expected: "ASC",
		},
		{
			name:     "descending - with prefix",
			sort:     "-created_at",
			expected: "DESC",
		},
		{
			name:     "ascending - title",
			sort:     "title",
			expected: "ASC",
		},
		{
			name:     "descending - title",
			sort:     "-title",
			expected: "DESC",
		},
		{
			name:     "empty string",
			sort:     "",
			expected: "ASC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Filters{Sort: tt.sort}
			got := f.sortDirection()
			if got != tt.expected {
				t.Errorf("sortDirection() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestFilters_limit(t *testing.T) {
	tests := []struct {
		name     string
		pageSize int
		expected int
	}{
		{
			name:     "page size 10",
			pageSize: 10,
			expected: 10,
		},
		{
			name:     "page size 1",
			pageSize: 1,
			expected: 1,
		},
		{
			name:     "page size 100",
			pageSize: 100,
			expected: 100,
		},
		{
			name:     "page size 50",
			pageSize: 50,
			expected: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Filters{PageSize: tt.pageSize}
			got := f.limit()
			if got != tt.expected {
				t.Errorf("limit() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestFilters_offset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		expected int
	}{
		{
			name:     "first page",
			page:     1,
			pageSize: 10,
			expected: 0,
		},
		{
			name:     "second page",
			page:     2,
			pageSize: 10,
			expected: 10,
		},
		{
			name:     "third page",
			page:     3,
			pageSize: 20,
			expected: 40,
		},
		{
			name:     "page 10 with size 25",
			page:     10,
			pageSize: 25,
			expected: 225,
		},
		{
			name:     "page 1 with size 100",
			page:     1,
			pageSize: 100,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Filters{
				Page:     tt.page,
				PageSize: tt.pageSize,
			}
			got := f.offset()
			if got != tt.expected {
				t.Errorf("offset() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestCalculateMetadata(t *testing.T) {
	tests := []struct {
		name         string
		totalRecords int
		page         int
		pageSize     int
		expected     Metadata
	}{
		{
			name:         "no records",
			totalRecords: 0,
			page:         1,
			pageSize:     10,
			expected:     Metadata{},
		},
		{
			name:         "single page",
			totalRecords: 5,
			page:         1,
			pageSize:     10,
			expected: Metadata{
				CurrentPage:  1,
				PageSize:     10,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 5,
			},
		},
		{
			name:         "multiple pages - first page",
			totalRecords: 25,
			page:         1,
			pageSize:     10,
			expected: Metadata{
				CurrentPage:  1,
				PageSize:     10,
				FirstPage:    1,
				LastPage:     3,
				TotalRecords: 25,
			},
		},
		{
			name:         "multiple pages - middle page",
			totalRecords: 100,
			page:         5,
			pageSize:     10,
			expected: Metadata{
				CurrentPage:  5,
				PageSize:     10,
				FirstPage:    1,
				LastPage:     10,
				TotalRecords: 100,
			},
		},
		{
			name:         "exact page boundary",
			totalRecords: 100,
			page:         1,
			pageSize:     10,
			expected: Metadata{
				CurrentPage:  1,
				PageSize:     10,
				FirstPage:    1,
				LastPage:     10,
				TotalRecords: 100,
			},
		},
		{
			name:         "partial last page",
			totalRecords: 23,
			page:         1,
			pageSize:     10,
			expected: Metadata{
				CurrentPage:  1,
				PageSize:     10,
				FirstPage:    1,
				LastPage:     3,
				TotalRecords: 23,
			},
		},
		{
			name:         "large dataset",
			totalRecords: 1000,
			page:         20,
			pageSize:     25,
			expected: Metadata{
				CurrentPage:  20,
				PageSize:     25,
				FirstPage:    1,
				LastPage:     40,
				TotalRecords: 1000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateMetadata(tt.totalRecords, tt.page, tt.pageSize)

			if got.CurrentPage != tt.expected.CurrentPage {
				t.Errorf("CurrentPage = %d, want %d", got.CurrentPage, tt.expected.CurrentPage)
			}
			if got.PageSize != tt.expected.PageSize {
				t.Errorf("PageSize = %d, want %d", got.PageSize, tt.expected.PageSize)
			}
			if got.FirstPage != tt.expected.FirstPage {
				t.Errorf("FirstPage = %d, want %d", got.FirstPage, tt.expected.FirstPage)
			}
			if got.LastPage != tt.expected.LastPage {
				t.Errorf("LastPage = %d, want %d", got.LastPage, tt.expected.LastPage)
			}
			if got.TotalRecords != tt.expected.TotalRecords {
				t.Errorf("TotalRecords = %d, want %d", got.TotalRecords, tt.expected.TotalRecords)
			}
		})
	}
}
