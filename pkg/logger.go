package pkg

import (
	"go.uber.org/zap"
)

// Logger on case of changing logger in the future
type Logger interface {
	Info(msg string)
	Fatal(msg string)
	Debug(msg string)
	Warn(msg string)
}

type logger struct {
	l *zap.Logger
}

func New(lg *zap.Logger) Logger {
	return &logger{l: lg}
}

func (l logger) Info(msg string) {
	l.l.Info(msg)
}

func (l logger) Fatal(msg string) {
	l.l.Fatal(msg)
}

func (l logger) Debug(msg string) {
	l.l.Debug(msg)
}

func (l logger) Warn(msg string) {
	l.l.Warn(msg)
}
