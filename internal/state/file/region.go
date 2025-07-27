// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rasorp/attila/internal/server/state"
)

type Region struct {
	store *Store
}

func (r *Region) Create(req *state.RegionCreateReq) (*state.RegionCreateResp, *state.ErrorResp) {
	r.store.lock.Lock()
	defer r.store.lock.Unlock()

	path := filepath.Join(r.store.regionDir, req.Region.Name+".json")

	if code, err := createStoreFile(path, req.Region); err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("state: %w", err), code)
	}

	return &state.RegionCreateResp{Region: req.Region}, nil
}

func (r *Region) Delete(req *state.RegionDeleteReq) (*state.RegionDeleteResp, *state.ErrorResp) {
	r.store.lock.Lock()
	defer r.store.lock.Unlock()

	path := filepath.Join(r.store.regionDir, req.RegionName+".json")

	if err := os.Remove(path); err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("state: %w", err), 500)
	}

	return &state.RegionDeleteResp{}, nil
}

func (r *Region) Get(req *state.RegionGetReq) (*state.RegionGetResp, *state.ErrorResp) {
	r.store.lock.RLock()
	defer r.store.lock.RUnlock()

	path := filepath.Join(r.store.regionDir, req.RegionName+".json")

	var decodedRegion state.Region

	if code, err := getStoreFile(path, &decodedRegion); err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("state: %w", err), code)
	}

	return &state.RegionGetResp{Region: &decodedRegion}, nil
}

func (r *Region) List(_ *state.RegionListReq) (*state.RegionListResp, *state.ErrorResp) {
	r.store.lock.RLock()
	defer r.store.lock.RUnlock()

	var resp state.RegionListResp

	err := listStoreFiles(r.store.regionDir, func(bytes []byte) error {

		var decodedRegion state.Region

		if err := json.Unmarshal(bytes, &decodedRegion); err != nil {
			return err
		}

		resp.Regions = append(resp.Regions, &decodedRegion)
		return nil
	})

	if err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("state: %w", err), 500)
	}
	return &resp, nil
}
