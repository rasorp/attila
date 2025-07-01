// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package mem

import (
	"fmt"

	"github.com/hashicorp/go-memdb"

	"github.com/rasorp/attila/internal/server/state"
)

type JobRegister struct {
	db *memdb.MemDB
}

func (j *JobRegister) Plan() state.JobRegisterPlanState { return &JobRegisterPlan{db: j.db} }

type JobRegisterPlan struct {
	db *memdb.MemDB
}

func (j *JobRegisterPlan) Create(req *state.JobRegisterPlanCreateReq) (*state.JobRegisterPlanCreateResp, *state.ErrorResp) {

	txn := j.db.Txn(true)
	defer txn.Abort()

	if err := txn.Insert(jobRegisterPlanTableName, req.Plan); err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to create job registration plan: %w", err), 500)
	}

	txn.Commit()
	return &state.JobRegisterPlanCreateResp{Plan: req.Plan}, nil
}

func (j *JobRegisterPlan) Delete(req *state.JobRegisterPlanDeleteReq) (*state.JobRegisterPlanDeleteResp, *state.ErrorResp) {

	txn := j.db.Txn(true)
	defer txn.Abort()

	existingPlan, err := txn.First(jobRegisterPlanTableName, indexID, req.ID)
	if err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to read job registration plan: %w", err), 500)
	}
	if existingPlan == nil {
		return nil, state.NewErrorResp(fmt.Errorf("job registration plan %q not found", req.ID), 404)
	}

	if err := txn.Delete(jobRegisterPlanTableName, existingPlan); err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to delete job registration plan: %w", err), 500)
	}

	txn.Commit()
	return &state.JobRegisterPlanDeleteResp{}, nil
}

func (j *JobRegisterPlan) Get(req *state.JobRegisterPlanGetReq) (*state.JobRegisterPlanGetResp, *state.ErrorResp) {

	txn := j.db.Txn(false)
	defer txn.Abort()

	existingPlan, err := txn.First(jobRegisterPlanTableName, indexID, req.ID)
	if err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to read job registration plan: %w", err), 500)
	}
	if existingPlan == nil {
		return nil, state.NewErrorResp(fmt.Errorf("job registration plan %q not found", req.ID.String()), 404)
	}

	txn.Commit()
	return &state.JobRegisterPlanGetResp{Plan: existingPlan.(*state.JobRegisterPlan)}, nil
}

func (j *JobRegisterPlan) List(req *state.JobRegisterPlanListReq) (*state.JobRegisterPlanListResp, *state.ErrorResp) {

	txn := j.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get(jobRegisterPlanTableName, indexID)
	if err != nil {
		return nil, state.NewErrorResp(fmt.Errorf("failed to list job registration plans: %w", err), 500)
	}

	var reply state.JobRegisterPlanListResp

	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		reply.Plans = append(reply.Plans, raw.(*state.JobRegisterPlan))
	}

	return &reply, nil
}
