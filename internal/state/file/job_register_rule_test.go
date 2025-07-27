// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package file

import (
	"testing"

	"github.com/shoenig/test/must"

	"github.com/rasorp/attila/internal/helper/test/mock"
	"github.com/rasorp/attila/internal/server/state"
)

func TestJobRegisterRule_Create(t *testing.T) {

	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	mockRule := mock.JobRegistrationRule()

	createResp1, errResp1 := testState.JobRegister().Rule().Create(
		&state.JobRegisterRuleCreateReq{Rule: mockRule},
	)
	must.Nil(t, errResp1)
	must.Eq(t, mockRule, createResp1.Rule)

	createResp2, errResp2 := testState.JobRegister().Rule().Create(
		&state.JobRegisterRuleCreateReq{Rule: mockRule},
	)
	must.NotNil(t, errResp2)
	must.Nil(t, createResp2)
}

func TestJobRegisterRule_Delete(t *testing.T) {

	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	mockRule := mock.JobRegistrationRule()

	// Create our test resource.
	createResp1, errResp1 := testState.JobRegister().Rule().Create(
		&state.JobRegisterRuleCreateReq{Rule: mockRule},
	)
	must.Nil(t, errResp1)
	must.Eq(t, mockRule, createResp1.Rule)

	// Attempt two delete calls. This first should succeed but the second should
	// fail as the file will no longer exist.
	deleteResp1, errResp1 := testState.JobRegister().Rule().Delete(
		&state.JobRegisterRuleDeleteReq{Name: mockRule.Name},
	)
	must.Nil(t, errResp1)
	must.Eq(t, &state.JobRegisterRuleDeleteResp{}, deleteResp1)

	deleteResp2, errResp2 := testState.JobRegister().Rule().Delete(
		&state.JobRegisterRuleDeleteReq{Name: mockRule.Name},
	)
	must.NotNil(t, errResp2)
	must.Nil(t, deleteResp2)
}

func TestJobRegisterRule_Get(t *testing.T) {

	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	// Attempt to read a rule that does not exist.
	getResp1, err := testState.JobRegister().Rule().Get(&state.JobRegisterRuleGetReq{Name: "rule"})
	must.Error(t, err)
	must.Nil(t, getResp1)

	mockRule := mock.JobRegistrationRule()

	// Create our test resource.
	createResp1, errResp1 := testState.JobRegister().Rule().Create(
		&state.JobRegisterRuleCreateReq{Rule: mockRule},
	)
	must.Nil(t, errResp1)
	must.Eq(t, mockRule, createResp1.Rule)

	getResp2, err := testState.JobRegister().Rule().Get(&state.JobRegisterRuleGetReq{Name: mockRule.Name})
	must.Nil(t, err)
	must.Eq(t, mockRule, getResp2.Rule)
}

func TestJobRegisterRule_List(t *testing.T) {

	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	// Perform an initial test with no rules in state.
	listResp1, err := testState.JobRegister().Rule().List(nil)
	must.Nil(t, err)
	must.Len(t, 0, listResp1.Rules)

	// Create a number of rules that we can list out.
	mockRules := make([]*state.JobRegisterRule, 5)

	for i := range mockRules {
		mockRules[i] = mock.JobRegistrationRule()
		createResp, err := testState.JobRegister().Rule().Create(
			&state.JobRegisterRuleCreateReq{Rule: mockRules[i]},
		)
		must.Nil(t, err)
		must.NotNil(t, createResp)
	}

	listResp2, err := testState.JobRegister().Rule().List(nil)
	must.Nil(t, err)
	must.SliceContainsAll(t, listResp2.Rules, mockRules)
}
