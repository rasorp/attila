// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package mem

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-memdb"

	"github.com/rasorp/attila/internal/domain"
)

type Namespace struct {
	db *memdb.MemDB
}

func (ns *Namespace) Create(req *domain.NamespaceCreateReq) (*domain.NamespaceCreateResp, domain.StateError) {
	txn := ns.db.Txn(true)
	defer txn.Abort()

	existingNs, err := txn.First(namespaceTableName, indexID, req.Namespace.Name)
	if err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to read namespace: %w", err), 500)
	}
	if existingNs != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("namespace %q already exists", req.Namespace.Name), 400)
	}

	if err := txn.Insert(namespaceTableName, req.Namespace); err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to create namespace: %w", err), 500)
	}

	txn.Commit()
	return &domain.NamespaceCreateResp{Namespace: req.Namespace}, nil
}

func (ns *Namespace) Delete(req *domain.NamespaceDeleteReq) (*domain.NamespaceDeleteResp, domain.StateError) {
	txn := ns.db.Txn(true)
	defer txn.Abort()

	existingNs, err := txn.First(namespaceTableName, indexID, req.Name)
	if err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to read namespace: %w", err), 500)
	}
	if existingNs == nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("namespace %q not found", req.Name), 404)
	}

	//
	iterMethods, err := txn.Get(jobRegisterMethodTableName, indexNamespace, req.Name)
	if err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to query state: %w", err), 500)
	}
	if iterMethods != nil {
		return nil, domain.NewStateErrorResp(errors.New("namespace not empty"), 400)
	}

	iterRules, err := txn.Get(jobRegisterRuleTableName, indexNamespace, req.Name)
	if err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to query state: %w", err), 500)
	}
	if iterRules.Next() != nil {
		return nil, domain.NewStateErrorResp(errors.New("namespace not empty"), 400)
	}

	iterPlans, err := txn.Get(jobRegisterPlanTableName, "namespace", req.Name)
	if err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to query state: %w", err), 500)
	}
	if iterPlans.Next() != nil {
		return nil, domain.NewStateErrorResp(errors.New("namespace not empty"), 400)
	}

	if err := txn.Delete(namespaceTableName, existingNs); err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to delete namespace: %w", err), 500)
	}

	txn.Commit()
	return &domain.NamespaceDeleteResp{}, nil
}

func (ns *Namespace) Get(req *domain.NamespaceGetReq) (*domain.NamespaceGetResp, domain.StateError) {
	txn := ns.db.Txn(false)
	defer txn.Abort()

	existingNs, err := txn.First(namespaceTableName, indexID, req.Name)
	if err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to read namespace: %w", err), 500)
	}
	if existingNs == nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("namespace %q not found", req.Name), 404)
	}

	txn.Commit()
	return &domain.NamespaceGetResp{Namespace: existingNs.(*domain.Namespace)}, nil
}

func (ns *Namespace) List(req *domain.NamespaceListReq) (*domain.NamespaceListResp, domain.StateError) {
	txn := ns.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get(namespaceTableName, indexID)
	if err != nil {
		return nil, domain.NewStateErrorResp(fmt.Errorf("failed to list namespaces: %w", err), 500)
	}

	var reply domain.NamespaceListResp

	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		reply.Namespaces = append(reply.Namespaces, raw.(*domain.Namespace))
	}

	return &reply, nil
}
