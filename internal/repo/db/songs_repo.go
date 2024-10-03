package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/JMURv/effectiveMobile/internal/repo"
	"github.com/JMURv/effectiveMobile/pkg/model"
	utils "github.com/JMURv/effectiveMobile/pkg/utils/db"
	"github.com/lib/pq"
	"strings"
)

func (r *Repository) ListSongs(ctx context.Context, page, size int, filters map[string]any) (*model.PaginatedSongs, error) {
	filterQ, args := utils.BuildFilterQuery(filters)

	var selectQ strings.Builder
	selectQ.WriteString(`
		SELECT group_name, song_name, release_date, link, lyrics 
		FROM songs
	`)
	selectQ.WriteString(filterQ)
	selectQ.WriteString(fmt.Sprintf(" LIMIT %v OFFSET %v", size, (page-1)*size))

	rows, err := r.conn.QueryContext(ctx, selectQ.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*model.Song, 0, size)
	for rows.Next() {
		song := &model.Song{}
		if err := rows.Scan(&song.Group, &song.Song, &song.ReleaseDate, &song.Link, pq.Array(&song.Lyrics)); err != nil {
			return nil, err
		}
		res = append(res, song)
	}

	var countQ strings.Builder
	countQ.WriteString("SELECT COUNT(*) FROM songs")
	countQ.WriteString(filterQ)

	var count int64
	if err := r.conn.QueryRowContext(ctx, countQ.String(), args...).Scan(&count); err != nil {
		return nil, err
	}

	totalPages := int((count + int64(size) - 1) / int64(size))
	return &model.PaginatedSongs{
		Data:        res,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: page < totalPages,
	}, nil
}

func (r *Repository) GetSong(ctx context.Context, id uint64, page, size int) (*model.PaginatedSongs, error) {
	var count int64
	res := &model.Song{}

	offset := (page - 1) * size
	err := r.conn.QueryRowContext(ctx, `
		SELECT group_name, song_name, release_date, link, lyrics[$2:$3], array_length(lyrics, 1) as count
		FROM songs
		WHERE id = $1
		`, id, offset+1, offset+size).
		Scan(&res.Group, &res.Song, &res.ReleaseDate, &res.Link, pq.Array(&res.Lyrics), &count)

	if err == sql.ErrNoRows {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	totalPages := int((count + int64(size) - 1) / int64(size))
	return &model.PaginatedSongs{
		Data:        res,
		Count:       count,
		TotalPages:  totalPages,
		CurrentPage: page,
		HasNextPage: page < totalPages,
	}, nil
}

func (r *Repository) CreateSong(ctx context.Context, req *model.Song) (uint64, error) {
	var idx uint64
	err := r.conn.QueryRow(`SELECT id FROM songs WHERE group_name=$1 AND song_name=$2`, req.Group, req.Song).Scan(&idx)
	if err == nil {
		return 0, repo.ErrAlreadyExists
	} else if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	var id uint64
	err = r.conn.QueryRowContext(
		ctx,
		`INSERT INTO songs (group_name, song_name, release_date, lyrics, link) 
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		req.Group, req.Song, req.ReleaseDate, pq.Array(req.Lyrics), req.Link,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) UpdateSong(ctx context.Context, req *model.Song) error {
	var idx uint64
	err := r.conn.
		QueryRow(`SELECT id FROM songs WHERE id = $1`, req.ID).
		Scan(&idx)

	if err == sql.ErrNoRows {
		return repo.ErrNotFound
	} else if err != nil {
		return err
	}

	if _, err = r.conn.ExecContext(ctx,
		`UPDATE songs SET group_name = $1, song_name = $2, release_date = $3, lyrics = $4, link = $5 WHERE id = $6`,
		req.Group, req.Song, req.ReleaseDate, pq.Array(req.Lyrics), req.Link, req.ID,
	); err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteSong(ctx context.Context, id uint64) error {
	var idx uint64
	err := r.conn.
		QueryRow(`SELECT id FROM songs WHERE id = $1`, id).
		Scan(&idx)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return repo.ErrNotFound
	} else if err != nil {
		return err
	}

	if _, err = r.conn.ExecContext(ctx, `DELETE FROM songs WHERE id = $1`, id); err != nil {
		return err
	}
	return nil
}
