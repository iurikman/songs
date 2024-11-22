package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/iurikman/songs/internal/models"
	log "github.com/sirupsen/logrus"
)

type HTTPResponse struct {
	Data  any    `json:"data"`
	Error string `json:"error"`
}

type service interface {
	CreateSong(ctx context.Context, song models.Song) (*models.Song, error)
}

func (s *Server) createSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var song models.Song

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
	}

	createSong, err := s.svc.CreateSong(r.Context(), song)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "internal server error")
		log.Warnf("s.svc.CreateSong(r.Context(), song) err: %v", err)

		return
	}

	writeOKResponse(w, http.StatusCreated, createSong)
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
