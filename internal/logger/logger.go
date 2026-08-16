// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	timeKey = "timestamp"
	nameKey = "component"
)

func New(cfg *Config) (*zap.Logger, error) {

	// Zap does not sanitize the input string, so do that here to avoid
	// case-sensetive config parameters.
	logLevel, err := zap.ParseAtomicLevel(strings.ToLower(cfg.Level))
	if err != nil {
		return nil, err
	}

	var encoder zapcore.Encoder

	switch cfg.Format {
	case formatHuman:
		encoder = newHumanEncoder(*cfg.Colour)
	default:
		encoder = newJSONEncoder()
	}

	// Accumulate our options; currently we only support adding the log line
	// detail.
	var opts []zap.Option

	if *cfg.IncludeLine {
		opts = append(opts, zap.AddCaller())
	}

	return zap.New(zapcore.NewCore(encoder, os.Stderr, logLevel), opts...), nil
}

func newHumanEncoder(colour bool) zapcore.Encoder {

	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = timeKey
	cfg.NameKey = nameKey
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder

	if colour {
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	}
	return zapcore.NewConsoleEncoder(cfg)
}

func newJSONEncoder() zapcore.Encoder {

	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = timeKey
	cfg.NameKey = nameKey
	cfg.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	return zapcore.NewJSONEncoder(cfg)
}
