package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/iurikman/songs/internal/config"
	"github.com/iurikman/songs/internal/rest"
	"github.com/iurikman/songs/internal/service"
	"github.com/iurikman/songs/internal/songdetails"
	"github.com/iurikman/songs/internal/store"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
)

// @title Songs API
// @version 1.0
// @description API for managing songs
// @contact.name Iurikman
// @host localhost:8080
// @BasePath /api/v1

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	defer cancel()

	cfg := config.NewConfig()

	log.Debug("configuration initialized")

	storeConfig := store.Config{
		PGUser:     cfg.PostgresUser,
		PGPassword: cfg.PostgresPassword,
		PGHost:     cfg.PostgresHost,
		PGPort:     cfg.PostgresPort,
		PGDatabase: cfg.PostgresDatabase,
	}

	db, err := store.New(ctx, storeConfig)
	if err != nil {
		log.Panicf("store.New(ctx, storeConfig) err: %v", err)
	}

	if err := db.Migrate(migrate.Up); err != nil {
		log.Panicf("db.Migrate(migrate.Up) err: %v", err)
	}

	log.Debug("successful migration")

	songDetails := songdetails.NewSongDetails(cfg.APIUrl + cfg.APIPort)

	svc := service.NewService(db, songDetails)

	log.Debug("service initialized")

	serverConfig := rest.SrvConfig{BindAddr: cfg.BindAddress}

	svr, err := rest.NewServer(serverConfig, svc)
	if err != nil {
		log.Panicf("rest.NewServer(serverConfig, svc) err: %v", err)
	}

	log.Debug("rest server initialized")

	if err := svr.Start(ctx); err != nil {
		log.Panicf("svr.Start() err: %v", err)
	}

	log.Info("server stopped")
}
