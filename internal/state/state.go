// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"errors"

	"github.com/rasorp/attila/internal/server/state"
	"github.com/rasorp/attila/internal/state/file"
	"github.com/rasorp/attila/internal/state/mem"
)

func NewBackend(cfg *Config) (state.State, error) {

	if cfg.Memory.Enabled() {
		return mem.New()
	}

	if cfg.File.Enabled() {
		return file.New(cfg.File.Path)
	}

	return nil, errors.New("no state backend configured")
}
