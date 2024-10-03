package main

import (
	"encoding/json"
	"fmt"
	ctrl "github.com/JMURv/effectiveMobile/internal/ctrl"
	"github.com/JMURv/effectiveMobile/internal/ctrl/external"
	hdl "github.com/JMURv/effectiveMobile/internal/hdl/http"
	db "github.com/JMURv/effectiveMobile/internal/repo/db"
	cfg "github.com/JMURv/effectiveMobile/pkg/config"
	"github.com/JMURv/effectiveMobile/pkg/model"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func mustRegisterLogger(mode string) {
	switch mode {
	case "prod":
		zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	case "dev":
		zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Panic("panic occurred", zap.Any("error", err))
			os.Exit(1)
		}
	}()

	conf := cfg.MustLoad()
	mustRegisterLogger(conf.Server.Mode)

	// Setting up main app
	repo := db.New(conf.DB)
	svc := ctrl.New(repo, external.New(conf.ExternalAPIPort))
	h := hdl.New(svc)

	// Start external API
	go startExternalAPI(conf.ExternalAPIPort)

	// Graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-c

		zap.L().Info("Shutting down gracefully...")

		if err := repo.Close(); err != nil {
			zap.L().Debug("Error closing repository", zap.Error(err))
		}
		if err := h.Close(); err != nil {
			zap.L().Debug("Error closing handler", zap.Error(err))
		}

		os.Exit(0)
	}()

	// Start service
	zap.L().Info(
		fmt.Sprintf("Starting server on %v://%v:%v", conf.Server.Scheme, conf.Server.Domain, conf.Server.Port),
	)

	h.Start(conf.Server.Port)
}
func startExternalAPI(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("group") == "" || r.URL.Query().Get("song") == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(&model.SongDetail{
				ReleaseDate: "16.07.2006",
				Link:        "https://example.com",
				Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
			}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	srv := &http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf(":%v", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  20 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		zap.L().Debug("Server error", zap.Error(err))
	}
}
