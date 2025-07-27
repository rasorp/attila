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

type JobRegisterPlan struct {
	store *Store
}

func (j *JobRegisterPlan) Create(req *state.JobRegisterPlanCreateReq) (*state.JobRegisterPlanCreateResp, *state.ErrorResp) {
	j.store.lock.Lock()
	defer j.store.lock.Unlock()

	path := filepath.Join(j.store.jobRegPlanDir, req.Plan.ID.String()+".json")

	if code, err := createStoreFile(path, req.Plan); err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("state: %w", err), code)
	}

	return &state.JobRegisterPlanCreateResp{Plan: req.Plan}, nil
}

func (j *JobRegisterPlan) Delete(req *state.JobRegisterPlanDeleteReq) (*state.JobRegisterPlanDeleteResp, *state.ErrorResp) {
	j.store.lock.Lock()
	defer j.store.lock.Unlock()

	path := filepath.Join(j.store.jobRegPlanDir, req.ID.String()+".json")

	if err := os.Remove(path); err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("state: %w", err), 500)
	}

	return &state.JobRegisterPlanDeleteResp{}, nil
}

func (j *JobRegisterPlan) Get(req *state.JobRegisterPlanGetReq) (*state.JobRegisterPlanGetResp, *state.ErrorResp) {
	j.store.lock.RLock()
	defer j.store.lock.RUnlock()

	path := filepath.Join(j.store.jobRegPlanDir, req.ID.String()+".json")

	var decodedPlan state.JobRegisterPlan

	if code, err := getStoreFile(path, &decodedPlan); err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("state: %w", err), code)
	}

	return &state.JobRegisterPlanGetResp{Plan: &decodedPlan}, nil
}

func (j *JobRegisterPlan) List(_ *state.JobRegisterPlanListReq) (*state.JobRegisterPlanListResp, *state.ErrorResp) {
	j.store.lock.RLock()
	defer j.store.lock.RUnlock()

	var resp state.JobRegisterPlanListResp

	err := listStoreFiles(j.store.jobRegPlanDir, func(bytes []byte) error {

		var decodedPlan state.JobRegisterPlan

		if err := json.Unmarshal(bytes, &decodedPlan); err != nil {
			return err
		}

		resp.Plans = append(resp.Plans, &decodedPlan)
		return nil
	})

	if err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("state: %w", err), 500)
	}
	return &resp, nil
}
