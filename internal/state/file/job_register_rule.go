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

type JobRegisterRule struct {
	store *Store
}

func (j *JobRegisterRule) Create(req *state.JobRegisterRuleCreateReq) (*state.JobRegisterRuleCreateResp, *state.ErrorResp) {
	j.store.lock.Lock()
	defer j.store.lock.Unlock()

	path := filepath.Join(j.store.jobRegRuleDir, req.Rule.Name+".json")

	if code, err := createStoreFile(path, req.Rule); err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("state: %w", err), code)
	}

	return &state.JobRegisterRuleCreateResp{Rule: req.Rule}, nil
}

func (j *JobRegisterRule) Delete(req *state.JobRegisterRuleDeleteReq) (*state.JobRegisterRuleDeleteResp, *state.ErrorResp) {
	j.store.lock.Lock()
	defer j.store.lock.Unlock()

	path := filepath.Join(j.store.jobRegRuleDir, req.Name+".json")

	if err := os.Remove(path); err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("state: %w", err), 500)
	}

	return &state.JobRegisterRuleDeleteResp{}, nil
}

func (j *JobRegisterRule) Get(req *state.JobRegisterRuleGetReq) (*state.JobRegisterRuleGetResp, *state.ErrorResp) {
	j.store.lock.RLock()
	defer j.store.lock.RUnlock()

	path := filepath.Join(j.store.jobRegRuleDir, req.Name+".json")

	var decodedRule state.JobRegisterRule

	if code, err := getStoreFile(path, &decodedRule); err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("state: %w", err), code)
	}

	return &state.JobRegisterRuleGetResp{Rule: &decodedRule}, nil
}

func (j *JobRegisterRule) List(_ *state.JobRegisterRuleListReq) (*state.JobRegisterRuleListResp, *state.ErrorResp) {
	j.store.lock.RLock()
	defer j.store.lock.RUnlock()

	var resp state.JobRegisterRuleListResp

	err := listStoreFiles(j.store.jobRegRuleDir, func(bytes []byte) error {

		var decodedRule state.JobRegisterRule

		if err := json.Unmarshal(bytes, &decodedRule); err != nil {
			return err
		}

		resp.Rules = append(resp.Rules, &decodedRule)
		return nil
	})

	if err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("state: %w", err), 500)
	}
	return &resp, nil
}
