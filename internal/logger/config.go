// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog"

	"github.com/rasorp/attila/internal/helper/pointer"
)

type Config struct {
	Level       string `hcl:"level,optional"`
	Format      string `hcl:"format,optional"`
	Colour      *bool  `hcl:"colour,optional"`
	IncludeLine *bool  `hcl:"include_line,optional"`
}

func DefaultConfig() *Config {
	return &Config{
		Level:       zerolog.LevelInfoValue,
		Format:      "json",
		Colour:      pointer.Of(false),
		IncludeLine: pointer.Of(false),
	}
}

func (c *Config) Validate() error {

	if c == nil {
		return errors.New("log config block required")
	}

	var mErr *multierror.Error

	if _, err := zerolog.ParseLevel(strings.ToLower(c.Level)); err != nil {
		mErr = multierror.Append(mErr, fmt.Errorf("failed to parse level: %w", err))
	}

	switch c.Format {
	case "human", "json":
	default:
		mErr = multierror.Append(mErr, fmt.Errorf("unsupported format: %q", c.Format))
	}

	return mErr.ErrorOrNil()
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
