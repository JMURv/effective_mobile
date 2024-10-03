package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/JMURv/effectiveMobile/internal/repo"
	"github.com/JMURv/effectiveMobile/pkg/model"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"regexp"
	"testing"
	"time"
)

// https://github.com/DATA-DOG/go-sqlmock/issues/201
func TestRepository_ListSongs(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repository := Repository{conn: db}

	//t.Run("Success", func(t *testing.T) {
	//	page := 1
	//	size := 2
	//	filters := map[string]any{"group": "test-group", "min_release_date": "2006-06-15"}
	//	lyrics := []string{"Lyric 1", "Lyric 2"}
	//	count := int64(3)
	//
	//	expectedArgs := []interface{}{"%test-group%", "2006-06-15"}
	//
	//	mock.ExpectQuery(regexp.QuoteMeta(`SELECT group_name, song_name, release_date, link, lyrics FROM songs WHERE group_name ILIKE $1 AND release_date >= $2 LIMIT 2 OFFSET 0`)).
	//		WithArgs(expectedArgs[0], expectedArgs[1]).
	//		WillReturnRows(sqlmock.NewRows([]string{"group_name", "song_name", "release_date", "link", "lyrics"}).
	//			AddRow("test-group", "test-song", time.Now(), "https://example.com", pq.Array(lyrics)).
	//			AddRow("test-group", "test-song 2", time.Now(), "https://example.com/2", pq.Array(lyrics)))
	//
	//	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM songs WHERE group_name ILIKE $1 AND release_date >= $2`)).
	//		WithArgs(expectedArgs[0], expectedArgs[1]).
	//		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
	//
	//	result, err := repository.ListSongs(context.Background(), page, size, filters)
	//	require.NoError(t, err)
	//	require.NotNil(t, result)
	//	assert.Equal(t, count, result.Count)
	//	assert.Equal(t, 2, len(result.Data.([]any)))
	//	assert.Equal(t, 2, result.TotalPages)
	//	assert.Equal(t, page, result.CurrentPage)
	//	assert.Equal(t, true, result.HasNextPage)
	//
	//	require.NoError(t, mock.ExpectationsWereMet())
	//})

	t.Run("DBErrorOnQuery", func(t *testing.T) {
		page := 1
		size := 2
		filters := map[string]any{"group": "test-group"}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT group_name, song_name, release_date, link, lyrics FROM songs`)).
			WillReturnError(errors.New("some database error"))

		result, err := repository.ListSongs(context.Background(), page, size, filters)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "some database error", err.Error())

		require.NoError(t, mock.ExpectationsWereMet())
	})

	//t.Run("DBErrorOnCount", func(t *testing.T) {
	//	page := 1
	//	size := 2
	//	filters := map[string]any{"group": "test-group"}
	//	lyrics := []string{"Lyric 1", "Lyric 2"}
	//
	//	mock.ExpectQuery(regexp.QuoteMeta(`SELECT group_name, song_name, release_date, link, lyrics FROM songs`)).
	//		WillReturnRows(sqlmock.NewRows([]string{"group_name", "song_name", "release_date", "link", "lyrics"}).
	//			AddRow("test-group", "test-song", "2006-07-16", "https://example.com", pq.Array(lyrics)))
	//
	//	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM songs`)).
	//		WillReturnError(errors.New("some count error"))
	//
	//	result, err := repository.ListSongs(context.Background(), page, size, filters)
	//	require.Error(t, err)
	//	assert.Nil(t, result)
	//	assert.Equal(t, "some count error", err.Error())
	//
	//	require.NoError(t, mock.ExpectationsWereMet())
	//})
}

// https://github.com/DATA-DOG/go-sqlmock/issues/201
func TestRepository_GetSong(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repository := Repository{conn: db}

	//t.Run("Success", func(t *testing.T) {
	//	id := uint64(1)
	//	page := 1
	//	size := 2
	//	lyrics := []string{"Lyric 1", "Lyric 2"}
	//	count := int64(len(lyrics))
	//
	//	releaseDate, err := time.Parse("2006-01-02", "2006-07-16")
	//	require.NoError(t, err)
	//
	//	mock.ExpectQuery(regexp.QuoteMeta(`SELECT group_name, song_name, release_date, link, lyrics[$2:$3], array_length(lyrics, 1) as count FROM songs WHERE id = $1`)).
	//		WithArgs(id, 1, size).
	//		WillReturnRows(sqlmock.NewRows([]string{"group_name", "song_name", "release_date", "link", "lyrics", "count"}).
	//			AddRow("test-group", "test-song", releaseDate, "https://example.com", lyrics, count))
	//
	//	result, err := repository.GetSong(context.Background(), id, page, size)
	//	require.NoError(t, err)
	//	require.NotNil(t, result)
	//	assert.Equal(t, count, result.Count)
	//	assert.Equal(t, 1, result.TotalPages)
	//	assert.Equal(t, page, result.CurrentPage)
	//	assert.Equal(t, false, result.HasNextPage)
	//
	//	require.NoError(t, mock.ExpectationsWereMet())
	//})

	t.Run("ErrNotFound", func(t *testing.T) {
		id := uint64(2)
		page := 1
		size := 2

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT group_name, song_name, release_date, link, lyrics[$2:$3], array_length(lyrics, 1) as count FROM songs WHERE id = $1`)).
			WithArgs(id, 1, size).
			WillReturnError(sql.ErrNoRows)

		result, err := repository.GetSong(context.Background(), id, page, size)
		require.Error(t, err)
		assert.Equal(t, repo.ErrNotFound, err)
		assert.Nil(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DBErrorOnQuery", func(t *testing.T) {
		id := uint64(3)
		page := 1
		size := 2

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT group_name, song_name, release_date, link, lyrics[$2:$3], array_length(lyrics, 1) as count FROM songs WHERE id = $1`)).
			WithArgs(id, 1, size).
			WillReturnError(errors.New("some database error"))

		result, err := repository.GetSong(context.Background(), id, page, size)
		require.Error(t, err)
		assert.Equal(t, "some database error", err.Error())
		assert.Nil(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_CreateSong(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repository := Repository{conn: db}

	t.Run("Success", func(t *testing.T) {
		req := &model.Song{
			Group:       "test-group",
			Song:        "test-song",
			ReleaseDate: time.Now(),
			Lyrics:      []string{"Lyric 1", "Lyric 2"},
			Link:        "https://example.com",
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM songs WHERE group_name=$1 AND song_name=$2`)).
			WithArgs(req.Group, req.Song).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO songs (group_name, song_name, release_date, lyrics, link) VALUES ($1, $2, $3, $4, $5) RETURNING id`)).
			WithArgs(req.Group, req.Song, req.ReleaseDate, pq.Array(req.Lyrics), req.Link).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		id, err := repository.CreateSong(context.Background(), req)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), id)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ErrAlreadyExists", func(t *testing.T) {
		req := &model.Song{
			Group:       "test-group",
			Song:        "test-song",
			ReleaseDate: time.Now(),
			Lyrics:      []string{"Lyric 1", "Lyric 2"},
			Link:        "https://example.com",
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM songs WHERE group_name=$1 AND song_name=$2`)).
			WithArgs(req.Group, req.Song).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		id, err := repository.CreateSong(context.Background(), req)
		require.Error(t, err)
		assert.Equal(t, repo.ErrAlreadyExists, err)
		assert.Equal(t, uint64(0), id)
	})

	t.Run("DBErrorOnSelect", func(t *testing.T) {
		req := &model.Song{
			Group:       "test-group",
			Song:        "test-song",
			ReleaseDate: time.Now(),
			Lyrics:      []string{"Lyric 1", "Lyric 2"},
			Link:        "https://example.com",
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM songs WHERE group_name=$1 AND song_name=$2`)).
			WithArgs(req.Group, req.Song).
			WillReturnError(errors.New("some database error"))

		id, err := repository.CreateSong(context.Background(), req)
		require.Error(t, err)
		assert.Equal(t, "some database error", err.Error())
		assert.Equal(t, uint64(0), id)
	})

	t.Run("DBErrorOnInsert", func(t *testing.T) {
		req := &model.Song{
			Group:       "test-group",
			Song:        "test-song",
			ReleaseDate: time.Now(),
			Lyrics:      []string{"Lyric 1", "Lyric 2"},
			Link:        "https://example.com",
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM songs WHERE group_name=$1 AND song_name=$2`)).
			WithArgs(req.Group, req.Song).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO songs (group_name, song_name, release_date, lyrics, link) VALUES ($1, $2, $3, $4, $5) RETURNING id`)).
			WithArgs(req.Group, req.Song, req.ReleaseDate, pq.Array(req.Lyrics), req.Link).
			WillReturnError(errors.New("some insert error"))

		id, err := repository.CreateSong(context.Background(), req)
		require.Error(t, err)
		assert.Equal(t, "some insert error", err.Error())
		assert.Equal(t, uint64(0), id)
	})

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_UpdateSong(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repository := Repository{conn: db}

	t.Run("Success", func(t *testing.T) {
		req := &model.Song{
			ID:          1,
			Group:       "test-group",
			Song:        "updated-test-song",
			ReleaseDate: time.Now(),
			Lyrics:      []string{"Updated Lyric 1", "Updated Lyric 2"},
			Link:        "https://example.com/updated",
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM songs WHERE id = $1`)).
			WithArgs(req.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(req.ID))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE songs SET group_name = $1, song_name = $2, release_date = $3, lyrics = $4, link = $5 WHERE id = $6`)).
			WithArgs(req.Group, req.Song, req.ReleaseDate, pq.Array(req.Lyrics), req.Link, req.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repository.UpdateSong(context.Background(), req)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		req := &model.Song{
			ID:          2,
			Group:       "test-group",
			Song:        "test-song",
			ReleaseDate: time.Now(),
			Lyrics:      []string{"Lyric 1", "Lyric 2"},
			Link:        "https://example.com",
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM songs WHERE id = $1`)).
			WithArgs(req.ID).
			WillReturnError(sql.ErrNoRows)

		err := repository.UpdateSong(context.Background(), req)
		require.Error(t, err)
		assert.Equal(t, repo.ErrNotFound, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DBErrorOnSelect", func(t *testing.T) {
		req := &model.Song{
			ID:          3,
			Group:       "test-group",
			Song:        "test-song",
			ReleaseDate: time.Now(),
			Lyrics:      []string{"Lyric 1", "Lyric 2"},
			Link:        "https://example.com",
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM songs WHERE id = $1`)).
			WithArgs(req.ID).
			WillReturnError(errors.New("some database error"))

		err := repository.UpdateSong(context.Background(), req)
		require.Error(t, err)
		assert.Equal(t, "some database error", err.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DBErrorOnUpdate", func(t *testing.T) {
		req := &model.Song{
			ID:          4,
			Group:       "test-group",
			Song:        "test-song",
			ReleaseDate: time.Now(),
			Lyrics:      []string{"Lyric 1", "Lyric 2"},
			Link:        "https://example.com",
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM songs WHERE id = $1`)).
			WithArgs(req.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(req.ID))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE songs SET group_name = $1, song_name = $2, release_date = $3, lyrics = $4, link = $5 WHERE id = $6`)).
			WithArgs(req.Group, req.Song, req.ReleaseDate, pq.Array(req.Lyrics), req.Link, req.ID).
			WillReturnError(errors.New("some update error"))

		err := repository.UpdateSong(context.Background(), req)
		require.Error(t, err)
		assert.Equal(t, "some update error", err.Error())

		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_DeleteSong(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repository := Repository{conn: db}

	t.Run("Success", func(t *testing.T) {
		id := uint64(1)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM songs WHERE id = $1`)).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM songs WHERE id = $1`)).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repository.DeleteSong(context.Background(), id)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		id := uint64(2)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM songs WHERE id = $1`)).
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)

		err := repository.DeleteSong(context.Background(), id)
		require.Error(t, err)
		assert.Equal(t, repo.ErrNotFound, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DBErrorOnSelect", func(t *testing.T) {
		id := uint64(3)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM songs WHERE id = $1`)).
			WithArgs(id).
			WillReturnError(errors.New("some database error"))

		err := repository.DeleteSong(context.Background(), id)
		require.Error(t, err)
		assert.Equal(t, "some database error", err.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DBErrorOnDelete", func(t *testing.T) {
		id := uint64(4)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM songs WHERE id = $1`)).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM songs WHERE id = $1`)).
			WithArgs(id).
			WillReturnError(errors.New("some delete error"))

		err := repository.DeleteSong(context.Background(), id)

		require.Error(t, err)
		assert.Equal(t, "some delete error", err.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
