// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package mem

import (
	"github.com/hashicorp/go-memdb"

	"github.com/rasorp/attila/internal/server/state"
)

type Store struct {
	db *memdb.MemDB
}

func New() (state.State, error) {

	db, err := memdb.NewMemDB(newTableSchema())
	if err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

func (s *Store) Region() state.RegionState { return &Region{db: s.db} }

func (s *Store) JobRegister() state.JobRegisterState { return &JobRegister{db: s.db} }

func (s *Store) Name() string { return "mem" }
