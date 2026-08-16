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

type JobRegisterMethod struct {
	store *Store
}

func (j *JobRegisterMethod) Create(req *store.JobRegisterMethodCreateReq) (*store.JobRegisterMethodCreateResp, *store.ErrorResp) {
	j.store.lock.Lock()
	defer j.store.lock.Unlock()

	filePath := filepath.Join(j.store.jobRegMethodDir, req.Method.Name+".json")

	if code, err := createStoreFile(filePath, req.Method); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("state: %w", err), code)
	}

	return &store.JobRegisterMethodCreateResp{Method: req.Method}, nil
}

func (j *JobRegisterMethod) Delete(req *store.JobRegisterMethodDeleteReq) (*store.JobRegisterMethodDeleteResp, *store.ErrorResp) {
	j.store.lock.Lock()
	defer j.store.lock.Unlock()

	path := filepath.Join(j.store.jobRegMethodDir, req.Name+".json")

	if err := os.Remove(path); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("state: %w", err), 500)
	}

	return &store.JobRegisterMethodDeleteResp{}, nil
}

func (j *JobRegisterMethod) Get(req *store.JobRegisterMethodGetReq) (*store.JobRegisterMethodGetResp, *store.ErrorResp) {
	j.store.lock.RLock()
	defer j.store.lock.RUnlock()

	path := filepath.Join(j.store.jobRegMethodDir, req.Name+".json")

	var decodedMethod domain.JobRegisterMethod

	if code, err := getStoreFile(path, &decodedMethod); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("state: %w", err), code)
	}

	return &store.JobRegisterMethodGetResp{Method: &decodedMethod}, nil
}

func (j *JobRegisterMethod) List(_ *store.JobRegisterMethodListReq) (*store.JobRegisterMethodListResp, *store.ErrorResp) {
	j.store.lock.RLock()
	defer j.store.lock.RUnlock()

	var resp store.JobRegisterMethodListResp

	err := listStoreFiles(j.store.jobRegMethodDir, func(bytes []byte) error {
		var decodedMethod domain.JobRegisterMethod

		if err := json.Unmarshal(bytes, &decodedMethod); err != nil {
			return err
		}

		resp.Methods = append(resp.Methods, &decodedMethod)
		return nil
	})

	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("state: %w", err), 500)
	}
	return &resp, nil
}
