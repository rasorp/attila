// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package mem

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-memdb"

	"github.com/rasorp/attila/internal/domain"
	"github.com/rasorp/attila/internal/server/state"
	"github.com/rasorp/attila/internal/store"
)

type Namespace struct {
	db *memdb.MemDB
}

func (ns *Namespace) Create(req *state.NamespaceCreateReq) (*state.NamespaceCreateResp, state.Error) {
	txn := ns.db.Txn(true)
	defer txn.Abort()

	existingNs, err := txn.First(namespaceTableName, indexID, req.Namespace.Name)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to read namespace: %w", err), 500)
	}
	if existingNs != nil {
		return nil, store.NewErrorResp(fmt.Errorf("namespace %q already exists", req.Namespace.Name), 400)
	}

	if err := txn.Insert(namespaceTableName, req.Namespace); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to create namespace: %w", err), 500)
	}

	txn.Commit()
	return &state.NamespaceCreateResp{Namespace: req.Namespace}, nil
}

func (ns *Namespace) Delete(req *state.NamespaceDeleteReq) (*state.NamespaceDeleteResp, state.Error) {
	txn := ns.db.Txn(true)
	defer txn.Abort()

	existingNs, err := txn.First(namespaceTableName, indexID, req.Name)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to read namespace: %w", err), 500)
	}
	if existingNs == nil {
		return nil, store.NewErrorResp(fmt.Errorf("namespace %q not found", req.Name), 404)
	}

	//
	iterMethods, err := txn.Get(jobRegisterMethodTableName, indexNamespace, req.Name)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to query state: %w", err), 500)
	}
	if iterMethods != nil {
		return nil, store.NewErrorResp(errors.New("namespace not empty"), 400)
	}

	iterRules, err := txn.Get(jobRegisterRuleTableName, indexNamespace, req.Name)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to query state: %w", err), 500)
	}
	if iterRules.Next() != nil {
		return nil, store.NewErrorResp(errors.New("namespace not empty"), 400)
	}

	iterPlans, err := txn.Get(jobRegisterPlanTableName, "namespace", req.Name)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to query state: %w", err), 500)
	}
	if iterPlans.Next() != nil {
		return nil, store.NewErrorResp(errors.New("namespace not empty"), 400)
	}

	if err := txn.Delete(namespaceTableName, existingNs); err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to delete namespace: %w", err), 500)
	}

	txn.Commit()
	return &state.NamespaceDeleteResp{}, nil
}

func (ns *Namespace) Get(req *state.NamespaceGetReq) (*state.NamespaceGetResp, state.Error) {
	txn := ns.db.Txn(false)
	defer txn.Abort()

	existingNs, err := txn.First(namespaceTableName, indexID, req.Name)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to read namespace: %w", err), 500)
	}
	if existingNs == nil {
		return nil, store.NewErrorResp(fmt.Errorf("namespace %q not found", req.Name), 404)
	}

	txn.Commit()
	return &state.NamespaceGetResp{Namespace: existingNs.(*domain.Namespace)}, nil
}

func (ns *Namespace) List(req *state.NamespaceListReq) (*state.NamespaceListResp, state.Error) {
	txn := ns.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get(namespaceTableName, indexID)
	if err != nil {
		return nil, store.NewErrorResp(fmt.Errorf("failed to list namespaces: %w", err), 500)
	}

	var reply state.NamespaceListResp

	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		reply.Namespaces = append(reply.Namespaces, raw.(*domain.Namespace))
	}

	return &reply, nil
}
