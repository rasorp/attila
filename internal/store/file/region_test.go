// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package file

import (
	"testing"

	"github.com/shoenig/test/must"

	"github.com/rasorp/attila/internal/domain"
	"github.com/rasorp/attila/internal/helper/test/mock"
	"github.com/rasorp/attila/internal/store"
)

func TestRegion_Create(t *testing.T) {
	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	mockRegion := mock.Region()

	createResp1, errResp1 := testState.Region().Create(&store.RegionCreateReq{Region: mockRegion})
	must.Nil(t, errResp1)
	must.Eq(t, mockRegion, createResp1.Region)

	createResp2, errResp2 := testState.Region().Create(&store.RegionCreateReq{Region: mockRegion})
	must.NotNil(t, errResp2)
	must.Nil(t, createResp2)
}

func TestRegion_Delete(t *testing.T) {
	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	mockRegion := mock.Region()

	createResp1, errResp1 := testState.Region().Create(&store.RegionCreateReq{Region: mockRegion})
	must.Nil(t, errResp1)
	must.Eq(t, mockRegion, createResp1.Region)

	deleteResp1, errResp1 := testState.Region().Delete(
		&store.RegionDeleteReq{RegionName: mockRegion.Name},
	)
	must.Nil(t, errResp1)
	must.Eq(t, &store.RegionDeleteResp{}, deleteResp1)

	deleteResp2, errResp2 := testState.Region().Delete(
		&store.RegionDeleteReq{RegionName: mockRegion.Name},
	)
	must.NotNil(t, errResp2)
	must.Nil(t, deleteResp2)
}

func TestRegion_Get(t *testing.T) {
	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	getResp1, err := testState.Region().Get(&store.RegionGetReq{RegionName: "region"})
	must.Error(t, err)
	must.Nil(t, getResp1)

	mockRegion := mock.Region()

	createResp1, errResp1 := testState.Region().Create(&store.RegionCreateReq{Region: mockRegion})
	must.Nil(t, errResp1)
	must.Eq(t, mockRegion, createResp1.Region)

	getResp2, err := testState.Region().Get(&store.RegionGetReq{RegionName: mockRegion.Name})
	must.Nil(t, err)
	must.Eq(t, mockRegion, getResp2.Region)
}

func TestRegion_List(t *testing.T) {
	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	listResp1, err := testState.Region().List(nil)
	must.Nil(t, err)
	must.Len(t, 0, listResp1.Regions)

	mockRegions := make([]*domain.Region, 5)

	for i := range mockRegions {
		mockRegions[i] = mock.Region()
		createResp, err := testState.Region().Create(&store.RegionCreateReq{Region: mockRegions[i]})
		must.Nil(t, err)
		must.NotNil(t, createResp)
	}

	listResp2, err := testState.Region().List(nil)
	must.Nil(t, err)
	must.SliceContainsAll(t, listResp2.Regions, mockRegions)
}
