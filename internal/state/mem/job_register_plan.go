// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package mem

import (
	"fmt"

	"github.com/hashicorp/go-memdb"

	"github.com/rasorp/attila/internal/domain"
)

type JobRegister struct {
	db *memdb.MemDB
}

func (j *JobRegister) Plan() domain.JobRegisterPlanState { return &JobRegisterPlan{db: j.db} }

type JobRegisterPlan struct {
	db *memdb.MemDB
}

func (j *JobRegisterPlan) Create(req *domain.JobRegisterPlanCreateReq) (*domain.JobRegisterPlanCreateResp, domain.StateError) {
	txn := j.db.Txn(true)
	defer txn.Abort()

	if err := txn.Insert(jobRegisterPlanTableName, req.Plan); err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to create job registration plan: %w", err), 500)
	}

	txn.Commit()
	return &domain.JobRegisterPlanCreateResp{Plan: req.Plan}, nil
}

func (j *JobRegisterPlan) Delete(req *domain.JobRegisterPlanDeleteReq) (*domain.JobRegisterPlanDeleteResp, domain.StateError) {
	txn := j.db.Txn(true)
	defer txn.Abort()

	existingPlan, err := txn.First(jobRegisterPlanTableName, indexID, req.ID)
	if err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to read job registration plan: %w", err), 500)
	}
	if existingPlan == nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("job registration plan %q not found", req.ID), 404)
	}

	if err := txn.Delete(jobRegisterPlanTableName, existingPlan); err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to delete job registration plan: %w", err), 500)
	}

	txn.Commit()
	return &domain.JobRegisterPlanDeleteResp{}, nil
}

func (j *JobRegisterPlan) Get(req *domain.JobRegisterPlanGetReq) (*domain.JobRegisterPlanGetResp, domain.StateError) {
	txn := j.db.Txn(false)
	defer txn.Abort()

	existingPlanRaw, err := txn.First(jobRegisterPlanTableName, indexID, req.ID)
	if err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to read job registration plan: %w", err), 500)
	}

	existingPlan := existingPlanRaw.(*domain.JobRegisterPlan)

	if existingPlan == nil || existingPlan.Namespace != req.Namespace {
		return nil, domain.NewStateErrorResp(fmt.Errorf("job registration plan %q not found", req.ID.String()), 404)
	}

	txn.Commit()
	return &domain.JobRegisterPlanGetResp{Plan: existingPlan}, nil
}

func (j *JobRegisterPlan) List(req *domain.JobRegisterPlanListReq) (*domain.JobRegisterPlanListResp, domain.StateError) {
	txn := j.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get(jobRegisterPlanTableName, indexID, req.Namespace)
	if err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to list job registration plans: %w", err), 500)
	}

	var reply domain.JobRegisterPlanListResp

	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		reply.Plans = append(reply.Plans, raw.(*domain.JobRegisterPlan))
	}

	return &reply, nil
}
