package ctrl

import (
	"context"
	"errors"
	"github.com/JMURv/effectiveMobile/internal/repo"
	"github.com/JMURv/effectiveMobile/mocks"
	"github.com/JMURv/effectiveMobile/pkg/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestController_ListSongs(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	svcRepo := mocks.NewMockSongsRepo(ctrlMock)
	extRepo := mocks.NewMockAPIRepo(ctrlMock)

	ctrl := New(svcRepo, extRepo)
	ctx := context.Background()

	page, size, filters := 1, 10, map[string]any{}

	t.Run("Success", func(t *testing.T) {
		svcRepo.EXPECT().ListSongs(gomock.Any(), page, size, filters).Return(&model.PaginatedSongs{}, nil).Times(1)

		res, err := ctrl.ListSongs(ctx, page, size, filters)
		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("ErrOther", func(t *testing.T) {
		newErr := errors.New("new error")
		svcRepo.EXPECT().ListSongs(gomock.Any(), page, size, filters).Return(nil, newErr).Times(1)

		res, err := ctrl.ListSongs(ctx, page, size, filters)
		assert.IsType(t, newErr, err)
		assert.Nil(t, res)
	})

}

func TestController_GetSong(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	svcRepo := mocks.NewMockSongsRepo(ctrlMock)
	extRepo := mocks.NewMockAPIRepo(ctrlMock)

	ctrl := New(svcRepo, extRepo)
	ctx := context.Background()

	idx, page, size := uint64(1), 1, 40

	t.Run("Success", func(t *testing.T) {
		svcRepo.EXPECT().GetSong(gomock.Any(), idx, page, size).Return(&model.PaginatedSongs{}, nil).Times(1)

		res, err := ctrl.GetSong(ctx, idx, page, size)
		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		svcRepo.EXPECT().GetSong(gomock.Any(), idx, page, size).Return(nil, repo.ErrNotFound).Times(1)

		res, err := ctrl.GetSong(ctx, idx, page, size)
		assert.IsType(t, ErrNotFound, err)
		assert.Nil(t, res)
	})

	t.Run("ErrOther", func(t *testing.T) {
		newErr := errors.New("new error")
		svcRepo.EXPECT().GetSong(gomock.Any(), idx, page, size).Return(nil, newErr).Times(1)

		res, err := ctrl.GetSong(ctx, idx, page, size)
		assert.IsType(t, newErr, err)
		assert.Nil(t, res)
	})
}

func TestController_CreateSong(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	svcRepo := mocks.NewMockSongsRepo(ctrlMock)
	extRepo := mocks.NewMockAPIRepo(ctrlMock)

	ctrl := New(svcRepo, extRepo)
	ctx := context.Background()
	details := &model.SongDetail{
		ReleaseDate: "16.07.2006",
		Text:        "test text\n\ntest text",
		Link:        "https://example.com",
	}

	expCreate := &model.Song{
		Lyrics:      []string{"test text", "test text"},
		Link:        "https://example.com",
		ReleaseDate: time.Date(2006, 7, 16, 0, 0, 0, 0, time.UTC),
	}

	t.Run("Success", func(t *testing.T) {
		extRepo.EXPECT().FetchSongDetail(gomock.Any(), gomock.Any()).Return(details, nil).Times(1)

		svcRepo.EXPECT().CreateSong(gomock.Any(), expCreate).Return(uint64(1), nil).Times(1)

		idx, err := ctrl.CreateSong(ctx, &model.Song{})
		assert.Nil(t, err)
		assert.Equal(t, uint64(1), idx)
	})

	t.Run("ErrAlreadyExists", func(t *testing.T) {
		extRepo.EXPECT().FetchSongDetail(gomock.Any(), gomock.Any()).Return(details, nil).Times(1)
		svcRepo.EXPECT().CreateSong(gomock.Any(), expCreate).Return(uint64(0), repo.ErrAlreadyExists).Times(1)

		idx, err := ctrl.CreateSong(ctx, &model.Song{})
		assert.IsType(t, repo.ErrAlreadyExists, err)
		assert.Equal(t, uint64(0), idx)
	})

	t.Run("ErrOther", func(t *testing.T) {
		newErr := errors.New("new error")
		extRepo.EXPECT().FetchSongDetail(gomock.Any(), gomock.Any()).Return(details, nil).Times(1)
		svcRepo.EXPECT().CreateSong(gomock.Any(), expCreate).Return(uint64(0), newErr).Times(1)

		idx, err := ctrl.CreateSong(ctx, &model.Song{})
		assert.IsType(t, newErr, err)
		assert.Equal(t, uint64(0), idx)
	})

	t.Run("ErrFetchDetails", func(t *testing.T) {
		newErr := errors.New("new error")
		extRepo.EXPECT().FetchSongDetail(gomock.Any(), gomock.Any()).Return(nil, newErr).Times(1)

		idx, err := ctrl.CreateSong(ctx, &model.Song{})
		assert.IsType(t, newErr, err)
		assert.Equal(t, uint64(0), idx)
	})

	t.Run("ErrParseDate", func(t *testing.T) {
		extRepo.EXPECT().FetchSongDetail(gomock.Any(), gomock.Any()).Return(&model.SongDetail{
			ReleaseDate: "bad-time",
			Text:        "test text\n\ntest text",
			Link:        "https://example.com",
		}, nil).Times(1)

		idx, err := ctrl.CreateSong(ctx, &model.Song{})
		assert.NotNil(t, err)
		assert.Equal(t, uint64(0), idx)
	})
}

func TestController_UpdateSong(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	svcRepo := mocks.NewMockSongsRepo(ctrlMock)
	extRepo := mocks.NewMockAPIRepo(ctrlMock)

	ctrl := New(svcRepo, extRepo)
	ctx := context.Background()

	req := &model.Song{
		ID:    1,
		Song:  "updated song",
		Group: "updated group",
	}

	t.Run("Success", func(t *testing.T) {
		svcRepo.EXPECT().UpdateSong(gomock.Any(), req).Return(nil).Times(1)

		err := ctrl.UpdateSong(ctx, req)
		assert.Nil(t, err)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		svcRepo.EXPECT().UpdateSong(gomock.Any(), req).Return(repo.ErrNotFound).Times(1)

		err := ctrl.UpdateSong(ctx, req)
		assert.IsType(t, ErrNotFound, err)
	})

	t.Run("ErrOther", func(t *testing.T) {
		newErr := errors.New("new error")
		svcRepo.EXPECT().UpdateSong(gomock.Any(), req).Return(newErr).Times(1)

		err := ctrl.UpdateSong(ctx, req)
		assert.IsType(t, newErr, err)
	})
}

func TestController_DeleteSong(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	svcRepo := mocks.NewMockSongsRepo(ctrlMock)
	extRepo := mocks.NewMockAPIRepo(ctrlMock)

	ctrl := New(svcRepo, extRepo)
	ctx := context.Background()
	idx := uint64(1)

	t.Run("Success", func(t *testing.T) {
		svcRepo.EXPECT().DeleteSong(gomock.Any(), idx).Return(nil).Times(1)

		err := ctrl.DeleteSong(ctx, idx)
		assert.Nil(t, err)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		svcRepo.EXPECT().DeleteSong(gomock.Any(), idx).Return(repo.ErrNotFound).Times(1)

		err := ctrl.DeleteSong(ctx, idx)
		assert.IsType(t, ErrNotFound, err)
	})

	t.Run("ErrOther", func(t *testing.T) {
		newErr := errors.New("new error")
		svcRepo.EXPECT().DeleteSong(gomock.Any(), idx).Return(newErr).Times(1)

		err := ctrl.DeleteSong(ctx, idx)
		assert.IsType(t, newErr, err)
	})

}
