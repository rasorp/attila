// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/rasorp/attila/internal/helper/test/mock"
	"github.com/rasorp/attila/internal/server/state"
)

func Test_New(t *testing.T) {

	t.Run("no state", func(t *testing.T) {

		dir := t.TempDir()

		fileStore, err := New(dir)
		must.NoError(t, err)
		must.NotNil(t, fileStore)

		for _, subDir := range []string{jobRegMethodDir, jobRegPlanDir, jobRegRuleDir, regionDir} {

			fileInfo, err := os.Stat(filepath.Join(dir, subDir))
			must.NoError(t, err)
			must.True(t, fileInfo.IsDir())
			must.Eq(t, os.ModeDir|os.FileMode(0700), fileInfo.Mode())
		}
	})

	t.Run("existing state", func(t *testing.T) {

		dir := t.TempDir()

		// Create an initial instance of the file store.
		fileStore, err := New(dir)
		must.NoError(t, err)
		must.NotNil(t, fileStore)

		// Write a region to our state store.
		mockRegion := mock.Region()

		createResp, err := fileStore.Region().Create(&state.RegionCreateReq{Region: mockRegion})
		must.Nil(t, err)
		must.Eq(t, createResp.Region, mockRegion)

		// Create another instance of the file store.
		fileStore, err = New(dir)
		must.NoError(t, err)
		must.NotNil(t, fileStore)

		// Ensure the region is still there.
		getResp, err := fileStore.Region().Get(&state.RegionGetReq{RegionName: mockRegion.Name})
		must.Nil(t, err)
		must.Eq(t, getResp.Region, mockRegion)
	})

	t.Run("irregular dir", func(t *testing.T) {

		exec, err := os.Executable()
		must.NoError(t, err)

		fileStore, err := New(exec)
		must.Error(t, err)
		must.Nil(t, fileStore)
	})
}

func TestStore_Name(t *testing.T) {
	fileStore := &Store{}
	must.Eq(t, "file", fileStore.Name())
}
