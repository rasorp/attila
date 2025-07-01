// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package mem

import (
	"fmt"

	"github.com/hashicorp/go-memdb"

	"github.com/rasorp/attila/internal/server/state"
)

func (j *JobRegister) Rule() state.JobRegisterRuleState { return &JobRegisterRule{db: j.db} }

type JobRegisterRule struct {
	db *memdb.MemDB
}

// Create implements state.JobRegisterRuleState.
func (j *JobRegisterRule) Create(req *state.JobRegisterRuleCreateReq) (*state.JobRegisterRuleCreateResp, *state.ErrorResp) {

	txn := j.db.Txn(true)
	defer txn.Abort()

	existingRegion, err := txn.First(jobRegisterRuleTableName, indexID, req.Rule.Name)
	if err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to read job register rule: %w", err), 500)
	}
	if existingRegion != nil {
		return nil, state.NewErrorResp(fmt.Errorf("job register rule %q already exists", req.Rule.Name), 400)
	}

	if err := txn.Insert(jobRegisterRuleTableName, req.Rule); err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to create job register rule: %w", err), 500)
	}

	txn.Commit()
	return &state.JobRegisterRuleCreateResp{Rule: req.Rule}, nil
}

// Delete implements state.JobRegisterRuleState.
func (j *JobRegisterRule) Delete(req *state.JobRegisterRuleDeleteReq) (*state.JobRegisterRuleDeleteResp, *state.ErrorResp) {

	txn := j.db.Txn(true)
	defer txn.Abort()

	existingMethod, err := txn.First(jobRegisterRuleTableName, indexID, req.Name)
	if err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to read job register rule: %w", err), 500)
	}
	if existingMethod == nil {
		return nil, state.NewErrorResp(fmt.Errorf("job register rule %q not found", req.Name), 404)
	}

	if err := txn.Delete(jobRegisterRuleTableName, existingMethod); err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to delete job register rule: %w", err), 500)
	}

	txn.Commit()
	return &state.JobRegisterRuleDeleteResp{}, nil
}

// Get implements state.JobRegisterRuleState.
func (j *JobRegisterRule) Get(req *state.JobRegisterRuleGetReq) (*state.JobRegisterRuleGetResp, *state.ErrorResp) {

	txn := j.db.Txn(false)
	defer txn.Abort()

	existingMethod, err := txn.First(jobRegisterRuleTableName, indexID, req.Name)
	if err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to read job register rule: %w", err), 500)
	}
	if existingMethod == nil {
		return nil, state.NewErrorResp(fmt.Errorf("job register rule %q not found", req.Name), 404)
	}

	txn.Commit()
	return &state.JobRegisterRuleGetResp{Rule: existingMethod.(*state.JobRegisterRule)}, nil
}

// List implements state.JobRegisterRuleState.
func (j *JobRegisterRule) List(req *state.JobRegisterRuleListReq) (*state.JobRegisterRuleListResp, *state.ErrorResp) {

	txn := j.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get(jobRegisterRuleTableName, indexID)
	if err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to list job register rules: %w", err), 500)
	}

	var reply state.JobRegisterRuleListResp

	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		reply.Rules = append(reply.Rules, raw.(*state.JobRegisterRule))
	}

	return &reply, nil
}
