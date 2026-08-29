// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package mem

import (
	"fmt"

	"github.com/hashicorp/go-memdb"

	"github.com/rasorp/attila/internal/domain"
	"github.com/rasorp/attila/internal/store"
)

func (j *JobRegister) Method() store.JobRegisterMethodState { return &JobRegisterMethod{db: j.db} }

type JobRegisterMethod struct {
	db *memdb.MemDB
}

func (j *JobRegisterMethod) Create(req *store.JobRegisterMethodCreateReq) (*store.JobRegisterMethodCreateResp, *store.ErrorResp) {
	txn := j.db.Txn(true)
	defer txn.Abort()

	existingMethod, err := txn.First(jobRegisterMethodTableName, indexID, req.Method.Namespace, req.Method.Name)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to read job register method: %w", err), 500)
	}
	if existingMethod != nil {
		return nil, store.NewErrorResp(fmt.Errorf("job register method %q already exists", req.Method.Name), 400)
	}

	// Inside the same transaction, ensure the namespace exists that the method
	// references.
	ns, err := txn.First(namespaceTableName, indexID, req.Method.Namespace)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to read namespace: %w", err), 500)
	}
	if ns == nil {
		return nil, store.NewErrorResp(fmt.Errorf("namespace %q not found", req.Method.Namespace), 404)
	}

	// Ensure the linked registerment rules exist within state.
	for _, ruleLink := range req.Method.Rules {
		registerRule, err := txn.First(jobRegisterRuleTableName, indexID, req.Method.Namespace, ruleLink.Name)
		if err != nil {
			return nil, store.NewErrorResp(fmt.Errorf("failed to read job register rule: %w", err), 500)
		}
		if registerRule == nil {
			return nil, store.NewErrorResp(fmt.Errorf("job register rule %q not found", ruleLink.Name), 400)
		}
	}

	if err := txn.Insert(jobRegisterMethodTableName, req.Method); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to create job register method: %w", err), 500)
	}

	txn.Commit()
	return &store.JobRegisterMethodCreateResp{Method: req.Method}, nil
}

func (j *JobRegisterMethod) Delete(req *store.JobRegisterMethodDeleteReq) (*store.JobRegisterMethodDeleteResp, *store.ErrorResp) {
	txn := j.db.Txn(true)
	defer txn.Abort()

	existingMethod, err := txn.First(jobRegisterMethodTableName, indexID, req.Namespace, req.Name)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to read job register method: %w", err), 500)
	}
	if existingMethod == nil {
		return nil, store.NewErrorResp(fmt.Errorf("job register method %q not found", req.Name), 404)
	}

	if err := txn.Delete(jobRegisterMethodTableName, existingMethod); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to delete job register method: %w", err), 500)
	}

	txn.Commit()
	return &store.JobRegisterMethodDeleteResp{}, nil
}

func (j *JobRegisterMethod) Get(req *store.JobRegisterMethodGetReq) (*store.JobRegisterMethodGetResp, *store.ErrorResp) {
	txn := j.db.Txn(false)
	defer txn.Abort()

	existingMethod, err := txn.First(jobRegisterMethodTableName, indexID, req.Namespace, req.Name)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to read job register method: %w", err), 500)
	}
	if existingMethod == nil {
		return nil, store.NewErrorResp(fmt.Errorf("job register method %q not found", req.Name), 404)
	}

	txn.Commit()
	return &store.JobRegisterMethodGetResp{Method: existingMethod.(*domain.JobRegisterMethod)}, nil
}

func (j *JobRegisterMethod) List(req *store.JobRegisterMethodListReq) (*store.JobRegisterMethodListResp, *store.ErrorResp) {
	txn := j.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get(jobRegisterMethodTableName, indexID+"_prefix", req.Namespace, "")
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to list job register methods: %w", err), 500)
	}

	var reply store.JobRegisterMethodListResp

	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		reply.Methods = append(reply.Methods, raw.(*domain.JobRegisterMethod))
	}

	return &reply, nil
}
