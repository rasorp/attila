// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package mem

import (
	"fmt"

	"github.com/hashicorp/go-memdb"

	"github.com/rasorp/attila/internal/domain"
)

func (j *JobRegister) Rule() domain.JobRegisterRuleState { return &JobRegisterRule{db: j.db} }

type JobRegisterRule struct {
	db *memdb.MemDB
}

func (j *JobRegisterRule) Create(req *domain.JobRegisterRuleCreateReq) (*domain.JobRegisterRuleCreateResp, domain.StateError) {
	txn := j.db.Txn(true)
	defer txn.Abort()

	existingRule, err := txn.First(jobRegisterRuleTableName, indexID, req.Rule.Namespace, req.Rule.Name)
	if err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to read job register rule: %w", err), 500)
	}
	if existingRule != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("job register rule %q already exists", req.Rule.Name), 400)
	}

	// Inside the same transaction, ensure the namespace exists that the rule
	// references.
	ns, err := txn.First(namespaceTableName, indexID, req.Rule.Namespace)
	if err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to read namespace: %w", err), 500)
	}
	if ns == nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("namespace %q not found", req.Rule.Namespace), 404)
	}

	if err := txn.Insert(jobRegisterRuleTableName, req.Rule); err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to create job register rule: %w", err), 500)
	}

	txn.Commit()
	return &domain.JobRegisterRuleCreateResp{Rule: req.Rule}, nil
}

func (j *JobRegisterRule) Delete(req *domain.JobRegisterRuleDeleteReq) (*domain.JobRegisterRuleDeleteResp, domain.StateError) {
	txn := j.db.Txn(true)
	defer txn.Abort()

	existingRule, err := txn.First(jobRegisterRuleTableName, indexID, req.Namespace, req.Name)
	if err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to read job register rule: %w", err), 500)
	}
	if existingRule == nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("job register rule %q not found", req.Name), 404)
	}

	if err := txn.Delete(jobRegisterRuleTableName, existingRule); err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to delete job register rule: %w", err), 500)
	}

	txn.Commit()
	return &domain.JobRegisterRuleDeleteResp{}, nil
}

func (j *JobRegisterRule) Get(req *domain.JobRegisterRuleGetReq) (*domain.JobRegisterRuleGetResp, domain.StateError) {
	txn := j.db.Txn(false)
	defer txn.Abort()

	existingRule, err := txn.First(jobRegisterRuleTableName, indexID, req.Namespace, req.Name)
	if err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to read job register rule: %w", err), 500)
	}
	if existingRule == nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("job register rule %q not found", req.Name), 404)
	}

	txn.Commit()
	return &domain.JobRegisterRuleGetResp{Rule: existingRule.(*domain.JobRegisterRule)}, nil
}

func (j *JobRegisterRule) List(req *domain.JobRegisterRuleListReq) (*domain.JobRegisterRuleListResp, domain.StateError) {
	txn := j.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get(jobRegisterRuleTableName, "namespace", req.Namespace)
	if err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to list job register rules: %w", err), 500)
	}

	var reply domain.JobRegisterRuleListResp

	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		reply.Rules = append(reply.Rules, raw.(*domain.JobRegisterRule))
	}

	return &reply, nil
}
