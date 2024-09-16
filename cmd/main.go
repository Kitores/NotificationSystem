package main

import (
	"NotificationSystem/internal/config"
	"NotificationSystem/internal/setupLogger"
	"log/slog"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger.SetupLogger(cfg.Env)

	log.Info("starting application", slog.String("env", cfg.Env))

	//storage, err :=

	//TODO: config
	//TODO: logging
	//TODO: database
	//TODO: server

}
