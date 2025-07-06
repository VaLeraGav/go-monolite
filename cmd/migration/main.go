package main

import (
	"flag"
	"go-monolite/internal/config"
	"go-monolite/internal/migrator"
	"go-monolite/internal/store"
	"go-monolite/pkg/helper"
	"go-monolite/pkg/logger"
	"os"
)

func main() {
	currentDir := helper.GetProjectPath()

	action := flag.String("action", "up", "action to perform: up or down")
	envFile := flag.String("env", "", "path to config file")
	flag.Parse()

	configPath := config.Path(currentDir, *envFile)

	config := config.MustInit(configPath)

	logPath := logger.LogPath(currentDir, "")

	// для отображении сообщений в миграции
	logger.InitLogger(logger.EnvDev, logPath)

	postgresStore, err := store.NewPostgres(config)
	if err != nil {
		logger.Fatal(err, "connectDb error")
	}

	migrator := migrator.MustGetNewMigrator(config.Db.NameDb, postgresStore.Db.DB)
	resp, err := migrator.Run(*action)
	if err != nil {
		logger.Fatal(err, "error start migrator")
		os.Exit(1)
	}
	logger.Info(resp)
}
