package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Syha-01/animeVerseAPI/internal/validator"
)

type UserAnimeList struct {
	ID                   string    `json:"id"`
	UserID               string    `json:"user_id"`
	AnimeID              int64     `json:"anime_id"`
	Status               string    `json:"status"`
	CurrentEpisode       int32     `json:"current_episode"`
	Score                int32     `json:"score,omitempty"`
	StartedWatchingDate  time.Time `json:"started_watching_date,omitempty"`
	FinishedWatchingDate time.Time `json:"finished_watching_date,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func ValidateUserAnimeList(v *validator.Validator, list *UserAnimeList) {
	v.Check(list.UserID != "", "user_id", "must be provided")
	v.Check(list.AnimeID > 0, "anime_id", "must be a positive integer")
	v.Check(list.Status != "", "status", "must be provided")
	v.Check(validator.PermittedValue(list.Status, "Watching", "Completed", "Dropped", "Watch later"), "status", "must be a valid status")
	v.Check(list.CurrentEpisode >= 0, "current_episode", "must be a positive integer")
	if list.Score != 0 {
		v.Check(list.Score >= 1 && list.Score <= 10, "score", "must be between 1 and 10")
	}
}

type UserAnimeListModel struct {
	DB *sql.DB
}

func (m UserAnimeListModel) InsertUserAnimeList(list *UserAnimeList) error {
	query := `
		INSERT INTO user_anime_list (user_id, anime_id, status, current_episode, score, started_watching_date, finished_watching_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	args := []any{
		list.UserID,
		list.AnimeID,
		list.Status,
		list.CurrentEpisode,
		list.Score,
		list.StartedWatchingDate,
		list.FinishedWatchingDate,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&list.ID, &list.CreatedAt, &list.UpdatedAt)
}

func (m UserAnimeListModel) GetUserAnimeList(id string) (*UserAnimeList, error) {
	query := `
		SELECT id, user_id, anime_id, status, current_episode, score, started_watching_date, finished_watching_date, created_at, updated_at
		FROM user_anime_list
		WHERE id = $1`

	var list UserAnimeList
	var score sql.NullInt32
	var startedWatchingDate, finishedWatchingDate sql.NullTime

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&list.ID,
		&list.UserID,
		&list.AnimeID,
		&list.Status,
		&list.CurrentEpisode,
		&score,
		&startedWatchingDate,
		&finishedWatchingDate,
		&list.CreatedAt,
		&list.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	if score.Valid {
		list.Score = score.Int32
	}
	if startedWatchingDate.Valid {
		list.StartedWatchingDate = startedWatchingDate.Time
	}
	if finishedWatchingDate.Valid {
		list.FinishedWatchingDate = finishedWatchingDate.Time
	}

	return &list, nil
}

func (m UserAnimeListModel) UpdateUserAnimeList(list *UserAnimeList) error {
	query := `
		UPDATE user_anime_list
		SET status = $1, current_episode = $2, score = $3, started_watching_date = $4, finished_watching_date = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at`

	args := []any{
		list.Status,
		list.CurrentEpisode,
		list.Score,
		list.StartedWatchingDate,
		list.FinishedWatchingDate,
		list.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&list.UpdatedAt)
}

func (m UserAnimeListModel) DeleteUserAnimeList(id string) error {
	query := `
		DELETE FROM user_anime_list
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m UserAnimeListModel) GetAllUserAnimeLists(userID string, status string, filters Filters) ([]*UserAnimeList, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, user_id, anime_id, status, current_episode, score, started_watching_date, finished_watching_date, created_at, updated_at
		FROM user_anime_list
		WHERE (user_id::text = $1 OR $1 = '')
		AND (status::text = $2 OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{userID, status, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	lists := []*UserAnimeList{}

	for rows.Next() {
		var list UserAnimeList
		var score sql.NullInt32
		var startedWatchingDate, finishedWatchingDate sql.NullTime

		err := rows.Scan(
			&totalRecords,
			&list.ID,
			&list.UserID,
			&list.AnimeID,
			&list.Status,
			&list.CurrentEpisode,
			&score,
			&startedWatchingDate,
			&finishedWatchingDate,
			&list.CreatedAt,
			&list.UpdatedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		if score.Valid {
			list.Score = score.Int32
		}
		if startedWatchingDate.Valid {
			list.StartedWatchingDate = startedWatchingDate.Time
		}
		if finishedWatchingDate.Valid {
			list.FinishedWatchingDate = finishedWatchingDate.Time
		}

		lists = append(lists, &list)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return lists, metadata, nil
}
