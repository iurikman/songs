package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/iurikman/songs/internal/config"
	"github.com/iurikman/songs/internal/rest"
	"github.com/iurikman/songs/internal/service"
	"github.com/iurikman/songs/internal/store"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	defer cancel()

	cfg := config.NewConfig()

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

	log.Info("successful migration")

	svc := service.NewService(db)

	serverConfig := rest.SrvConfig{BindAddr: cfg.BindAddress}

	svr, err := rest.NewServer(serverConfig, svc)
	if err != nil {
		log.Panicf("rest.NewServer(serverConfig, svc) err: %v", err)
	}

	if err := svr.Start(ctx); err != nil {
		log.Panicf("svr.Start() err: %v", err)
	}

	log.Info("server started")
}
