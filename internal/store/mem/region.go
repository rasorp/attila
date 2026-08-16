// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package mem

import (
	"fmt"

	"github.com/hashicorp/go-memdb"

	"github.com/rasorp/attila/internal/domain"
	"github.com/rasorp/attila/internal/store"
)

type Region struct {
	db *memdb.MemDB
}

func (ar *Region) Create(req *store.RegionCreateReq) (*store.RegionCreateResp, *store.ErrorResp) {
	txn := ar.db.Txn(true)
	defer txn.Abort()

	existingRegion, err := txn.First(regionTableName, indexID, req.Region.Name)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to read region: %w", err), 500)
	}
	if existingRegion != nil {
		return nil, store.NewErrorResp(fmt.Errorf("region %q already exists", req.Region.Name), 400)
	}

	if err := txn.Insert(regionTableName, req.Region); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to create region: %w", err), 500)
	}

	txn.Commit()
	return &store.RegionCreateResp{Region: req.Region}, nil
}

func (ar *Region) Delete(req *store.RegionDeleteReq) (*store.RegionDeleteResp, *store.ErrorResp) {
	txn := ar.db.Txn(true)
	defer txn.Abort()

	existingRegion, err := txn.First(regionTableName, indexID, req.RegionName)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to read region: %w", err), 500)
	}
	if existingRegion == nil {
		return nil, store.NewErrorResp(fmt.Errorf("region %q not found", req.RegionName), 404)
	}

	if err := txn.Delete(regionTableName, existingRegion); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to delete region: %w", err), 500)
	}

	txn.Commit()
	return &store.RegionDeleteResp{}, nil
}

func (ar *Region) Get(req *store.RegionGetReq) (*store.RegionGetResp, *store.ErrorResp) {
	txn := ar.db.Txn(false)
	defer txn.Abort()

	existingRegion, err := txn.First(regionTableName, indexID, req.RegionName)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to read region: %w", err), 500)
	}
	if existingRegion == nil {
		return nil, store.NewErrorResp(fmt.Errorf("region %q not found", req.RegionName), 404)
	}

	txn.Commit()
	return &store.RegionGetResp{Region: existingRegion.(*domain.Region)}, nil
}

func (ar *Region) List(req *store.RegionListReq) (*store.RegionListResp, *store.ErrorResp) {
	txn := ar.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get(regionTableName, indexID)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to list regions: %w", err), 500)
	}

	var reply store.RegionListResp

	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		reply.Regions = append(reply.Regions, raw.(*domain.Region))
	}

	return &reply, nil
}
