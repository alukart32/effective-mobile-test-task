// Package zerologx provides a custom zerolog.
package zerologx

import (
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var (
	once sync.Once
	log  zerolog.Logger
)

// Get returns zerolog.Logger.
func Get() *zerolog.Logger {
	once.Do(func() {
		logLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
		if err != nil {
			logLevel = int(zerolog.DebugLevel) // default to DEBUG
		}

		zerolog.LevelFieldName = "lvl"
		zerolog.MessageFieldName = "msg"
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

		zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
			return strings.ToUpper(l.String())
		}

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		log = zerolog.New(output).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Logger()
	})

	return &log
}

// Disabled returns a new disabled Zerolog.
func Disabled() *zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	l := zerolog.New(os.Stdout)
	return &l
}
