package logger

import (
	"github.com/maxwellzp/golang-chat-api/internal/config"
	"go.uber.org/zap"
)

type Logger struct {
	log *zap.SugaredLogger
}

func NewLogger(cfg *config.Config) (*Logger, error) {
	var zapLogger *zap.Logger
	var err error

	if cfg.Application.AppEnv == "prod" {
		zapLogger, err = zap.NewProduction()
	} else {
		zapLogger, err = zap.NewDevelopment()
	}
	if err != nil {
		return nil, err
	}
	return &Logger{log: zapLogger.Sugar()}, nil
}

func (l *Logger) Fatalw(msg string, keysAndValues ...any) {
	l.log.Fatalw(msg, keysAndValues...)
}

func (l *Logger) Infow(msg string, keysAndValues ...any) {
	l.log.Infow(msg, keysAndValues...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...any) {
	l.log.Warnw(msg, keysAndValues...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...any) {
	l.log.Errorw(msg, keysAndValues...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...any) {
	l.log.Debugw(msg, keysAndValues...)
}
func (l *Logger) Panicw(msg string, keysAndValues ...any) {
	l.log.Panicw(msg, keysAndValues...)
}
