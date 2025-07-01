// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-set/v3"
	"github.com/rs/zerolog"

	"github.com/rasorp/attila/internal/logger"
	"github.com/rasorp/attila/internal/state"
)

type Config struct {
	Log   *logger.Config `hcl:"log,optional"`
	State *state.Config  `hcl:"state,optional"`
	HTTP  *HTTPConfig    `hcl:"http,optional"`
}

func (c *Config) Merge(z *Config) *Config {

	if c == nil {
		return z
	}

	result := *c
	result.Log = c.Log.Merge(z.Log)
	result.State = c.State.Merge(z.State)
	result.HTTP = c.HTTP.Merge(z.HTTP)

	return &result
}

func (c *Config) Validate() error {

	var mErr *multierror.Error

	if err := c.Log.Validate(); err != nil {
		mErr = multierror.Append(mErr, err)
	}

	if err := c.State.Validate(); err != nil {
		mErr = multierror.Append(mErr, err)
	}

	if err := c.HTTP.Validate(); err != nil {
		mErr = multierror.Append(mErr, err)
	}

	return mErr.ErrorOrNil()
}

type HTTPConfig struct {
	Binds          []*BindConfig `hcl:"bind,optional"`
	AccessLogLevel string        `hcl:"access_log_level,optional"`
}

type BindConfig struct {
	Addr string `hcl:"addr,optional"`
}

func (h *HTTPConfig) Validate() error {

	if h == nil {
		return errors.New("http config block required")
	}

	var mErr *multierror.Error

	if len(h.Binds) < 1 {
		mErr = multierror.Append(mErr, errors.New("http bind address required"))
	}
	if _, err := zerolog.ParseLevel(strings.ToLower(h.AccessLogLevel)); err != nil {
		mErr = multierror.Append(mErr, fmt.Errorf("failed to parse access log level: %w", err))
	}

	for _, bind := range h.Binds {
		parsedURL, err := url.Parse(bind.Addr)
		if err != nil {
			mErr = multierror.Append(mErr, fmt.Errorf("failed to parse bind address: %w", err))
		}

		switch parsedURL.Scheme {
		case "unix", "http", "https":
		default:
			mErr = multierror.Append(mErr, fmt.Errorf("unsupported bind protocol %q", parsedURL.Scheme))
		}
	}

	return nil
}

func (h *HTTPConfig) Merge(z *HTTPConfig) *HTTPConfig {

	if h == nil {
		return z
	}

	result := *h

	if z.AccessLogLevel != "" {
		result.AccessLogLevel = z.AccessLogLevel
	}

	// Use a set to deduplicate the bind addresses, so it does the heavy
	// lifitng and ensures accuracy.
	bindSet := set.New[*BindConfig](0)
	bindSet.InsertSlice(result.Binds)
	bindSet.InsertSlice(h.Binds)
	result.Binds = bindSet.Slice()

	return &result
}

// DefaultConfig returns a fully populated server config which is perfectly
// suitable for being used without modification.
func DefaultConfig() *Config {
	return &Config{
		Log:   logger.DefaultConfig(),
		State: state.DefaultConfig(),
		HTTP: &HTTPConfig{
			AccessLogLevel: zerolog.LevelInfoValue,
			Binds: []*BindConfig{
				{
					Addr: "http://127.0.0.1:8080",
				},
			},
		},
	}
}
