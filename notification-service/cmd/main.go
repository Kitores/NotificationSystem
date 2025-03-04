package main

import (
	"github.com/Kitores/NotificationSystem/notification-service/internal/config"
	grpc_server "github.com/Kitores/NotificationSystem/notification-service/internal/grpc/grpc-server"
	"github.com/Kitores/NotificationSystem/notification-service/internal/setupLogger"
	"log/slog"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger.SetupLogger(cfg.Env)

	log.Info("starting application", slog.String("env", cfg.Env))

	grpc_server.RunGRPCServe(log, cfg)

	//TODO: config
	//TODO: logging
	//TODO: database
	//TODO: server
}
