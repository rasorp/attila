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

type JobRegisterMethod struct {
	store *Store
}

func (j *JobRegisterMethod) Create(req *state.JobRegisterMethodCreateReq) (*state.JobRegisterMethodCreateResp, *state.ErrorResp) {
	j.store.lock.Lock()
	defer j.store.lock.Unlock()

	filePath := filepath.Join(j.store.jobRegMethodDir, req.Method.Name+".json")

	if code, err := createStoreFile(filePath, req.Method); err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("state: %w", err), code)
	}

	return &state.JobRegisterMethodCreateResp{Method: req.Method}, nil
}

func (j *JobRegisterMethod) Delete(req *state.JobRegisterMethodDeleteReq) (*state.JobRegisterMethodDeleteResp, *state.ErrorResp) {
	j.store.lock.Lock()
	defer j.store.lock.Unlock()

	path := filepath.Join(j.store.jobRegMethodDir, req.Name+".json")

	if err := os.Remove(path); err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("state: %w", err), 500)
	}

	return &state.JobRegisterMethodDeleteResp{}, nil
}

func (j *JobRegisterMethod) Get(req *state.JobRegisterMethodGetReq) (*state.JobRegisterMethodGetResp, *state.ErrorResp) {
	j.store.lock.RLock()
	defer j.store.lock.RUnlock()

	path := filepath.Join(j.store.jobRegMethodDir, req.Name+".json")

	var decodedMethod state.JobRegisterMethod

	if code, err := getStoreFile(path, &decodedMethod); err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("state: %w", err), code)
	}

	return &state.JobRegisterMethodGetResp{Method: &decodedMethod}, nil
}

func (j *JobRegisterMethod) List(_ *state.JobRegisterMethodListReq) (*state.JobRegisterMethodListResp, *state.ErrorResp) {
	j.store.lock.RLock()
	defer j.store.lock.RUnlock()

	var resp state.JobRegisterMethodListResp

	err := listStoreFiles(j.store.jobRegMethodDir, func(bytes []byte) error {

		var decodedMethod state.JobRegisterMethod

		if err := json.Unmarshal(bytes, &decodedMethod); err != nil {
			return err
		}

		resp.Methods = append(resp.Methods, &decodedMethod)
		return nil
	})

	if err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("state: %w", err), 500)
	}
	return &resp, nil
}
