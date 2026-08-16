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

func TestJobRegisterRule_Create(t *testing.T) {
	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	mockRule := mock.JobRegistrationRule()

	createResp1, errResp1 := testState.JobRegister().Rule().Create(
		&store.JobRegisterRuleCreateReq{Rule: mockRule},
	)
	must.Nil(t, errResp1)
	must.Eq(t, mockRule, createResp1.Rule)

	createResp2, errResp2 := testState.JobRegister().Rule().Create(
		&store.JobRegisterRuleCreateReq{Rule: mockRule},
	)
	must.NotNil(t, errResp2)
	must.Nil(t, createResp2)
}

func TestJobRegisterRule_Delete(t *testing.T) {
	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	mockRule := mock.JobRegistrationRule()

	createResp1, errResp1 := testState.JobRegister().Rule().Create(
		&store.JobRegisterRuleCreateReq{Rule: mockRule},
	)
	must.Nil(t, errResp1)
	must.Eq(t, mockRule, createResp1.Rule)

	deleteResp1, errResp1 := testState.JobRegister().Rule().Delete(
		&store.JobRegisterRuleDeleteReq{Name: mockRule.Name},
	)
	must.Nil(t, errResp1)
	must.Eq(t, &store.JobRegisterRuleDeleteResp{}, deleteResp1)

	deleteResp2, errResp2 := testState.JobRegister().Rule().Delete(
		&store.JobRegisterRuleDeleteReq{Name: mockRule.Name},
	)
	must.NotNil(t, errResp2)
	must.Nil(t, deleteResp2)
}

func TestJobRegisterRule_Get(t *testing.T) {
	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	getResp1, err := testState.JobRegister().Rule().Get(&store.JobRegisterRuleGetReq{Name: "rule"})
	must.Error(t, err)
	must.Nil(t, getResp1)

	mockRule := mock.JobRegistrationRule()

	createResp1, errResp1 := testState.JobRegister().Rule().Create(
		&store.JobRegisterRuleCreateReq{Rule: mockRule},
	)
	must.Nil(t, errResp1)
	must.Eq(t, mockRule, createResp1.Rule)

	getResp2, err := testState.JobRegister().Rule().Get(&store.JobRegisterRuleGetReq{Name: mockRule.Name})
	must.Nil(t, err)
	must.Eq(t, mockRule, getResp2.Rule)
}

func TestJobRegisterRule_List(t *testing.T) {
	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	listResp1, err := testState.JobRegister().Rule().List(nil)
	must.Nil(t, err)
	must.Len(t, 0, listResp1.Rules)

	mockRules := make([]*domain.JobRegisterRule, 5)

	for i := range mockRules {
		mockRules[i] = mock.JobRegistrationRule()
		createResp, err := testState.JobRegister().Rule().Create(
			&store.JobRegisterRuleCreateReq{Rule: mockRules[i]},
		)
		must.Nil(t, err)
		must.NotNil(t, createResp)
	}

	listResp2, err := testState.JobRegister().Rule().List(nil)
	must.Nil(t, err)
	must.SliceContainsAll(t, listResp2.Rules, mockRules)
}
