// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"errors"

	"github.com/rasorp/attila/internal/domain"
	"github.com/rasorp/attila/internal/state/mem"
)

func New(cfg *Config) (domain.State, error) {

	var (
		backend domain.State
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

func namespaceInit(backend domain.State) error {

	_, err := backend.Namespace().Get(&domain.NamespaceGetReq{Name: domain.NamespaceDefaultName})
	if err == nil {
		return nil
	}

	if err.StatusCode() == 404 {
		if _, err := backend.Namespace().Create(&domain.NamespaceCreateReq{Namespace: domain.DefaultNamespace()}); err == nil {
			return nil
		}
	}

	return err
}
