package app

import (
	"go-monolite/internal/config"
	"go-monolite/internal/server"
	"go-monolite/internal/store"
	"go-monolite/pkg/logger"
	"time"
)

func InitApp(config *config.Config, currentDir string) {
	logPath := logger.LogPath(currentDir, "")
	logger.InitLogger(config.Env, logPath)

	postgresStore, err := store.NewPostgres(config)
	if err != nil {
		logger.Fatal(err, "connectDb error")
	}

	s := server.NewServer(config, postgresStore)
	httpServer := s.StartServer()

	GracefulShutdown(
		10*time.Second,
		httpServer,
		postgresStore,
	)
}
