package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"github.com/iurikman/songs/internal/models"
	log "github.com/sirupsen/logrus"
)

const (
	standardPage = 10
)

type HTTPResponse struct {
	Data  any    `json:"data"`
	Error string `json:"error"`
}

type service interface {
	CreateSong(ctx context.Context, song models.Song) (*models.Song, error)
	GetSongs(ctx context.Context, params models.Params) ([]*models.Song, error)
	GetText(ctx context.Context, id uuid.UUID, verse int) (*string, error)
	DeleteSong(ctx context.Context, id uuid.UUID) error
	UpdateSong(ctx context.Context, id uuid.UUID, song models.Song) (*models.Song, error)
}

func (s *Server) createSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())

		return
	}

	createSong, err := s.svc.CreateSong(r.Context(), song)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "internal server error")
		log.Warnf("s.svc.CreateSong(r.Context(), song) err: %v", err)

		return
	}

	writeOKResponse(w, http.StatusCreated, createSong)
}

func (s *Server) getSongs(w http.ResponseWriter, r *http.Request) {
	params, err := parseParams(r.URL.Query())
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid query parameters")

		return
	}

	songs, err := s.svc.GetSongs(r.Context(), *params)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "internal server error")

		return
	}

	writeOKResponse(w, http.StatusOK, songs)
}

func (s *Server) getText(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid id")

		return
	}

	verse, err := strconv.Atoi(r.URL.Query().Get("verse"))
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid verse")

		return
	}

	text, err := s.svc.GetText(r.Context(), id, verse)

	switch {
	case errors.Is(err, models.ErrSongNotFound):
		writeErrorResponse(w, http.StatusNotFound, err.Error())

		return
	case errors.Is(err, models.ErrVerseIsNotValid):
		writeErrorResponse(w, http.StatusBadRequest, err.Error())

		return
	case err != nil:
		writeErrorResponse(w, http.StatusInternalServerError, "internal server error")

		return
	}

	writeOKResponse(w, http.StatusOK, text)
}

func (s *Server) deleteSong(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid id")

		return
	}

	err = s.svc.DeleteSong(r.Context(), id)

	switch {
	case errors.Is(err, models.ErrSongNotFound):
		writeErrorResponse(w, http.StatusNotFound, err.Error())

		return

	case err != nil:
		writeErrorResponse(w, http.StatusInternalServerError, "internal server error")

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) updateSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid id")

		return
	}

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
	}

	updatedSong, err := s.svc.UpdateSong(r.Context(), id, song)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "internal server error")

		return
	}

	writeOKResponse(w, http.StatusOK, updatedSong)
}

func parseParams(value url.Values) (*models.Params, error) {
	decoder := schema.NewDecoder()
	params := &models.Params{}

	err := decoder.Decode(params, value)
	if err != nil {
		return nil, fmt.Errorf("decoder.Decode(params, value): %w", err)
	}

	if params.Limit == 0 {
		params.Limit = standardPage
	}

	return params, nil
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, description string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(HTTPResponse{Error: description}); err != nil {
		log.Warnf("json.NewEncoder(w).Encode(HTTPResponse{Error: description}) err: %v", err)
	}
}

func writeOKResponse(w http.ResponseWriter, statusCode int, respData any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(HTTPResponse{Data: respData}); err != nil {
		log.Warnf("json.NewEncoder(w).Encode(HTTPResponse{Error: respData}) err: %v", err)
	}
}
