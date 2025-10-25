package logger

import (
	"os"
	"strings"

	"github.com/audryus/steganocc/config"
	"github.com/rs/zerolog"
)

type iLog interface {
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
}

type Log struct {
	logger zerolog.Logger
}

var _ iLog = (*Log)(nil)

func New(cfg config.Config) *Log {
	l := zerolog.InfoLevel
	zerolog.SetGlobalLevel(l)
	skipFrameCount := 3

	var level zerolog.Level

	switch strings.ToLower(cfg.Logger.Level) {
	case "error":
		level = zerolog.ErrorLevel
	case "warn":
		level = zerolog.WarnLevel
	case "debug":
		level = zerolog.DebugLevel
	default:
		level = zerolog.InfoLevel
	}

	return &Log{
		logger: zerolog.New(os.Stdout).
			With().
			Timestamp().
			CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
			Logger().Level(level),
	}
}

func (l *Log) Engine() *zerolog.Logger {
	return &l.logger
}

func (l *Log) Info(message string, args ...interface{}) {
	l.log(zerolog.InfoLevel, message, args...)
}

func (l *Log) Debug(message string, args ...interface{}) {
	l.log(zerolog.DebugLevel, message, args...)
}

func (l *Log) Warn(message string, args ...interface{}) {
	l.log(zerolog.WarnLevel, message, args...)
}

func (l *Log) Error(message string, args ...interface{}) {
	l.log(zerolog.ErrorLevel, message, args...)
}

func (l *Log) log(level zerolog.Level, message string, args ...interface{}) {
	var event *zerolog.Event

	switch level {
	case zerolog.ErrorLevel:
		event = l.logger.Error()
	case zerolog.DebugLevel:
		event = l.logger.Debug()
	case zerolog.WarnLevel:
		event = l.logger.Warn()
	default:
		event = l.logger.Info()
	}

	if len(args) == 0 {
		event.Msg(message)
	} else {
		event.Msgf(message, args...)
	}
}
