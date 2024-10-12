package main

import (
	"astrologerService/api"
	"astrologerService/config"
	"astrologerService/migrations"
	"astrologerService/storage"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	ctx := context.Background()

	conn, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		logger.Error("connect to database", "error", err)
		return
	}

	logger.Info("connected to postgres -ok")
	logger.Info("ping Postgres...")

	err = conn.Ping(ctx)
	if err != nil {
		logger.Error("ping to database", "error", err)
		return
	}
	defer conn.Close()

	logger.Info("up migrations...")

	err = migrations.UpMigrations(conn)
	if err != nil {
		log.Panic(err)
	}

	st := storage.NewStorage(conn)
	handl := api.NewHandler(st, logger)
	mw := api.NewMiddleware(logger)

	ctx, _ = context.WithCancel(ctx)

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err = handl.Information(ctx)
				if err != nil {
					logger.Error("get information", "error", err)
				}
			}
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/apod", handl.Apod)
	mux.HandleFunc("GET /api/apod/date", handl.ApodByDate)
	server := http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           mw.Logging(mux),
		ReadTimeout:       time.Second,
		ReadHeaderTimeout: time.Second,
		WriteTimeout:      time.Second,
	}
	err = server.ListenAndServe()
	if err != nil {
		logger.Error("listen and serve", "error", err)
		return
	}

	logger.Info("server started!", "port", cfg.Port)
}
