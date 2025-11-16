package data

import (
	"github.com/Syha-01/animeVerseAPI/internal/validator"
)

type Filters struct {
	Page     int // which page number the client wants
	PageSize int // how many records per page
}

// Metadata struct holds pagination metadata.
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// ValidateFilters checks the page and page_size parameters.
func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 1000, "page", "must be a maximum of 1000") // Adjusted for larger potential datasets
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
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
