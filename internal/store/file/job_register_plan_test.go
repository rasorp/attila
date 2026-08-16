// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package file

import (
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/shoenig/test/must"

	"github.com/rasorp/attila/internal/domain"
	"github.com/rasorp/attila/internal/helper/test/mock"
	"github.com/rasorp/attila/internal/store"
)

func TestJobRegisterPlan_Create(t *testing.T) {
	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	mockPlan := mock.JobRegistrationPlan()

	createResp1, errResp1 := testState.JobRegister().Plan().Create(
		&store.JobRegisterPlanCreateReq{Plan: mockPlan},
	)
	must.Nil(t, errResp1)
	must.Eq(t, mockPlan, createResp1.Plan)

	createResp2, errResp2 := testState.JobRegister().Plan().Create(
		&store.JobRegisterPlanCreateReq{Plan: mockPlan},
	)
	must.NotNil(t, errResp2)
	must.Nil(t, createResp2)
}

func TestJobRegisterPlan_Delete(t *testing.T) {
	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	mockPlan := mock.JobRegistrationPlan()

	createResp1, errResp1 := testState.JobRegister().Plan().Create(
		&store.JobRegisterPlanCreateReq{Plan: mockPlan},
	)
	must.Nil(t, errResp1)
	must.Eq(t, mockPlan, createResp1.Plan)

	deleteResp1, errResp1 := testState.JobRegister().Plan().Delete(
		&store.JobRegisterPlanDeleteReq{ID: mockPlan.ID},
	)
	must.Nil(t, errResp1)
	must.Eq(t, &store.JobRegisterPlanDeleteResp{}, deleteResp1)

	deleteResp2, errResp2 := testState.JobRegister().Plan().Delete(
		&store.JobRegisterPlanDeleteReq{ID: mockPlan.ID},
	)
	must.NotNil(t, errResp2)
	must.Nil(t, deleteResp2)
}

func TestJobRegisterPlan_Get(t *testing.T) {
	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	getResp1, err := testState.JobRegister().Plan().Get(&store.JobRegisterPlanGetReq{ID: ulid.Make()})
	must.Error(t, err)
	must.Nil(t, getResp1)

	mockPlan := mock.JobRegistrationPlan()

	createResp1, errResp1 := testState.JobRegister().Plan().Create(
		&store.JobRegisterPlanCreateReq{Plan: mockPlan},
	)
	must.Nil(t, errResp1)
	must.Eq(t, mockPlan, createResp1.Plan)

	getResp2, err := testState.JobRegister().Plan().Get(&store.JobRegisterPlanGetReq{ID: mockPlan.ID})
	must.Nil(t, err)
	must.Eq(t, mockPlan, getResp2.Plan)
}

func TestJobRegisterPlan_List(t *testing.T) {
	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	listResp1, err := testState.JobRegister().Plan().List(nil)
	must.Nil(t, err)
	must.Len(t, 0, listResp1.Plans)

	mockPlans := make([]*domain.JobRegisterPlan, 5)

	for i := range mockPlans {
		mockPlans[i] = mock.JobRegistrationPlan()
		createResp, err := testState.JobRegister().Plan().Create(
			&store.JobRegisterPlanCreateReq{Plan: mockPlans[i]},
		)
		must.Nil(t, err)
		must.NotNil(t, createResp)
	}

	listResp2, err := testState.JobRegister().Plan().List(nil)
	must.Nil(t, err)
	must.SliceContainsAll(t, listResp2.Plans, mockPlans)
}
