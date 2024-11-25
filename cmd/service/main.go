package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/iurikman/songs/internal/rest"
	"github.com/iurikman/songs/internal/service"
	"github.com/iurikman/songs/internal/store"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	defer cancel()

	err := godotenv.Load(".env")
	if err != nil {
		log.Panicf("Error loading .env file")
	}

	storeConfig := store.Config{
		PGUser:     os.Getenv("POSTGRES_USER"),
		PGPassword: os.Getenv("POSTGRES_PASSWORD"),
		PGHost:     os.Getenv("POSTGRES_HOST"),
		PGPort:     os.Getenv("POSTGRES_PORT"),
		PGDatabase: os.Getenv("POSTGRES_DATABASE"),
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

	serverConfig := rest.SrvConfig{BindAddr: os.Getenv("BIND_ADDRESS")}

	svr, err := rest.NewServer(serverConfig, svc)
	if err != nil {
		log.Panicf("rest.NewServer(serverConfig, svc) err: %v", err)
	}

	if err := svr.Start(ctx); err != nil {
		log.Panicf("svr.Start() err: %v", err)
	}
}
