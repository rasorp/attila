// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package backend

import (
	"errors"

	"github.com/rasorp/attila/internal/domain"
	"github.com/rasorp/attila/internal/server/state"
	"github.com/rasorp/attila/internal/store"
	"github.com/rasorp/attila/internal/store/mem"
)

func New(cfg *Config) (store.State, error) {

	var (
		backend store.State
		err     error
	)

	switch cfg.Provider {
	case ProviderMemory:
		backend, err = mem.New()
	default:
		return nil, errors.New("no state backend configured")
	}

	if err := namespaceInit(backend); err != nil {
		return nil, err
	}
	return backend, err
}

func namespaceInit(backend store.State) error {

	_, err := backend.Namespace().Get(&state.NamespaceGetReq{Name: domain.NamespaceDefaultName})
	if err == nil {
		return nil
	}

	if err.StatusCode() == 404 {
		if _, err := backend.Namespace().Create(&state.NamespaceCreateReq{Namespace: domain.DefaultNamespace()}); err == nil {
			return nil
		}
	}

	return err
}
