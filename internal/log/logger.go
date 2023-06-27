package log

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	log.Logger = log.Logger.With().Caller().Logger()
	zerolog.DefaultContextLogger = &log.Logger
}

func Default() *zerolog.Logger {
	return &log.Logger
}

func WithContext(ctx context.Context) *zerolog.Logger {
	return log.Ctx(ctx)
}

func CreateContext(ctx context.Context, key string, value interface{}) context.Context {
	return log.Logger.With().
		Interface(key, value).
		Logger().
		WithContext(ctx)
}
