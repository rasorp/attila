// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package file

import (
	"testing"

	"github.com/shoenig/test/must"

	"github.com/rasorp/attila/internal/helper/test/mock"
	"github.com/rasorp/attila/internal/server/state"
)

func TestRegion_Create(t *testing.T) {

	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	mockRegion := mock.Region()

	createResp1, errResp1 := testState.Region().Create(&state.RegionCreateReq{Region: mockRegion})
	must.Nil(t, errResp1)
	must.Eq(t, mockRegion, createResp1.Region)

	createResp2, errResp2 := testState.Region().Create(&state.RegionCreateReq{Region: mockRegion})
	must.NotNil(t, errResp2)
	must.Nil(t, createResp2)
}

func TestRegion_Delete(t *testing.T) {

	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	mockRegion := mock.Region()

	// Create our test resource.
	createResp1, errResp1 := testState.Region().Create(&state.RegionCreateReq{Region: mockRegion})
	must.Nil(t, errResp1)
	must.Eq(t, mockRegion, createResp1.Region)

	// Attempt two delete calls. This first should succeed but the second should
	// fail as the file will no longer exist.
	deleteResp1, errResp1 := testState.Region().Delete(
		&state.RegionDeleteReq{RegionName: mockRegion.Name},
	)
	must.Nil(t, errResp1)
	must.Eq(t, &state.RegionDeleteResp{}, deleteResp1)

	deleteResp2, errResp2 := testState.Region().Delete(
		&state.RegionDeleteReq{RegionName: mockRegion.Name},
	)
	must.NotNil(t, errResp2)
	must.Nil(t, deleteResp2)
}

func TestRegion_Get(t *testing.T) {

	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	// Attempt to read a region that does not exist.
	getResp1, err := testState.Region().Get(&state.RegionGetReq{RegionName: "region"})
	must.Error(t, err)
	must.Nil(t, getResp1)

	mockRegion := mock.Region()

	// Create our test resource.
	createResp1, errResp1 := testState.Region().Create(&state.RegionCreateReq{Region: mockRegion})
	must.Nil(t, errResp1)
	must.Eq(t, mockRegion, createResp1.Region)

	getResp2, err := testState.Region().Get(&state.RegionGetReq{RegionName: mockRegion.Name})
	must.Nil(t, err)
	must.Eq(t, mockRegion, getResp2.Region)
}

func TestRegion_List(t *testing.T) {

	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	// Perform an initial test with no regions in state.
	listResp1, err := testState.Region().List(nil)
	must.Nil(t, err)
	must.Len(t, 0, listResp1.Regions)

	// Create a number of regions that we can list out.
	mockRegions := make([]*state.Region, 5)

	for i := range mockRegions {
		mockRegions[i] = mock.Region()
		createResp, err := testState.Region().Create(&state.RegionCreateReq{Region: mockRegions[i]})
		must.Nil(t, err)
		must.NotNil(t, createResp)
	}

	listResp2, err := testState.Region().List(nil)
	must.Nil(t, err)
	must.SliceContainsAll(t, listResp2.Regions, mockRegions)
}
