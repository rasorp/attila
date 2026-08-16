// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/rasorp/attila/internal/helper/pointer"
)

const (
	formatJSON  = "json"
	formatHuman = "human"
)

type Config struct {
	Level       string `hcl:"level,optional"`
	Format      string `hcl:"format,optional"`
	Colour      *bool  `hcl:"colour,optional"`
	IncludeLine *bool  `hcl:"include_line,optional"`
}

func DefaultConfig() *Config {
	return &Config{
		Level:       zap.InfoLevel.String(),
		Format:      formatJSON,
		Colour:      pointer.Of(false),
		IncludeLine: pointer.Of(false),
	}
}

func (c *Config) Validate() error {

	if c == nil {
		return errors.New("log config block required")
	}

	var errs []error

	if _, err := zap.ParseAtomicLevel(strings.ToLower(c.Level)); err != nil {
		errs = append(errs, fmt.Errorf("failed to parse level: %w", err))
	}

	switch c.Format {
	case formatHuman, formatJSON:
	default:
		errs = append(errs, fmt.Errorf("unsupported format: %q", c.Format))
	}

	return errors.Join(errs...)
}

func (c *Config) Merge(z *Config) *Config {

	if c == nil {
		return z
	}

	result := *c

	if z.Level != "" {
		result.Level = z.Level
	}
	if z.Format != "" {
		result.Format = z.Format
	}
	if z.Colour != nil {
		result.Colour = z.Colour
	}
	if z.IncludeLine != nil {
		result.IncludeLine = z.IncludeLine
	}

	return &result
}
