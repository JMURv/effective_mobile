package external

import (
	"encoding/json"
	errs "github.com/JMURv/effectiveMobile/internal/ctrl"
	"github.com/JMURv/effectiveMobile/pkg/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestController_FetchSongDetail(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/info", r.URL.Path)
		assert.Equal(t, "test-group", r.URL.Query().Get("group"))
		assert.Equal(t, "test-song", r.URL.Query().Get("song"))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&model.SongDetail{
			ReleaseDate: "16.07.2006",
			Text:        "Ooh baby, don't you know I suffer?",
			Link:        "https://example.com",
		})
	}))
	defer server.Close()

	s, _ := strconv.Atoi(strings.Split(server.URL, ":")[2])
	ctrl := New(s)

	t.Run("Success", func(t *testing.T) {
		result, err := ctrl.FetchSongDetail("test-group", "test-song")
		assert.Nil(t, err)

		expectedResponse := &model.SongDetail{
			ReleaseDate: "16.07.2006",
			Text:        "Ooh baby, don't you know I suffer?",
			Link:        "https://example.com",
		}

		assert.Equal(t, expectedResponse, result)
	})

	t.Run("StatusBadRequest", func(t *testing.T) {
		server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		})

		result, err := ctrl.FetchSongDetail("test-group", "test-song")
		assert.NotNil(t, err)
		assert.Equal(t, errs.ErrBadExtReq, err)
		assert.Nil(t, result)
	})

	t.Run("StatusInternalServerError", func(t *testing.T) {
		server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})

		result, err := ctrl.FetchSongDetail("test-group", "test-song")
		assert.NotNil(t, err)
		assert.Equal(t, errs.ExtSrvErr, err)
		assert.Nil(t, result)
	})

	t.Run("ErrorDecoding", func(t *testing.T) {
		server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid json"))
		})

		result, err := ctrl.FetchSongDetail("test-group", "test-song")
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("UnexpectedStatusCode", func(t *testing.T) {
		server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		})

		result, err := ctrl.FetchSongDetail("test-group", "test-song")
		assert.NotNil(t, err)
		assert.Equal(t, errs.ErrExtUnreachable, err)
		assert.Nil(t, result)
	})
}
