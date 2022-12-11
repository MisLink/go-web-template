package logger

import (
	"io"
	"os"
	"time"

	stdlog "log"

	"github.com/google/wire"
	"github.com/knadh/koanf"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type Options struct {
	Level string
	Json  bool
}

func NewOptions(k *koanf.Koanf) (*Options, error) {
	o := new(Options)
	if err := k.Unmarshal("logger", o); err != nil {
		return nil, err
	}
	return o, nil
}

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano
}

func New(opt *Options) zerolog.Logger {
	level, err := zerolog.ParseLevel(opt.Level)
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(level)
	out := os.Stderr
	var writer io.Writer = os.Stderr
	if !opt.Json {
		writer = zerolog.NewConsoleWriter(
			func(w *zerolog.ConsoleWriter) { w.Out = out },
			func(w *zerolog.ConsoleWriter) {
				if isatty.IsTerminal(out.Fd()) {
					w.NoColor = false
				} else {
					w.NoColor = true
				}
			},
		)
	}
	logger := zerolog.New(writer).With().Caller().Timestamp().Logger()

	log.Logger = logger
	stdlog.SetFlags(0)
	stdlog.SetOutput(logger)
	return logger
}

var ProviderSet = wire.NewSet(New, NewOptions)
