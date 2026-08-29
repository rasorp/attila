// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"errors"
	"fmt"
)

type Config struct {
	Provider string `hcl:"provider,block"`
}

const (
	ProviderMemory = "mem"
)

// DefaultConfig returns the default configuration for the Attila storage
// backend. This enables the in-memory backend by default as it's currently the
// only backend.
func DefaultConfig() *Config {
	return &Config{
		Provider: ProviderMemory,
	}
}

// Validate performs validation on the config object and all nested
// configuration blocks. The function can be called safely without checking if
// the object is nil. The returned error could wrap multiple errors and should
// indicate a terminal error in the process which intends to use the config
// object.
func (c *Config) Validate() error {
	if c == nil {
		return errors.New("state config block required")
	}

	var errs []error

	switch c.Provider {
	case ProviderMemory:
	default:
		errs = append(errs, fmt.Errorf("unsupported state provider %q", c.Provider))
	}

	return errors.Join(errs...)
}

func (c *Config) Merge(z *Config) *Config {
	if c == nil {
		return z
	}
	if z == nil {
		return c
	}

	result := *c

	if z.Provider != "" {
		result.Provider = z.Provider
	}

	return &result
}
