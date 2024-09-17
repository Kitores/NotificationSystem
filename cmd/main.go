package main

import (
	"NotificationSystem/internal/config"
	"NotificationSystem/internal/lib/logger/sl"
	"NotificationSystem/internal/setupLogger"
	"NotificationSystem/internal/storage/postgres"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger.SetupLogger(cfg.Env)

	log.Info("starting application", slog.String("env", cfg.Env))

	conn := &config.Storage{
		Host:     cfg.Storage.Host,
		Port:     cfg.Storage.Port,
		UserDb:   cfg.Storage.UserDb,
		Password: cfg.Storage.Password,
		Dbname:   cfg.Storage.Dbname,
		SSLmode:  cfg.Storage.SSLmode,
	}
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", conn.Host, conn.Port, conn.UserDb, conn.Password, conn.Dbname, conn.SSLmode)
	storage, err := postgres.NewPg(connStr)
	if err != nil {
		log.Error("failed to initialize storage: %v", sl.Err(err))
		os.Exit(1)
	}
	fmt.Println(storage)

	//TODO: config
	//TODO: logging
	//TODO: database
	//TODO: server

}
