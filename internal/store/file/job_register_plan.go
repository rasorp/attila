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

type JobRegisterPlan struct {
	store *Store
}

func (j *JobRegisterPlan) Create(req *store.JobRegisterPlanCreateReq) (*store.JobRegisterPlanCreateResp, *store.ErrorResp) {
	j.store.lock.Lock()
	defer j.store.lock.Unlock()

	path := filepath.Join(j.store.jobRegPlanDir, req.Plan.ID.String()+".json")

	if code, err := createStoreFile(path, req.Plan); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("state: %w", err), code)
	}

	return &store.JobRegisterPlanCreateResp{Plan: req.Plan}, nil
}

func (j *JobRegisterPlan) Delete(req *store.JobRegisterPlanDeleteReq) (*store.JobRegisterPlanDeleteResp, *store.ErrorResp) {
	j.store.lock.Lock()
	defer j.store.lock.Unlock()

	path := filepath.Join(j.store.jobRegPlanDir, req.ID.String()+".json")

	if err := os.Remove(path); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("state: %w", err), 500)
	}

	return &store.JobRegisterPlanDeleteResp{}, nil
}

func (j *JobRegisterPlan) Get(req *store.JobRegisterPlanGetReq) (*store.JobRegisterPlanGetResp, *store.ErrorResp) {
	j.store.lock.RLock()
	defer j.store.lock.RUnlock()

	path := filepath.Join(j.store.jobRegPlanDir, req.ID.String()+".json")

	var decodedPlan domain.JobRegisterPlan

	if code, err := getStoreFile(path, &decodedPlan); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("state: %w", err), code)
	}

	return &store.JobRegisterPlanGetResp{Plan: &decodedPlan}, nil
}

func (j *JobRegisterPlan) List(_ *store.JobRegisterPlanListReq) (*store.JobRegisterPlanListResp, *store.ErrorResp) {
	j.store.lock.RLock()
	defer j.store.lock.RUnlock()

	var resp store.JobRegisterPlanListResp

	err := listStoreFiles(j.store.jobRegPlanDir, func(bytes []byte) error {
		var decodedPlan domain.JobRegisterPlan

		if err := json.Unmarshal(bytes, &decodedPlan); err != nil {
			return err
		}

		resp.Plans = append(resp.Plans, &decodedPlan)
		return nil
	})

	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("state: %w", err), 500)
	}
	return &resp, nil
}
