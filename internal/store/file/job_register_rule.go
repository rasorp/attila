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

type JobRegisterRule struct {
	store *Store
}

func (j *JobRegisterRule) Create(req *store.JobRegisterRuleCreateReq) (*store.JobRegisterRuleCreateResp, *store.ErrorResp) {
	j.store.lock.Lock()
	defer j.store.lock.Unlock()

	path := filepath.Join(j.store.jobRegRuleDir, req.Rule.Name+".json")

	if code, err := createStoreFile(path, req.Rule); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("state: %w", err), code)
	}

	return &store.JobRegisterRuleCreateResp{Rule: req.Rule}, nil
}

func (j *JobRegisterRule) Delete(req *store.JobRegisterRuleDeleteReq) (*store.JobRegisterRuleDeleteResp, *store.ErrorResp) {
	j.store.lock.Lock()
	defer j.store.lock.Unlock()

	path := filepath.Join(j.store.jobRegRuleDir, req.Name+".json")

	if err := os.Remove(path); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("state: %w", err), 500)
	}

	return &store.JobRegisterRuleDeleteResp{}, nil
}

func (j *JobRegisterRule) Get(req *store.JobRegisterRuleGetReq) (*store.JobRegisterRuleGetResp, *store.ErrorResp) {
	j.store.lock.RLock()
	defer j.store.lock.RUnlock()

	path := filepath.Join(j.store.jobRegRuleDir, req.Name+".json")

	var decodedRule domain.JobRegisterRule

	if code, err := getStoreFile(path, &decodedRule); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("state: %w", err), code)
	}

	return &store.JobRegisterRuleGetResp{Rule: &decodedRule}, nil
}

func (j *JobRegisterRule) List(_ *store.JobRegisterRuleListReq) (*store.JobRegisterRuleListResp, *store.ErrorResp) {
	j.store.lock.RLock()
	defer j.store.lock.RUnlock()

	var resp store.JobRegisterRuleListResp

	err := listStoreFiles(j.store.jobRegRuleDir, func(bytes []byte) error {
		var decodedRule domain.JobRegisterRule

		if err := json.Unmarshal(bytes, &decodedRule); err != nil {
			return err
		}

		resp.Rules = append(resp.Rules, &decodedRule)
		return nil
	})

	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("state: %w", err), 500)
	}
	return &resp, nil
}
