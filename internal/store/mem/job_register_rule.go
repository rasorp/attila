// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package mem

import (
	"fmt"

	"github.com/hashicorp/go-memdb"

	"github.com/rasorp/attila/internal/domain"
	"github.com/rasorp/attila/internal/store"
)

func (j *JobRegister) Rule() store.JobRegisterRuleState { return &JobRegisterRule{db: j.db} }

type JobRegisterRule struct {
	db *memdb.MemDB
}

func (j *JobRegisterRule) Create(req *store.JobRegisterRuleCreateReq) (*store.JobRegisterRuleCreateResp, *store.ErrorResp) {
	txn := j.db.Txn(true)
	defer txn.Abort()

	existingRegion, err := txn.First(jobRegisterRuleTableName, indexID, req.Rule.Name)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to read job register rule: %w", err), 500)
	}
	if existingRegion != nil {
		return nil, store.NewErrorResp(fmt.Errorf("job register rule %q already exists", req.Rule.Name), 400)
	}

	if err := txn.Insert(jobRegisterRuleTableName, req.Rule); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to create job register rule: %w", err), 500)
	}

	txn.Commit()
	return &store.JobRegisterRuleCreateResp{Rule: req.Rule}, nil
}

func (j *JobRegisterRule) Delete(req *store.JobRegisterRuleDeleteReq) (*store.JobRegisterRuleDeleteResp, *store.ErrorResp) {
	txn := j.db.Txn(true)
	defer txn.Abort()

	existingMethod, err := txn.First(jobRegisterRuleTableName, indexID, req.Name)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to read job register rule: %w", err), 500)
	}
	if existingMethod == nil {
		return nil, store.NewErrorResp(fmt.Errorf("job register rule %q not found", req.Name), 404)
	}

	if err := txn.Delete(jobRegisterRuleTableName, existingMethod); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to delete job register rule: %w", err), 500)
	}

	txn.Commit()
	return &store.JobRegisterRuleDeleteResp{}, nil
}

func (j *JobRegisterRule) Get(req *store.JobRegisterRuleGetReq) (*store.JobRegisterRuleGetResp, *store.ErrorResp) {
	txn := j.db.Txn(false)
	defer txn.Abort()

	existingMethod, err := txn.First(jobRegisterRuleTableName, indexID, req.Name)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to read job register rule: %w", err), 500)
	}
	if existingMethod == nil {
		return nil, store.NewErrorResp(fmt.Errorf("job register rule %q not found", req.Name), 404)
	}

	txn.Commit()
	return &store.JobRegisterRuleGetResp{Rule: existingMethod.(*domain.JobRegisterRule)}, nil
}

func (j *JobRegisterRule) List(req *store.JobRegisterRuleListReq) (*store.JobRegisterRuleListResp, *store.ErrorResp) {
	txn := j.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get(jobRegisterRuleTableName, indexID)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to list job register rules: %w", err), 500)
	}

	var reply store.JobRegisterRuleListResp

	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		reply.Rules = append(reply.Rules, raw.(*domain.JobRegisterRule))
	}

	return &reply, nil
}
