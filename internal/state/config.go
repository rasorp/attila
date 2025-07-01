// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"errors"

	"github.com/rasorp/attila/internal/helper/pointer"
)

type Config struct {
	Memory *MemoryConfig `hcl:"memory,optional"`
}

type MemoryConfig struct {
	Enabled *bool `hcl:"enabled,optional"`
}

func DefaultConfig() *Config {
	return &Config{
		Memory: &MemoryConfig{Enabled: pointer.Of(true)},
	}
}

func (c *Config) Validate() error {

	if c == nil {
		return errors.New("state config block required")
	}
	if c.Memory == nil || !*c.Memory.Enabled {
		return errors.New("memory state must be enabled")
	}
	return nil
}

func (c *Config) Merge(z *Config) *Config {
	if c == nil {
		return z
	}

	result := *c

	if c.Memory != nil {
		if c.Memory.Enabled != nil {
			result.Memory.Enabled = c.Memory.Enabled
		}
	}

	return &result
}
