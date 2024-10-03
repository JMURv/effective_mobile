package http

import (
	"context"
	"fmt"
	_ "github.com/JMURv/effectiveMobile/docs"
	"github.com/JMURv/effectiveMobile/internal/hdl"
	"github.com/JMURv/effectiveMobile/pkg/model"
	utils "github.com/JMURv/effectiveMobile/pkg/utils/http"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Ctrl interface {
	ListSongs(ctx context.Context, page int, size int, filters map[string]any) (*model.PaginatedSongs, error)
	GetSong(ctx context.Context, id uint64, page int, size int) (*model.PaginatedSongs, error)
	CreateSong(ctx context.Context, req *model.Song) (uint64, error)
	UpdateSong(ctx context.Context, req *model.Song) error
	DeleteSong(ctx context.Context, id uint64) error
}

type Handler struct {
	srv  *http.Server
	ctrl Ctrl
}

func New(ctrl Ctrl) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

func (h *Handler) Start(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("/api/health-check", func(w http.ResponseWriter, r *http.Request) {
		utils.SuccessResponse(w, http.StatusOK, "OK")
	})

	mux.HandleFunc("/api/songs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.ListSongs(w, r)
		case http.MethodPost:
			h.CreateSong(w, r)
		default:
			utils.ErrResponse(w, http.StatusMethodNotAllowed, hdl.ErrMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/songs/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetSong(w, r)
		case http.MethodPut:
			h.UpdateSong(w, r)
		case http.MethodDelete:
			h.DeleteSong(w, r)
		default:
			utils.ErrResponse(w, http.StatusMethodNotAllowed, hdl.ErrMethodNotAllowed)
		}
	})

	h.srv = &http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf(":%v", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  20 * time.Second,
	}

	if err := h.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		zap.L().Debug("Server error", zap.Error(err))
	}
}

func (h *Handler) Close() error {
	if err := h.srv.Shutdown(context.Background()); err != nil {
		return err
	}
	return nil
}
