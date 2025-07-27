// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/rasorp/attila/internal/helper/file"
	"github.com/rasorp/attila/internal/server/state"
)

type Store struct {
	dir             string
	jobRegMethodDir string
	jobRegPlanDir   string
	jobRegRuleDir   string
	regionDir       string
	lock            sync.RWMutex
}

const (
	jobRegMethodDir = "job/registration/method"
	jobRegPlanDir   = "job/registration/plan"
	jobRegRuleDir   = "job/registration/rule"
	regionDir       = "region"
)

func New(dir string) (state.State, error) {

	s := Store{
		dir:             dir,
		jobRegMethodDir: filepath.Join(dir, jobRegMethodDir),
		jobRegPlanDir:   filepath.Join(dir, jobRegPlanDir),
		jobRegRuleDir:   filepath.Join(dir, jobRegRuleDir),
		regionDir:       filepath.Join(dir, regionDir),
	}

	for _, subDir := range []string{s.jobRegPlanDir, s.jobRegMethodDir, s.jobRegRuleDir, s.regionDir} {

		// Check the existence of directory. Any error is terminal, except one
		// indicating the directory doesn't exist, as this is normal expected
		// behaviour of a new server.
		fileInfo, err := os.Stat(subDir)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("failed to stat dir: %w", err)
		}

		// Check that any object found looks like a directory.
		if fileInfo != nil && !fileInfo.IsDir() {
			return nil, fmt.Errorf("path %q is file not dir", subDir)
		}

		// Each directory needs rwx permissions rather than just rw, otherwise
		// we will not be able to create subdirectories.
		if err := os.MkdirAll(subDir, 0700); err != nil {
			return nil, fmt.Errorf("failed to create dir: %w", err)
		}
	}

	return &s, nil
}

func (s *Store) JobRegister() state.JobRegisterState { return &JobRegister{store: s} }

func (s *Store) Region() state.RegionState { return &Region{store: s} }

func (s *Store) Name() string { return "file" }

func createStoreFile(path string, data any) (int, error) {

	existing, err := os.Stat(path)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return 500, err
	}
	if existing != nil {
		return 400, errors.New("resource already exists")
	}

	objBytes, err := json.Marshal(data)
	if err != nil {
		return 500, err
	}

	if err := file.AtomicWrite(path, objBytes, 0600); err != nil {
		return 500, err
	}

	return 0, nil
}

func getStoreFile(path string, obj any) (int, error) {

	existing, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return 404, errors.New("resource not found")
		} else {
			return 500, err
		}
	}

	if err := json.Unmarshal(existing, obj); err != nil {
		return 500, err
	}

	return 0, nil
}

func listStoreFiles(path string, decodeFn func([]byte) error) error {
	return filepath.WalkDir(path, func(path string, d fs.DirEntry, e error) error {

		if d == nil {
			return fmt.Errorf("path %q not found", path)
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".json" {
			return nil
		}

		fileContent, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return decodeFn(fileContent)
	})
}
