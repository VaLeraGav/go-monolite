package main

import (
	"flag"
	"go-monolite/internal/app"
	"go-monolite/internal/config"
	"go-monolite/pkg/helper"
)

func main() {
	currentDir := helper.GetProjectPath()

	envFile := flag.String("env", "", "path to config file")
	flag.Parse()

	configPath := config.Path(currentDir, *envFile)

	config := config.MustInit(configPath)

	app.InitApp(config, currentDir)
}
