// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"errors"

	"github.com/rasorp/attila/internal/server/state"
	"github.com/rasorp/attila/internal/state/mem"
)

func NewBackend(cfg *Config) (state.State, error) {

	if cfg.Memory != nil && *cfg.Memory.Enabled {
		backend, err := mem.New()
		if err != nil {
			return nil, err
		}
		return backend, nil
	}

	return nil, errors.New("no state backend configured")
}
