package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

type Logger interface {
	Info(ctx context.Context, msg string, keyValues ...any)
	Debug(ctx context.Context, msg string, keyValues ...any)
	Error(ctx context.Context, err error, msg string, keyValues ...any)
	Warn(ctx context.Context, err error, msg string, keyValues ...any)
	Fatal(ctx context.Context, err error, msg string, keyValues ...any)

	WithLogger(ctx context.Context) context.Context
	FromContext(ctx context.Context) *zerolog.Logger
	WithContext(ctx context.Context, keyValues ...any) context.Context
}

type ctxKey string

const (
	loggerKey ctxKey = "logger"
)

type Adapter struct {
	Logger zerolog.Logger
}

var globalLog *Adapter

// Инициализация глобального адаптера
func InitLogger(env, logPath string) *Adapter {
	logger := ConfigureLogger(env, logPath)
	globalLog = &Adapter{
		Logger: logger,
	}
	return globalLog
}

func GetZerologLogger() *zerolog.Logger {
	return &globalLog.Logger
}

// Добавление логгера в контекст
func (a *Adapter) WithLogger(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey, a.Logger)
}

// Получение логгера из контекста
func (a *Adapter) FromContext(ctx context.Context) *zerolog.Logger {
	logger, ok := ctx.Value(loggerKey).(zerolog.Logger)
	if !ok {
		return &a.Logger
	}
	return &logger
}

// Добавление атрибутов в контекст и обновление логгера
func (a *Adapter) WithContext(ctx context.Context, keyValues ...any) context.Context {
	currentLogger := a.FromContext(ctx)

	ev := currentLogger.With()
	for i := 0; i < len(keyValues); i += 2 {
		if i+1 < len(keyValues) {
			key, ok := keyValues[i].(string)
			if !ok {
				continue
			}
			ev = ev.Interface(key, keyValues[i+1])
		}
	}
	newLogger := ev.Logger()
	return context.WithValue(ctx, loggerKey, newLogger)
}

// Логирование уровня INFO
func (a *Adapter) Info(ctx context.Context, msg string, keyValues ...any) {
	currentLogger := a.FromContext(ctx)
	event := currentLogger.Info()
	a.addFields(event, keyValues...).Msg(msg)
}

// Логирование уровня DEBUG
func (a *Adapter) Debug(ctx context.Context, msg string, keyValues ...any) {
	currentLogger := a.FromContext(ctx)
	event := currentLogger.Debug()
	a.addFields(event, keyValues...).Msg(msg)
}

// Логирование уровня ERROR
func (a *Adapter) Error(ctx context.Context, err error, msg string, keyValues ...any) {
	currentLogger := a.FromContext(ctx)
	event := currentLogger.Error().Err(err)
	a.addFields(event, keyValues...).Msg(msg)
}

// Логирование уровня WARN
func (a *Adapter) Warn(ctx context.Context, err error, msg string, keyValues ...any) {
	currentLogger := a.FromContext(ctx)
	event := currentLogger.Warn().Err(err)
	a.addFields(event, keyValues...).Msg(msg)
}

// Логирование уровня FATAL
func (a *Adapter) Fatal(ctx context.Context, err error, msg string, keyValues ...any) {
	currentLogger := a.FromContext(ctx)
	event := currentLogger.Fatal().Err(err)
	a.addFields(event, keyValues...).Msg(msg)
	os.Exit(1)
}

// Вспомогательная функция добавления полей в Event
func (a *Adapter) addFields(event *zerolog.Event, keyValues ...any) *zerolog.Event {
	for i := 0; i < len(keyValues); i += 2 {
		if i+1 < len(keyValues) {

			key, ok := keyValues[i].(string)
			if !ok {
				continue
			}

			event = event.Interface(key, keyValues[i+1])
		}
	}
	return event
}

func WithContext(ctx context.Context, keyValues ...any) context.Context {
	return globalLog.WithContext(ctx, keyValues...)
}

// Info логирует информационное сообщение
func InfoCtx(ctx context.Context, msg string, keyValues ...any) {
	globalLog.Info(ctx, msg, keyValues...)
}

// Debug логирует отладочное сообщение
func DebugCtx(ctx context.Context, msg string, keyValues ...any) {
	globalLog.Debug(ctx, msg, keyValues...)
}

// Error логирует сообщение об ошибке
func ErrorCtx(ctx context.Context, err error, msg string, keyValues ...any) {
	globalLog.Error(ctx, err, msg, keyValues...)
}

// Warn логирует предупреждение
func WarnCtx(ctx context.Context, err error, msg string, keyValues ...any) {
	globalLog.Warn(ctx, err, msg, keyValues...)
}

// Fatal логирует фатальную ошибку
func FatalCtx(ctx context.Context, err error, msg string, keyValues ...any) {
	globalLog.Fatal(ctx, err, msg, keyValues...)
}

func Info(msg string, keyValues ...any) {
	event := globalLog.Logger.Info()
	globalLog.addFields(event, keyValues...).Msg(msg)
}

func Debug(msg string, keyValues ...any) {
	event := globalLog.Logger.Debug()
	globalLog.addFields(event, keyValues...).Msg(msg)
}

func Error(err error, msg string, keyValues ...any) {
	event := globalLog.Logger.Error().Err(err)
	globalLog.addFields(event, keyValues...).Msg(msg)
}

func Warn(err error, msg string, keyValues ...any) {
	event := globalLog.Logger.Warn().Err(err)
	globalLog.addFields(event, keyValues...).Msg(msg)
}

func Fatal(err error, msg string, keyValues ...any) {
	event := globalLog.Logger.Fatal().Err(err)
	globalLog.addFields(event, keyValues...).Msg(msg)
}

func WithLogger(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey, globalLog)
}
