// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package backend

import (
	"errors"

	"github.com/rasorp/attila/internal/store"
	"github.com/rasorp/attila/internal/store/file"
	"github.com/rasorp/attila/internal/store/mem"
)

func NewBackend(cfg *Config) (store.State, error) {
	if cfg.Memory.Enabled() {
		return mem.New()
	}

	if cfg.File.Enabled() {
		return file.New(cfg.File.Path)
	}

	return nil, errors.New("no state backend configured")
}
