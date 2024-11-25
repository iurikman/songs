package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

const (
	readHeaderTimeout = 5 * time.Second
	maxHeaderBytes    = 1 << 20
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Panicf("Error loading .env file")
	}

	router := chi.NewRouter()
	srv := &http.Server{
		Addr:              os.Getenv("API_PORT"),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	router.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		respData := SongDetail{
			ReleaseDate: "16.07.2006",
			Text: "Ooh baby, don't you know I suffer?\n" +
				"Ooh baby, can you hear me moan?\n" +
				"You caught me under false pretenses\n" +
				"How long before you let me go?" +
				"\n\n" +
				"Ooh\n" +
				"You set my soul alight\n" +
				"Ooh\\" +
				"nYou set my soul alight",
			Link: "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
		}

		if err := json.NewEncoder(w).Encode(respData); err != nil {
			log.Warnf("json.NewEncoder(w).Encode(HTTPResponse{Error: respData}) err: %v", err)
		}
	})

	err := srv.ListenAndServe()
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
