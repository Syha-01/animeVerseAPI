package data

import "database/sql"

type Models struct {
	Animes AnimeModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Animes: AnimeModel{DB: db},
	}
}
