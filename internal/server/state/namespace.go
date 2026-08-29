// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package state

import "github.com/rasorp/attila/internal/domain"

type Namespace interface {
	Create(*NamespaceCreateReq) (*NamespaceCreateResp, Error)
	Delete(*NamespaceDeleteReq) (*NamespaceDeleteResp, Error)
	Get(*NamespaceGetReq) (*NamespaceGetResp, Error)
	List(*NamespaceListReq) (*NamespaceListResp, Error)
}

type NamespaceCreateReq struct {
	Namespace *domain.Namespace
}

type NamespaceCreateResp struct {
	Namespace *domain.Namespace `json:"namespace"`
}

type NamespaceDeleteReq struct {
	Name string
}

type NamespaceDeleteResp struct{}

type NamespaceGetReq struct {
	Name string
}

type NamespaceGetResp struct {
	Namespace *domain.Namespace `json:"namespace"`
}

type NamespaceListReq struct{}

type NamespaceListResp struct {
	Namespaces []*domain.Namespace `json:"namespaces"`
}
