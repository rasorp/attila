// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"net/http"
)

const (
	NamespaceDefaultName     = "default"
	NamespaceRegionsWildcard = "*"
)

type Namespace struct {
	Name        string   `hcl:"name" json:"name"`
	Description string   `hcl:"description,optional" json:"description"`
	Regions     []string `hcl:"regions" json:"regions"`
}

type NamespaceStub struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type NamespaceCreateReq struct {
	Namespace *Namespace `json:"namespace"`
}

type NamespaceCreateResp struct {
	Namespace *Namespace `json:"namespace"`
}

type NamespaceListResp struct {
	Namespaces []*NamespaceStub `json:"namespaces"`
}

type NamespaceGetResp struct {
	Namespace *Namespace `json:"namespace"`
}

type Namespaces struct {
	client *Client
}

func (c *Client) Namespaces() *Namespaces { return &Namespaces{client: c} }

func (n *Namespaces) Create(ctx context.Context, req *NamespaceCreateReq) (*NamespaceCreateResp, *Response, error) {

	httpReq, err := n.client.NewRequest(http.MethodPost, "/v1alpha1/namespaces", req)
	if err != nil {
		return nil, nil, err
	}

	var nsResp NamespaceCreateResp

	resp, err := n.client.Do(ctx, httpReq, &nsResp)
	if err != nil {
		return nil, nil, err
	}

	return &nsResp, resp, nil
}

func (n *Namespaces) Delete(ctx context.Context, name string) (*Response, error) {

	req, err := n.client.NewRequest(http.MethodDelete, "/v1alpha1/namespaces/"+name, nil)
	if err != nil {
		return nil, err
	}

	resp, err := n.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (n *Namespaces) Get(ctx context.Context, name string) (*NamespaceGetResp, *Response, error) {

	var nsResp NamespaceGetResp

	req, err := n.client.NewRequest(http.MethodGet, "/v1alpha1/namespaces/"+name, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := n.client.Do(ctx, req, &nsResp)
	if err != nil {
		return nil, resp, err
	}

	return &nsResp, resp, nil
}

func (n *Namespaces) List(ctx context.Context) (*NamespaceListResp, *Response, error) {

	var nsResp NamespaceListResp

	req, err := n.client.NewRequest(http.MethodGet, "/v1alpha1/namespaces", nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := n.client.Do(ctx, req, &nsResp)
	if err != nil {
		return nil, resp, err
	}

	return &nsResp, resp, nil
}
