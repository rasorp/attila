// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package mem

import (
	"fmt"

	"github.com/hashicorp/go-memdb"

	"github.com/rasorp/attila/internal/domain"
)

type JobRegisterRun struct {
	db *memdb.MemDB
}

func (j *JobRegister) Run() domain.JobRegisterRunState {
	return &JobRegisterRun{db: j.db}
}

func (r *JobRegisterRun) Create(req *domain.JobRegisterRunCreateReq) (*domain.JobRegisterRunCreateResp, domain.StateError) {
	txn := r.db.Txn(true)
	defer txn.Abort()

	if err := txn.Insert(jobRegisterRunTableName, req.Run); err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to create run: %w", err), 500)
	}

	txn.Commit()
	return &domain.JobRegisterRunCreateResp{Run: req.Run}, nil
}

func (j *JobRegisterRun) Delete(req *domain.JobRegisterRunDeleteReq) (*domain.JobRegisterRunDeleteResp, domain.StateError) {
	txn := j.db.Txn(true)
	defer txn.Abort()

	existingRun, err := txn.First(jobRegisterRunTableName, indexID, req.ID)
	if err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to read job register run: %w", err), 500)
	}
	if existingRun == nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("job register run %q not found", req.ID), 404)
	}

	if err := txn.Delete(jobRegisterRunTableName, existingRun); err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to delete job register run: %w", err), 500)
	}

	txn.Commit()
	return &domain.JobRegisterRunDeleteResp{}, nil
}

// Get retrieves a plan run record by its ID and namespace.
func (r *JobRegisterRun) Get(req *domain.JobRegisterRunGetReq) (*domain.JobRegisterRunGetResp, domain.StateError) {
	txn := r.db.Txn(false)
	defer txn.Abort()

	raw, err := txn.First(jobRegisterRunTableName, indexID, req.ID)
	if err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to read plan run: %w", err), 500)
	}
	if raw == nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("plan run %q not found", req.ID.String()), 404)
	}

	run := raw.(*domain.JobRegisterRun)

	if run.Namespace != req.Namespace {
		return nil, domain.NewStateErrorResp(fmt.Errorf("plan run %q not found", req.ID.String()), 404)
	}

	return &domain.JobRegisterRunGetResp{Run: run}, nil
}

// List returns all plan run records for the given namespace.
func (r *JobRegisterRun) List(req *domain.JobRegisterRunListReq) (*domain.JobRegisterRunListResp, domain.StateError) {
	txn := r.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get(jobRegisterRunTableName, indexNamespace, req.Namespace)
	if err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to list plan runs: %w", err), 500)
	}

	var reply domain.JobRegisterRunListResp

	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		reply.Runs = append(reply.Runs, raw.(*domain.JobRegisterRun))
	}

	return &reply, nil
}
