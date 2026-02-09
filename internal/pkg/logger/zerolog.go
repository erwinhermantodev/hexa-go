package logger

import (
	"os"

	"github.com/erwinhermantodev/hexa-go/internal/pkg/interfaces"
	"github.com/rs/zerolog"
)

type zerologWrapper struct {
	logger zerolog.Logger
}

func NewZeroLog() interfaces.Logger {
	l := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return &zerologWrapper{logger: l}
}

func (l *zerologWrapper) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Debug().Fields(keysAndValues).Msg(msg)
}

func (l *zerologWrapper) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Info().Fields(keysAndValues).Msg(msg)
}

func (l *zerologWrapper) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Warn().Fields(keysAndValues).Msg(msg)
}

func (l *zerologWrapper) Error(err error, msg string, keysAndValues ...interface{}) {
	l.logger.Error().Err(err).Fields(keysAndValues).Msg(msg)
}

func (l *zerologWrapper) Fatal(err error, msg string, keysAndValues ...interface{}) {
	l.logger.Fatal().Err(err).Fields(keysAndValues).Msg(msg)
}

func (l *zerologWrapper) With(keysAndValues ...interface{}) interfaces.Logger {
	return &zerologWrapper{logger: l.logger.With().Fields(keysAndValues).Logger()}
}
