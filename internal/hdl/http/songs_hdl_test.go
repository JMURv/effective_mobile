package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/JMURv/effectiveMobile/internal/ctrl"
	"github.com/JMURv/effectiveMobile/mocks"
	"github.com/JMURv/effectiveMobile/pkg/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_ListSongs(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	ctrlRepo := mocks.NewMockCtrl(ctrlMock)
	hdl := New(ctrlRepo)

	ctx := context.Background()
	page := 1
	size := 40
	filters := map[string]any{"group": "rock"}

	success := &model.PaginatedSongs{
		Data: []*model.Song{
			{ID: 1, Group: "rock", Song: "song1"},
			{ID: 2, Group: "rock", Song: "song2"},
		},
		Count:       2,
		TotalPages:  1,
		CurrentPage: 1,
		HasNextPage: false,
	}

	t.Run("Success", func(t *testing.T) {
		ctrlRepo.EXPECT().ListSongs(ctx, page, size, filters).Return(success, nil).Times(1)

		req := httptest.NewRequest(http.MethodGet, "/api/songs?page=1&size=40&group=rock", nil)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.ListSongs(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("ErrInternalError", func(t *testing.T) {
		var ErrOther = errors.New("internal error")
		ctrlRepo.EXPECT().ListSongs(ctx, page, size, filters).Return(nil, ErrOther).Times(1)

		req := httptest.NewRequest(http.MethodGet, "/api/songs?page=1&size=40&group=rock", nil)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.ListSongs(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("MissingPageAndSize", func(t *testing.T) {
		ctrlRepo.EXPECT().ListSongs(ctx, 1, 40, filters).Return(success, nil).Times(1)

		req := httptest.NewRequest(http.MethodGet, "/api/songs?group=rock", nil)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.ListSongs(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("InvalidPage", func(t *testing.T) {
		ctrlRepo.EXPECT().ListSongs(ctx, 1, 40, filters).Return(success, nil).Times(1)
		req := httptest.NewRequest(http.MethodGet, "/api/songs?page=invalid&size=40&group=rock", nil) // Use default value for page
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.ListSongs(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("InvalidSize", func(t *testing.T) {
		ctrlRepo.EXPECT().ListSongs(ctx, 1, 40, filters).Return(success, nil).Times(1)
		req := httptest.NewRequest(http.MethodGet, "/api/songs?page=1&size=invalid&group=rock", nil) // Use default value for size
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.ListSongs(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})
}

func TestHandler_GetSong(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	ctrlRepo := mocks.NewMockCtrl(ctrlMock)
	hdl := New(ctrlRepo)

	ctx := context.Background()
	songID := uint64(1)
	page := 1
	size := 40
	songDetails := &model.PaginatedSongs{}

	t.Run("Success", func(t *testing.T) {
		ctrlRepo.EXPECT().GetSong(ctx, songID, page, size).Return(songDetails, nil).Times(1)

		req := httptest.NewRequest(http.MethodGet, "/api/songs/1?page=1&size=40", nil)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.GetSong(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		ctrlRepo.EXPECT().GetSong(ctx, songID, page, size).Return(nil, ctrl.ErrNotFound).Times(1)

		req := httptest.NewRequest(http.MethodGet, "/api/songs/1?page=1&size=40", nil)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.GetSong(w, req)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})

	t.Run("ErrInternalError", func(t *testing.T) {
		var ErrOther = errors.New("other error")
		ctrlRepo.EXPECT().GetSong(ctx, songID, page, size).Return(nil, ErrOther).Times(1)

		req := httptest.NewRequest(http.MethodGet, "/api/songs/1?page=1&size=40", nil)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.GetSong(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("InvalidSongID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/songs/invalid?page=1&size=40", nil)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.GetSong(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("MissingPageAndSize", func(t *testing.T) {
		ctrlRepo.EXPECT().GetSong(ctx, songID, 1, 40).Return(songDetails, nil).Times(1)

		req := httptest.NewRequest(http.MethodGet, "/api/songs/1", nil)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.GetSong(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})
}

func TestHandler_CreateSong(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	ctrlRepo := mocks.NewMockCtrl(ctrlMock)
	hdl := New(ctrlRepo)

	ctx := context.Background()
	success := &model.Song{
		Group: "group",
		Song:  "song",
	}

	t.Run("Success", func(t *testing.T) {
		ctrlRepo.EXPECT().CreateSong(ctx, success).Return(uint64(1), nil).Times(1)

		payload, _ := json.Marshal(success)
		req := httptest.NewRequest(http.MethodPost, "/api/songs", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.CreateSong(w, req)
		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
	})

	t.Run("ErrAlreadyExists", func(t *testing.T) {
		ctrlRepo.EXPECT().CreateSong(ctx, success).Return(uint64(0), ctrl.ErrAlreadyExists).Times(1)

		payload, _ := json.Marshal(success)
		req := httptest.NewRequest(http.MethodPost, "/api/songs", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.CreateSong(w, req)
		assert.Equal(t, http.StatusConflict, w.Result().StatusCode)
	})

	t.Run("ErrInternalError", func(t *testing.T) {
		var ErrOther = errors.New("other error")
		ctrlRepo.EXPECT().CreateSong(ctx, success).Return(uint64(0), ErrOther).Times(1)

		payload, _ := json.Marshal(success)
		req := httptest.NewRequest(http.MethodPost, "/api/songs", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.CreateSong(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("ErrDecodeRequest", func(t *testing.T) {
		payload, _ := json.Marshal(map[string]string{"group": "test-group"})
		req := httptest.NewRequest(http.MethodPost, "/api/songs", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.CreateSong(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		payload, _ := json.Marshal(map[string]any{"group": 123})
		req := httptest.NewRequest(http.MethodPost, "/api/songs", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.CreateSong(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})
}

func TestHandler_UpdateSong(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	ctrlRepo := mocks.NewMockCtrl(ctrlMock)
	hdl := New(ctrlRepo)

	ctx := context.Background()
	success := &model.Song{
		ID:    1,
		Group: "group",
		Song:  "song",
	}

	t.Run("Success", func(t *testing.T) {
		ctrlRepo.EXPECT().UpdateSong(ctx, success).Return(nil).Times(1)

		payload, _ := json.Marshal(success)
		req := httptest.NewRequest(http.MethodPut, "/api/songs/1", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.UpdateSong(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		ctrlRepo.EXPECT().UpdateSong(ctx, success).Return(ctrl.ErrNotFound).Times(1)

		payload, _ := json.Marshal(success)
		req := httptest.NewRequest(http.MethodPut, "/api/songs/1", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.UpdateSong(w, req)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})

	t.Run("ErrInternalError", func(t *testing.T) {
		var ErrOther = errors.New("other error")
		ctrlRepo.EXPECT().UpdateSong(ctx, success).Return(ErrOther).Times(1)

		payload, _ := json.Marshal(success)
		req := httptest.NewRequest(http.MethodPut, "/api/songs/1", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.UpdateSong(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("ErrDecodeRequest", func(t *testing.T) {
		payload, _ := json.Marshal(map[string]string{"group": "test-group"})
		req := httptest.NewRequest(http.MethodPut, "/api/songs/1", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.UpdateSong(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("InvalidSongID", func(t *testing.T) {
		payload, _ := json.Marshal(success)
		req := httptest.NewRequest(http.MethodPut, "/api/songs/invalid", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.UpdateSong(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	// Test case: Invalid JSON
	t.Run("InvalidJSON", func(t *testing.T) {
		payload := []byte(`{"group":123}`)
		req := httptest.NewRequest(http.MethodPut, "/api/songs/1", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.UpdateSong(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})
}

func TestHandler_DeleteSong(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	ctrlRepo := mocks.NewMockCtrl(ctrlMock)
	hdl := New(ctrlRepo)

	ctx := context.Background()
	songID := uint64(1)

	t.Run("Success", func(t *testing.T) {
		ctrlRepo.EXPECT().DeleteSong(ctx, songID).Return(nil).Times(1)

		req := httptest.NewRequest(http.MethodDelete, "/api/songs/1", nil)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.DeleteSong(w, req)
		assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		ctrlRepo.EXPECT().DeleteSong(ctx, songID).Return(ctrl.ErrNotFound).Times(1)

		req := httptest.NewRequest(http.MethodDelete, "/api/songs/1", nil)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.DeleteSong(w, req)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})

	t.Run("ErrInternalError", func(t *testing.T) {
		var ErrOther = errors.New("other error")
		ctrlRepo.EXPECT().DeleteSong(ctx, songID).Return(ErrOther).Times(1)

		req := httptest.NewRequest(http.MethodDelete, "/api/songs/1", nil)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.DeleteSong(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("InvalidSongID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/songs/invalid", nil)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		hdl.DeleteSong(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})
}

func TestHandler_StartAndClose(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	ctrlRepo := mocks.NewMockCtrl(ctrlMock)
	hdl := New(ctrlRepo)

	go hdl.Start(8080)
	time.Sleep(time.Second)

	err := hdl.Close()
	assert.Nil(t, err)
}
