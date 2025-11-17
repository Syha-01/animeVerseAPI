package data

import (
	"strings"

	"github.com/Syha-01/animeVerseAPI/internal/validator"
)

// The Filters struct contains fields for pagination and sorting.
type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string // A list of allowed sort fields
}

// Metadata struct holds pagination metadata.
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// ValidateFilters checks the page, page_size, and sort parameters.
func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 1000, "page", "must be a maximum of 1000")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

	// Check if the sort field is a valid value from the safelist.
	v.Check(validator.PermittedValue(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}

// sortColumn checks the client-provided sort field against the safelist
// and returns the column name. It panics if the value is unsafe to prevent SQL injection.
func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

// sortDirection returns the sort direction ("ASC" or "DESC") based on the "-" prefix.
func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

// limit calculates how many records to send back.
func (f Filters) limit() int {
	return f.PageSize
}

// offset calculates the starting point for the records to be sent.
func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

// calculateMetadata computes the metadata values.
func calculateMetadata(totalRecords int, page int, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}
	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     (totalRecords + pageSize - 1) / pageSize,
		TotalRecords: totalRecords,
	}
}
