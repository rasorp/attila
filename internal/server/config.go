// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-set/v3"
	"github.com/rs/zerolog"

	"github.com/rasorp/attila/internal/logger"
	storebackend "github.com/rasorp/attila/internal/store/backend"
)

type Config struct {
	Log   *logger.Config       `hcl:"log,optional"`
	State *storebackend.Config `hcl:"state,optional"`
	HTTP  *HTTPConfig          `hcl:"http,optional"`
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

	var errs []error

	if err := c.Log.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := c.State.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := c.HTTP.Validate(); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
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

	var errs []error

	if len(h.Binds) < 1 {
		errs = append(errs, errors.New("http bind address required"))
	}
	if _, err := zerolog.ParseLevel(strings.ToLower(h.AccessLogLevel)); err != nil {
		errs = append(errs, fmt.Errorf("failed to parse access log level: %w", err))
	}

	for _, bind := range h.Binds {
		parsedURL, err := url.Parse(bind.Addr)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to parse bind address: %w", err))
		} else {
			switch parsedURL.Scheme {
			case "unix", "http", "https":
			default:
				errs = append(errs, fmt.Errorf("unsupported bind protocol %q", parsedURL.Scheme))
			}
		}
	}

	return errors.Join(errs...)
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
		State: storebackend.DefaultConfig(),
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
