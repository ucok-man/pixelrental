package logging

import (
	"bytes"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type Meta map[string]any

type Logger struct {
	log *zerolog.Logger
}

func New() *Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	consoleLog := zerolog.ConsoleWriter{
		Out:           os.Stdout,
		TimeFormat:    time.Kitchen,
		FieldsExclude: []string{zerolog.ErrorStackFieldName},
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.MessageFieldName,
		},
		FormatExtra: func(m map[string]interface{}, b *bytes.Buffer) error {
			obj, ok := m[zerolog.ErrorStackFieldName]
			if !ok {
				return nil
			}
			if stack, ok := obj.(string); ok {
				b.WriteString("\nstack:\n")
				b.WriteString(stack)
				return nil
			}
			return nil
		},
	}

	log := zerolog.New(consoleLog).With().Timestamp().Logger()
	return &Logger{
		log: &log,
	}
}

func (l Logger) Info(msg string, properties Meta) {
	cEvent := l.log.Info()
	for key, val := range properties {
		cEvent.Any(key, val)
	}
	cEvent.Msg(msg)
}

func (l Logger) Error(err error, msg string, properties Meta) {
	cEvent := l.log.Error().Err(errors.Wrap(err, ""))
	for key, val := range properties {
		cEvent.Any(key, val)
	}
	cEvent.Msg(msg)
}

func (l Logger) Fatal(err error, msg string, properties Meta) {
	cEvent := l.log.Fatal().Err(errors.Wrap(err, ""))
	for key, val := range properties {
		cEvent.Any(key, val)
	}
	cEvent.Msg(msg)
}
