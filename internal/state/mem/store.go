// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package mem

import (
	"github.com/hashicorp/go-memdb"

	"github.com/rasorp/attila/internal/domain"
)

type Store struct {
	db *memdb.MemDB
}

func New() (domain.State, error) {
	db, err := memdb.NewMemDB(newTableSchema())
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) Namespace() domain.NamespaceState     { return &Namespace{db: s.db} }
func (s *Store) Region() domain.RegionState           { return &Region{db: s.db} }
func (s *Store) JobRegister() domain.JobRegisterState { return &JobRegister{db: s.db} }
func (s *Store) Name() string                         { return "mem" }
