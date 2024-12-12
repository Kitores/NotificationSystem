package main

import (
	"NotificationSystem/internal/config"
	grpc_server "NotificationSystem/internal/grpc/grpc-server"
	"NotificationSystem/internal/setupLogger"
	//"NotificationSystem/lib/logger/sl"
	//"fmt"

	"log/slog"
	//"os"
)

const (
	port = ":50051"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger.SetupLogger(cfg.Env)

	log.Info("starting application", slog.String("env", cfg.Env))

	//}
	//connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", conn.Host, conn.Port, conn.UserDb, conn.Password, conn.Dbname, conn.SSLmode)
	//storage, err := postgres.NewPg(connStr)
	//if err != nil {
	//	log.Error("failed to initialize storage: %v", sl.Err(err))
	//	os.Exit(1)
	//}
	//fmt.Println(storage)

	//in case of a remote database
	grpc_server.RunGRPCServe(log)

	//TODO: config
	//TODO: logging
	//TODO: database
	//TODO: server

}
