package data

import "database/sql"

type Models struct {
	Animes        AnimeModel
	UserAnimeList UserAnimeListModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Animes:        AnimeModel{DB: db},
		UserAnimeList: UserAnimeListModel{DB: db},
	}
}
