// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package file

import (
	"testing"

	"github.com/shoenig/test/must"

	"github.com/rasorp/attila/internal/helper/test/mock"
	"github.com/rasorp/attila/internal/server/state"
)

func TestJobRegisterMethod_Create(t *testing.T) {

	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	mockMethod := mock.JobRegistrationMethod()

	createResp1, errResp1 := testState.JobRegister().Method().Create(
		&state.JobRegisterMethodCreateReq{Method: mockMethod},
	)
	must.Nil(t, errResp1)
	must.Eq(t, mockMethod, createResp1.Method)

	createResp2, errResp2 := testState.JobRegister().Method().Create(
		&state.JobRegisterMethodCreateReq{Method: mockMethod},
	)
	must.NotNil(t, errResp2)
	must.Nil(t, createResp2)
}

func TestJobRegisterMethod_Delete(t *testing.T) {

	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	mockMethod := mock.JobRegistrationMethod()

	// Create our test resource.
	createResp1, errResp1 := testState.JobRegister().Method().Create(
		&state.JobRegisterMethodCreateReq{Method: mockMethod},
	)
	must.Nil(t, errResp1)
	must.Eq(t, mockMethod, createResp1.Method)

	// Attempt two delete calls. This first should succeed but the second should
	// fail as the file will no longer exist.
	deleteResp1, errResp1 := testState.JobRegister().Method().Delete(
		&state.JobRegisterMethodDeleteReq{Name: mockMethod.Name},
	)
	must.Nil(t, errResp1)
	must.Eq(t, &state.JobRegisterMethodDeleteResp{}, deleteResp1)

	deleteResp2, errResp2 := testState.JobRegister().Method().Delete(
		&state.JobRegisterMethodDeleteReq{Name: mockMethod.Name},
	)
	must.NotNil(t, errResp2)
	must.Nil(t, deleteResp2)
}

func TestJobRegisterMethod_Get(t *testing.T) {

	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	// Attempt to read a method that does not exist.
	getResp1, err := testState.JobRegister().Method().Get(&state.JobRegisterMethodGetReq{Name: "method"})
	must.Error(t, err)
	must.Nil(t, getResp1)

	mockMethod := mock.JobRegistrationMethod()

	// Create our test resource.
	createResp1, errResp1 := testState.JobRegister().Method().Create(
		&state.JobRegisterMethodCreateReq{Method: mockMethod},
	)
	must.Nil(t, errResp1)
	must.Eq(t, mockMethod, createResp1.Method)

	getResp2, err := testState.JobRegister().Method().Get(&state.JobRegisterMethodGetReq{Name: mockMethod.Name})
	must.Nil(t, err)
	must.Eq(t, mockMethod, getResp2.Method)
}

func TestJobRegisterMethod_List(t *testing.T) {

	testState, err := New(t.TempDir())
	must.NoError(t, err)
	must.NotNil(t, testState)

	// Perform an initial test with no methods in state.
	listResp1, err := testState.JobRegister().Method().List(nil)
	must.Nil(t, err)
	must.Len(t, 0, listResp1.Methods)

	// Create a number of methods that we can list out.
	mockMethods := make([]*state.JobRegisterMethod, 5)

	for i := range mockMethods {
		mockMethods[i] = mock.JobRegistrationMethod()
		createResp, err := testState.JobRegister().Method().Create(
			&state.JobRegisterMethodCreateReq{Method: mockMethods[i]},
		)
		must.Nil(t, err)
		must.NotNil(t, createResp)
	}

	listResp2, err := testState.JobRegister().Method().List(nil)
	must.Nil(t, err)
	must.SliceContainsAll(t, listResp2.Methods, mockMethods)
}
