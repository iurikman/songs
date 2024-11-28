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

// createSong godoc
// @Summary Create a new song
// @Description Create a new song with the provided details
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Song Data"
// @Success 201 {object} models.Song
// @Failure 400 {object} HTTPResponse
// @Failure 409 {object} HTTPResponse
// @Failure 500 {object} HTTPResponse
// @Router /songs [post].
func (s *Server) createSong(w http.ResponseWriter, r *http.Request) {
	log.Debug("createSong: handler invoked")

	var song models.Song

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())

		return
	}

	log.Debugf("Attempting to create song: %+v", song)

	createSong, err := s.svc.CreateSong(r.Context(), song)

	switch {
	case errors.Is(err, models.ErrDuplicateSong):
		writeErrorResponse(w, http.StatusConflict, err.Error())

		return
	case err != nil:
		writeErrorResponse(w, http.StatusInternalServerError, "internal server error")

		return
	}

	writeOKResponse(w, http.StatusCreated, createSong)
}

// getSongs godoc
// @Summary Get list of songs
// @Description Retrieve a list of songs based on filter and sorting parameters
// @Tags songs
// @Produce json
// @Param filter query string false "Filter by song name"
// @Param sorting query string false "Sort by field (e.g., name)"
// @Param descending query bool false "Sort in descending order"
// @Param offset query int false "Offset for pagination"
// @Param limit query int false "Limit number of songs"
// @Success 200 {array} models.Song
// @Failure 400 {object} HTTPResponse
// @Failure 500 {object} HTTPResponse
// @Router /songs [get].
func (s *Server) getSongs(w http.ResponseWriter, r *http.Request) {
	log.Debug("getSongs: handler invoked")

	params, err := parseParams(r.URL.Query())
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid query parameters")

		return
	}

	log.Debugf("Fetching songs with params: %+v", params)

	songs, err := s.svc.GetSongs(r.Context(), *params)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "internal server error")

		return
	}

	writeOKResponse(w, http.StatusOK, songs)
}

// getText godoc
// @Summary Get song text
// @Description Retrieve the text of a song by ID and verse
// @Tags songs
// @Produce json
// @Param id path string true "Song ID"
// @Param offset query int true "Verse offset"
// @Success 200 {string} string "Song text"
// @Failure 400 {object} HTTPResponse
// @Failure 404 {object} HTTPResponse
// @Failure 500 {object} HTTPResponse
// @Router /songs/{id} [get].
func (s *Server) getText(w http.ResponseWriter, r *http.Request) {
	log.Debug("getText: handler invoked")

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid id")

		return
	}

	verse, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid verse")

		return
	}

	log.Debugf("Retrieving text for song ID: %s, verse offset: %d", id, verse)
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

// deleteSong godoc
// @Summary Delete a song
// @Description Mark a song as deleted by its ID
// @Tags songs
// @Param id path string true "Song ID"
// @Success 204
// @Failure 400 {object} HTTPResponse
// @Failure 404 {object} HTTPResponse
// @Failure 500 {object} HTTPResponse
// @Router /songs/{id} [delete].
func (s *Server) deleteSong(w http.ResponseWriter, r *http.Request) {
	log.Debug("deleteSong: handler invoked")

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid id")

		return
	}

	log.Debugf("Attempting to delete song with ID: %s", id)
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

// updateSong godoc
// @Summary Update a song
// @Description Update details of a song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path string true "Song ID"
// @Param song body models.Song true "Song Data"
// @Success 200 {object} models.Song
// @Failure 400 {object} HTTPResponse
// @Failure 500 {object} HTTPResponse
// @Router /songs/{id} [put].
func (s *Server) updateSong(w http.ResponseWriter, r *http.Request) {
	log.Debug("updateSong: handler invoked")

	var song models.Song

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid id")

		return
	}

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
	}

	log.Debugf("Attempting to update song with ID: %s to: %+v", id, song)

	updatedSong, err := s.svc.UpdateSong(r.Context(), id, song)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "internal server error")

		return
	}

	writeOKResponse(w, http.StatusOK, updatedSong)
}

func parseParams(values url.Values) (*models.Params, error) {
	decoder := schema.NewDecoder()
	params := &models.Params{}

	log.Debug("Parsing query parameters")

	err := decoder.Decode(params, values)
	if err != nil {
		return nil, fmt.Errorf("decoder.Decode(params, values): %w", err)
	}

	if params.Limit == 0 {
		params.Limit = standardPage
	}

	log.Infof("Parsed parameters: %+v", params)

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
		log.Warnf("json.NewEncoder(w).Encode(HTTPResponse{Data: respData}) err: %v", err)
	}
}
