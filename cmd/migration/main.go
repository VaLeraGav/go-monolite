package main

import (
	"flag"
	"go-monolite/internal/config"
	"go-monolite/internal/store"
	"go-monolite/pkg/helper"
	"go-monolite/pkg/logger"
	"go-monolite/pkg/migrator"
)

func main() {
	currentDir := helper.GetProjectPath()

	action := flag.String("action", "up", "action to perform: up, down, create")
	modulePath := flag.String("modulePath", "", "module path for migration (e.g. some_dir)")
	name := flag.String("name", "", "name for migration (e.g. add_users_table)")

	flag.Parse()

	configPath := config.Path(currentDir, "")
	config := config.MustInit(configPath)

	logPath := logger.LogPath(currentDir, "")
	logger.InitLogger(logger.EnvLocal, logPath)

	postgresStore, err := store.NewPostgres(config)
	if err != nil {
		logger.Fatal(err, "connectDb error")
	}

	migrationsDir := "module"
	migrator, err := migrator.New(migrationsDir, config.Db.NameDb, postgresStore.Db.DB)
	if err != nil {
		logger.Fatal(err, "error get migrator")
	}

	switch *action {
	case "create":
		if *modulePath == "" || *name == "" {
			logger.Fatal(nil, "flags -modulePath and -name are required for creating migration")
		}
		if err := migrator.CreateMigration(*modulePath, *name); err != nil {
			logger.Fatal(err, "failed to create migration")
		}
		logger.Info("Migration files created successfully")

	case "up", "down":
		resp, err := migrator.Run(*action)
		if err != nil {
			logger.Fatal(err, "error running migrator")
		}
		logger.Info(resp)

	case "version":
		err := migrator.Version()
		if err != nil {
			logger.Fatal(err, "error running version")
		}

	default:
		logger.Fatal(nil, "unsupported action, use up, down or create")
	}
}
