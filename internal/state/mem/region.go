// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package mem

import (
	"fmt"

	"github.com/hashicorp/go-memdb"

	"github.com/rasorp/attila/internal/server/state"
)

type Region struct {
	db *memdb.MemDB
}

func (ar *Region) Create(req *state.RegionCreateReq) (*state.RegionCreateResp, *state.ErrorResp) {

	txn := ar.db.Txn(true)
	defer txn.Abort()

	existingRegion, err := txn.First(regionTableName, indexID, req.Region.Name)
	if err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to read region: %w", err), 500)
	}
	if existingRegion != nil {
		return nil, state.NewErrorResp(fmt.Errorf("region %q already exists", req.Region.Name), 400)
	}

	if err := txn.Insert(regionTableName, req.Region); err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to create region: %w", err), 500)
	}

	txn.Commit()
	return &state.RegionCreateResp{Region: req.Region}, nil
}

func (ar *Region) Delete(req *state.RegionDeleteReq) (*state.RegionDeleteResp, *state.ErrorResp) {

	txn := ar.db.Txn(true)
	defer txn.Abort()

	existingRegion, err := txn.First(regionTableName, indexID, req.RegionName)
	if err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to read region: %w", err), 500)
	}
	if existingRegion == nil {
		return nil, state.NewErrorResp(fmt.Errorf("region %q not found", req.RegionName), 404)
	}

	if err := txn.Delete(regionTableName, existingRegion); err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to delete region: %w", err), 500)
	}

	txn.Commit()
	return &state.RegionDeleteResp{}, nil
}

func (ar *Region) Get(req *state.RegionGetReq) (*state.RegionGetResp, *state.ErrorResp) {

	txn := ar.db.Txn(false)
	defer txn.Abort()

	existingRegion, err := txn.First(regionTableName, indexID, req.RegionName)
	if err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to read region: %w", err), 500)
	}
	if existingRegion == nil {
		return nil, state.NewErrorResp(fmt.Errorf("region %q not found", req.RegionName), 404)
	}

	txn.Commit()
	return &state.RegionGetResp{Region: existingRegion.(*state.Region)}, nil
}

func (ar *Region) List(req *state.RegionListReq) (*state.RegionListResp, *state.ErrorResp) {

	txn := ar.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get(regionTableName, indexID)
	if err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to list regions: %w", err), 500)
	}

	var reply state.RegionListResp

	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		reply.Regions = append(reply.Regions, raw.(*state.Region))
	}

	return &reply, nil
}
