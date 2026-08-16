// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rasorp/attila/internal/domain"
	"github.com/rasorp/attila/internal/store"
)

type Region struct {
	store *Store
}

func (r *Region) Create(req *store.RegionCreateReq) (*store.RegionCreateResp, *store.ErrorResp) {
	r.store.lock.Lock()
	defer r.store.lock.Unlock()

	path := filepath.Join(r.store.regionDir, req.Region.Name+".json")

	if code, err := createStoreFile(path, req.Region); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("state: %w", err), code)
	}

	return &store.RegionCreateResp{Region: req.Region}, nil
}

func (r *Region) Delete(req *store.RegionDeleteReq) (*store.RegionDeleteResp, *store.ErrorResp) {
	r.store.lock.Lock()
	defer r.store.lock.Unlock()

	path := filepath.Join(r.store.regionDir, req.RegionName+".json")

	if err := os.Remove(path); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("state: %w", err), 500)
	}

	return &store.RegionDeleteResp{}, nil
}

func (r *Region) Get(req *store.RegionGetReq) (*store.RegionGetResp, *store.ErrorResp) {
	r.store.lock.RLock()
	defer r.store.lock.RUnlock()

	path := filepath.Join(r.store.regionDir, req.RegionName+".json")

	var decodedRegion domain.Region

	if code, err := getStoreFile(path, &decodedRegion); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("state: %w", err), code)
	}

	return &store.RegionGetResp{Region: &decodedRegion}, nil
}

func (r *Region) List(_ *store.RegionListReq) (*store.RegionListResp, *store.ErrorResp) {
	r.store.lock.RLock()
	defer r.store.lock.RUnlock()

	var resp store.RegionListResp

	err := listStoreFiles(r.store.regionDir, func(bytes []byte) error {
		var decodedRegion domain.Region

		if err := json.Unmarshal(bytes, &decodedRegion); err != nil {
			return err
		}

		resp.Regions = append(resp.Regions, &decodedRegion)
		return nil
	})

	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("state: %w", err), 500)
	}
	return &resp, nil
}
