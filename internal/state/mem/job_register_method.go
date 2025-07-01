// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package mem

import (
	"fmt"

	"github.com/hashicorp/go-memdb"

	"github.com/rasorp/attila/internal/server/state"
)

func (j *JobRegister) Method() state.JobRegisterMethodState { return &JobRegisterMethod{db: j.db} }

type JobRegisterMethod struct {
	db *memdb.MemDB
}

// Create implements state.JobRegisterMethodState.
func (j *JobRegisterMethod) Create(req *state.JobRegisterMethodCreateReq) (*state.JobRegisterMethodCreateResp, *state.ErrorResp) {

	txn := j.db.Txn(true)
	defer txn.Abort()

	existingRegion, err := txn.First(jobRegisterMethodTableName, indexID, req.Method.Name)
	if err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to read job register method: %w", err), 500)
	}
	if existingRegion != nil {
		return nil, state.NewErrorResp(fmt.Errorf("job register method %q already exists", req.Method.Name), 400)
	}

	// Ensure the linked registerment rules exist within state.
	for _, ruleLink := range req.Method.Rules {
		registerRule, err := txn.First(jobRegisterRuleTableName, indexID, ruleLink.Name)
		if err != nil {
			return nil, state.NewErrorResp(fmt.Errorf("failed to read job register rule: %w", err), 500)
		}
		if registerRule == nil {
			return nil, state.NewErrorResp(fmt.Errorf("job register rule %q not found", ruleLink.Name), 400)
		}
	}

	if err := txn.Insert(jobRegisterMethodTableName, req.Method); err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to create job register method: %w", err), 500)
	}

	txn.Commit()
	return &state.JobRegisterMethodCreateResp{Method: req.Method}, nil
}

// Delete implements state.JobRegisterMethodState.
func (j *JobRegisterMethod) Delete(req *state.JobRegisterMethodDeleteReq) (*state.JobRegisterMethodDeleteResp, *state.ErrorResp) {

	txn := j.db.Txn(true)
	defer txn.Abort()

	existingMethod, err := txn.First(jobRegisterMethodTableName, indexID, req.Name)
	if err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to read job register method: %w", err), 500)
	}
	if existingMethod == nil {
		return nil, state.NewErrorResp(fmt.Errorf("job register method %q not found", req.Name), 404)
	}

	if err := txn.Delete(jobRegisterMethodTableName, existingMethod); err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to delete job register method: %w", err), 500)
	}

	txn.Commit()
	return &state.JobRegisterMethodDeleteResp{}, nil
}

// Get implements state.JobRegisterMethodState.
func (j *JobRegisterMethod) Get(req *state.JobRegisterMethodGetReq) (*state.JobRegisterMethodGetResp, *state.ErrorResp) {

	txn := j.db.Txn(false)
	defer txn.Abort()

	existingMethod, err := txn.First(jobRegisterMethodTableName, indexID, req.Name)
	if err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to read job register method: %w", err), 500)
	}
	if existingMethod == nil {
		return nil, state.NewErrorResp(fmt.Errorf("job register method %q not found", req.Name), 404)
	}

	txn.Commit()
	return &state.JobRegisterMethodGetResp{Method: existingMethod.(*state.JobRegisterMethod)}, nil
}

// List implements state.JobRegisterMethodState.
func (j *JobRegisterMethod) List(*state.JobRegisterMethodListReq) (*state.JobRegisterMethodListResp, *state.ErrorResp) {

	txn := j.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get(jobRegisterMethodTableName, indexID)
	if err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to list job register methods: %w", err), 500)
	}

	var reply state.JobRegisterMethodListResp

	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		reply.Methods = append(reply.Methods, raw.(*state.JobRegisterMethod))
	}

	return &reply, nil
}
