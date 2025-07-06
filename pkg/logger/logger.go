package logger

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
	EnvTest  = "test"
)

func ConfigureLogger(env string, logPath string) zerolog.Logger {
	var logger zerolog.Logger
	timeFormat := "2006-01-02 15:04:05"

	switch env {
	case EnvTest:
		zerolog.SetGlobalLevel(zerolog.Disabled)
		logger = zerolog.New(zerolog.Nop())
		return logger

	case EnvLocal:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		logFile := openLogFile(logPath)

		consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: timeFormat}
		fileWriter := zerolog.ConsoleWriter{Out: logFile, TimeFormat: timeFormat, NoColor: true}

		multiWriter := zerolog.MultiLevelWriter(fileWriter, consoleWriter)
		logger = zerolog.New(multiWriter)

	case EnvDev:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: timeFormat})

	case EnvProd:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: timeFormat, NoColor: true})

	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: timeFormat})
		logger.Warn().Msg("Неизвестная среда, используем уровень по умолчанию: Info")
	}

	logger = logger.With().Timestamp().Logger()

	return logger
}

func openLogFile(logPath string) *os.File {
	dir := filepath.Dir(logPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal().Err(err).Msg("Could not create log directory")
	}

	logFile, err := os.OpenFile(
		logPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not open log file")
	}

	return logFile
}

func LogPath(currentDir, logPath string) string {
	if logPath == "" {
		return filepath.Join(currentDir, "logs", "local.log")
	}

	return filepath.Join(currentDir, logPath)
}
