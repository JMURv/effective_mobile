package ctrl

import (
	"context"
	"errors"
	"github.com/JMURv/effectiveMobile/internal/repo"
	"github.com/JMURv/effectiveMobile/pkg/model"
	"go.uber.org/zap"
	"strings"
	"time"
)

type SongsRepo interface {
	ListSongs(ctx context.Context, page, size int, filters map[string]any) (*model.PaginatedSongs, error)
	GetSong(ctx context.Context, id uint64, page, size int) (*model.PaginatedSongs, error)
	CreateSong(ctx context.Context, req *model.Song) (uint64, error)
	UpdateSong(ctx context.Context, req *model.Song) error
	DeleteSong(ctx context.Context, id uint64) error
}

type APIRepo interface {
	FetchSongDetail(group, song string) (*model.SongDetail, error)
}

type Controller struct {
	repo SongsRepo
	api  APIRepo
}

func New(repo SongsRepo, api APIRepo) *Controller {
	return &Controller{
		repo: repo,
		api:  api,
	}
}

func (c *Controller) ListSongs(ctx context.Context, page, size int, filters map[string]any) (*model.PaginatedSongs, error) {
	const op = "songs.ListSongs.ctrl"

	res, err := c.repo.ListSongs(ctx, page, size, filters)
	if err != nil {
		zap.L().Debug(
			"failed to list songs",
			zap.Error(err), zap.String("op", op),
			zap.Int("page", page), zap.Int("size", size),
		)
		return nil, err
	}

	return res, nil
}

func (c *Controller) GetSong(ctx context.Context, id uint64, page, size int) (*model.PaginatedSongs, error) {
	const op = "songs.GetSong.ctrl"

	res, err := c.repo.GetSong(ctx, id, page, size)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug(
			"failed to find song",
			zap.Error(err), zap.String("op", op),
			zap.Uint64("ID", id), zap.Int("page", page), zap.Int("size", size),
		)
		return nil, ErrNotFound
	} else if err != nil {
		zap.L().Debug(
			"failed to get song",
			zap.Error(err), zap.String("op", op),
			zap.Uint64("ID", id), zap.Int("page", page), zap.Int("size", size),
		)
		return nil, err
	}

	return res, nil
}

func (c *Controller) CreateSong(ctx context.Context, req *model.Song) (uint64, error) {
	const op = "songs.CreateSong.ctrl"

	details, err := c.api.FetchSongDetail(req.Group, req.Song)
	if err != nil {
		zap.L().Debug(
			"failed to fetch song details",
			zap.Error(err), zap.String("op", op),
		)
		return 0, err
	}

	parsedDate, err := time.Parse("02.01.2006", details.ReleaseDate)
	if err != nil {
		zap.L().Debug(
			"failed to parse date",
			zap.Error(err), zap.String("op", op),
		)
		return 0, err
	}
	req.ReleaseDate = parsedDate
	req.Lyrics = strings.Split(details.Text, "\n\n")
	req.Link = details.Link

	res, err := c.repo.CreateSong(ctx, req)
	if err != nil && errors.Is(err, repo.ErrAlreadyExists) {
		zap.L().Debug(
			"song already exists",
			zap.Error(err), zap.String("op", op),
		)
		return 0, ErrAlreadyExists
	} else if err != nil {
		zap.L().Debug(
			"failed to create song",
			zap.Error(err), zap.String("op", op),
		)
		return 0, err
	}

	return res, nil
}

func (c *Controller) UpdateSong(ctx context.Context, req *model.Song) error {
	const op = "songs.UpdateSong.ctrl"

	err := c.repo.UpdateSong(ctx, req)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug(
			"failed to find song",
			zap.Error(err), zap.String("op", op),
			zap.Uint64("ID", req.ID), zap.String("song", req.Song),
		)
		return ErrNotFound
	} else if err != nil {
		zap.L().Debug(
			"failed to update song",
			zap.Error(err), zap.String("op", op),
			zap.Uint64("ID", req.ID), zap.String("song", req.Song),
		)
		return err
	}

	return nil
}

func (c *Controller) DeleteSong(ctx context.Context, id uint64) error {
	const op = "songs.DeleteSong.ctrl"

	err := c.repo.DeleteSong(ctx, id)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug(
			"failed to find song",
			zap.Error(err), zap.String("op", op),
			zap.Uint64("ID", id),
		)
		return ErrNotFound
	} else if err != nil {
		zap.L().Debug(
			"failed to delete song",
			zap.Error(err), zap.String("op", op),
			zap.Uint64("ID", id),
		)
		return err
	}

	return nil
}
