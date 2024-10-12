package api

import (
	"astrologerService/entity"
	"astrologerService/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Handler struct {
	storage *storage.Storage
	logger  *slog.Logger
}

func NewHandler(s *storage.Storage, l *slog.Logger) *Handler {
	return &Handler{storage: s, logger: l}
}

const (
	api = "https://api.nasa.gov/planetary/apod?api_key=BipbZ0NGtursCd4rCnXGZV8eFbrvqAJdlnJSm9wQ&date="
)

func (h *Handler) Information(ctx context.Context) error {
	t := time.Now()
	date := t.Format("2006-1-2")

	response, err := http.Get(fmt.Sprintf("%s%s", api, date))
	if err != nil {
		return err
	}

	defer response.Body.Close()

	var metadata entity.Metadata

	err = json.NewDecoder(response.Body).Decode(&metadata)
	if err != nil {
		return err
	}

	err = h.storage.SaveData(ctx, metadata)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) Apod(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dataBases, err := h.storage.GetMetaDataBases(ctx)
	if err != nil {
		http.Error(w, "get meta DataBases", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(dataBases)
	if err != nil {
		http.Error(w, "encoding", http.StatusInternalServerError)
		return
	}

}

func (h *Handler) ApodByDate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	date := r.URL.Query().Get("date")

	data, err := h.storage.GetMetaData(ctx, date)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
		}

		http.Error(w, "get metadata", http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "encoding", http.StatusInternalServerError)
		return
	}
}
