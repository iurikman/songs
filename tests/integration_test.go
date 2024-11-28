package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/iurikman/songs/internal/config"
	"github.com/iurikman/songs/internal/models"
	"github.com/iurikman/songs/internal/rest"
	"github.com/iurikman/songs/internal/service"
	"github.com/iurikman/songs/internal/songdetails"
	"github.com/iurikman/songs/internal/store"
	_ "github.com/jackc/pgx/v5/stdlib"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/suite"
)

const bindAddress = "http://localhost:8080/api/v1/songs"

type IntegrationTestSuite struct {
	suite.Suite
	cancel     context.CancelFunc
	store      *store.Postgres
	service    *service.Service
	server     *rest.Server
	mockserver *httptest.Server
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	cfg := config.NewConfig()

	db, err := store.New(ctx, store.Config{
		PGUser:     cfg.PostgresUser,
		PGPassword: cfg.PostgresPassword,
		PGHost:     cfg.PostgresHost,
		PGPort:     cfg.PostgresPort,
		PGDatabase: cfg.PostgresDatabase,
	})
	s.Require().NoError(err)

	s.store = db

	err = s.store.Migrate(migrate.Up)
	s.Require().NoError(err)

	err = s.store.Truncate(ctx, "songs")
	s.Require().NoError(err)

	s.mockserver = httptest.NewServer(http.HandlerFunc(handler))

	songDetails := songdetails.NewSongDetails(s.mockserver.URL)

	s.service = service.NewService(db, songDetails)

	s.server, err = rest.NewServer(rest.SrvConfig{BindAddr: os.Getenv("BIND_ADDRESS")}, s.service)
	s.Require().NoError(err)

	go func() {
		err := s.server.Start(ctx)
		s.Require().NoError(err)
	}()
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.cancel()
}

func (s *IntegrationTestSuite) sendRequest(ctx context.Context, method, endpoint string, body interface{}, dest interface{}) *http.Response {
	s.T().Helper()

	reqBody, err := json.Marshal(body)
	s.Require().NoError(err)

	req, err := http.NewRequestWithContext(ctx, method, bindAddress+endpoint, bytes.NewBuffer(reqBody))
	s.Require().NoError(err)

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)

	defer func() {
		err = resp.Body.Close()
		s.Require().NoError(err)
	}()

	if dest != nil {
		err = json.NewDecoder(resp.Body).Decode(&dest)
		s.Require().NoError(err)
	}

	return resp
}

func handler(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")

	if group == "" || song == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Симуляция данных для примера
	detail := models.SongDetails{
		ReleaseDate: "16.07.2006",
		Text: `Ooh baby, don't you know I suffer?
Ooh baby, can you hear me moan?
You caught me under false pretenses
How long before you let me go?

Ooh
You set my soul alight
Ooh
You set my soul alight`,
		Link: "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(detail)
	if err != nil {
		return
	}
}
