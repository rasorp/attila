// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package mem

import (
	"fmt"

	"github.com/hashicorp/go-memdb"

	"github.com/rasorp/attila/internal/domain"
	"github.com/rasorp/attila/internal/store"
)

type JobRegister struct {
	db *memdb.MemDB
}

func (j *JobRegister) Plan() store.JobRegisterPlanState { return &JobRegisterPlan{db: j.db} }

type JobRegisterPlan struct {
	db *memdb.MemDB
}

func (j *JobRegisterPlan) Create(req *store.JobRegisterPlanCreateReq) (*store.JobRegisterPlanCreateResp, *store.ErrorResp) {
	txn := j.db.Txn(true)
	defer txn.Abort()

	if err := txn.Insert(jobRegisterPlanTableName, req.Plan); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to create job registration plan: %w", err), 500)
	}

	txn.Commit()
	return &store.JobRegisterPlanCreateResp{Plan: req.Plan}, nil
}

func (j *JobRegisterPlan) Delete(req *store.JobRegisterPlanDeleteReq) (*store.JobRegisterPlanDeleteResp, *store.ErrorResp) {
	txn := j.db.Txn(true)
	defer txn.Abort()

	existingPlan, err := txn.First(jobRegisterPlanTableName, indexID, req.ID)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to read job registration plan: %w", err), 500)
	}
	if existingPlan == nil {
		return nil, store.NewErrorResp(fmt.Errorf("job registration plan %q not found", req.ID), 404)
	}

	if err := txn.Delete(jobRegisterPlanTableName, existingPlan); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to delete job registration plan: %w", err), 500)
	}

	txn.Commit()
	return &store.JobRegisterPlanDeleteResp{}, nil
}

func (j *JobRegisterPlan) Get(req *store.JobRegisterPlanGetReq) (*store.JobRegisterPlanGetResp, *store.ErrorResp) {
	txn := j.db.Txn(false)
	defer txn.Abort()

	existingPlan, err := txn.First(jobRegisterPlanTableName, indexID, req.ID)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to read job registration plan: %w", err), 500)
	}
	if existingPlan == nil {
		return nil, store.NewErrorResp(fmt.Errorf("job registration plan %q not found", req.ID.String()), 404)
	}

	txn.Commit()
	return &store.JobRegisterPlanGetResp{Plan: existingPlan.(*domain.JobRegisterPlan)}, nil
}

func (j *JobRegisterPlan) List(req *store.JobRegisterPlanListReq) (*store.JobRegisterPlanListResp, *store.ErrorResp) {
	txn := j.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get(jobRegisterPlanTableName, indexID)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to list job registration plans: %w", err), 500)
	}

	var reply store.JobRegisterPlanListResp

	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		reply.Plans = append(reply.Plans, raw.(*domain.JobRegisterPlan))
	}

	return &reply, nil
}
