package crontab

import "github.com/rs/zerolog"

type Logger struct {
	logger  zerolog.Logger
	logInfo bool
}

func (l *Logger) Info(msg string, keysAndValues ...any) {
	if l.logInfo {
		l.logger.Info().Fields(keysAndValues).Msg(msg)
	}
}

func (l *Logger) Error(err error, msg string, keysAndValues ...any) {
	l.logger.Err(err).Fields(keysAndValues).Msg(msg)
}
