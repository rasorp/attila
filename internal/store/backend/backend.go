// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package backend

import (
	"errors"

	"github.com/rasorp/attila/internal/store"
	"github.com/rasorp/attila/internal/store/mem"
)

func New(cfg *Config) (store.State, error) {
	switch cfg.Provider {
	case ProviderMemory:
		return mem.New()
	default:
		return nil, errors.New("no state backend configured")
	}
}
