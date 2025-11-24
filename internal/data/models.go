package data

import "database/sql"

type Models struct {
	Animes        AnimeModel
	UserAnimeList UserAnimeListModel
	Users         UserModel
	Tokens        TokenModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Animes:        AnimeModel{DB: db},
		UserAnimeList: UserAnimeListModel{DB: db},
		Users:         UserModel{DB: db},
		Tokens:        TokenModel{DB: db},
	}
}
