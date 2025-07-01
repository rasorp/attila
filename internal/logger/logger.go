// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func New(cfg *Config) (*zerolog.Logger, error) {

	// Parsing the log level is also performed when validating the config
	// options, but it's easy enough to check. We don't need to add anything to
	// the error, as the caller will handle wrapping it and it contains all the
	// information needed already.
	zerologLevel, err := zerolog.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		return nil, err
	}

	// Generate the base logger which will include the timestamp and send all
	// logs to stderr.
	zerologLogger := zerolog.New(os.Stderr).
		Level(zerologLevel).
		With().Timestamp().
		Logger()

	// If the user specified human logging format, set the output accordingly.
	// The default formatting of the timestamp and level is "8:46AM INF" which
	// is not the best, so these are changed.
	//
	// The level format is chosen to given a consistent indentation to the logs
	// which makes it easier to grok.
	if cfg.Format == "human" {
		zerologLogger = zerologLogger.Output(zerolog.ConsoleWriter{
			FormatLevel: func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("%-5s", i))
			},
			TimeFormat: time.RFC3339,
			Out:        os.Stderr,
			NoColor:    !*cfg.Colour,
		})
	}

	if *cfg.IncludeLine {
		zerologLogger = zerologLogger.With().Caller().Logger()
	}

	return &zerologLogger, nil
}
