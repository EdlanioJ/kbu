package log

import (
	"context"
	stdLog "log"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Level zerolog.Level

func SetLevel(levelString string) {
	level, _ := strconv.Atoi(levelString)
	Level = zerolog.Level(level)
}

func New(ctx context.Context, id string) context.Context {
	zerolog.SetGlobalLevel(Level)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano})
	logger := log.With().Str("request_id", id).Logger()
	ctx = logger.WithContext(ctx)
	return ctx
}

func Insert(ctx context.Context, key, value string) context.Context {
	logger := log.Ctx(ctx).With().Str(key, value).Logger()
	ctx = logger.WithContext(ctx)
	return ctx
}

func Debug(ctx context.Context, message string) {
	if ctx == nil {
		stdLog.Println(message)
		return
	}
	if e := log.Debug(); !e.Enabled() {
		return
	}
	log.Ctx(ctx).Debug().Msg(message)
}

func Debugf(ctx context.Context, message string, v ...interface{}) {
	if ctx == nil {
		stdLog.Println(message)
		return
	}
	if e := log.Debug(); !e.Enabled() {
		return
	}
	log.Ctx(ctx).Debug().Msgf(message, v...)
}

func Info(ctx context.Context, message string) {
	if ctx == nil {
		stdLog.Println(message)
		return
	}
	if e := log.Info(); !e.Enabled() {
		return
	}
	log.Ctx(ctx).Info().Msg(message)
}

func Infof(ctx context.Context, message string, v ...interface{}) {
	if ctx == nil {
		stdLog.Println(message)
		return
	}
	if e := log.Info(); !e.Enabled() {
		return
	}
	log.Ctx(ctx).Info().Msgf(message, v...)
}

func Warning(ctx context.Context, message string) {
	if ctx == nil {
		stdLog.Println(message)
		return
	}
	if e := log.Warn(); !e.Enabled() {
		return
	}
	log.Ctx(ctx).Warn().Msg(message)
}

func Warningf(ctx context.Context, message string, v ...interface{}) {
	if ctx == nil {
		stdLog.Println(message)
		return
	}
	if e := log.Warn(); !e.Enabled() {
		return
	}
	log.Ctx(ctx).Warn().Msgf(message, v...)
}

func Error(ctx context.Context, err error, message string) {
	if ctx == nil {
		stdLog.Println(message)
		return
	}
	if e := log.Err(nil); !e.Enabled() {
		return
	}
	log.Ctx(ctx).Err(err).Msg(message)
}

func Errorf(ctx context.Context, err error, message string, v ...interface{}) {
	if ctx == nil {
		stdLog.Println(message)
		return
	}
	if e := log.Err(nil); !e.Enabled() {
		return
	}
	log.Ctx(ctx).Err(err).Msgf(message, v...)
}

func Fatal(ctx context.Context, message string) {
	if ctx == nil {
		stdLog.Println(message)
		return
	}
	if e := log.Fatal(); !e.Enabled() {
		return
	}
	log.Ctx(ctx).Fatal().Msg(message)
}
