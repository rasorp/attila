// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"net/http"
)

type Region struct {
	Name     string       `hcl:"name" json:"name"`
	Group    string       `hcl:"group,optional" json:"group"`
	Auth     *RegionAuth  `hcl:"auth,block" json:"auth"`
	API      []*RegionAPI `hcl:"api,block" json:"api"`
	TLS      *RegionTLS   `hcl:"tls,optional" json:"tls,omitempty"`
	Metadata *Metadata    `hcl:"metadata" json:"metadata"`
}

type RegionAuth struct {
	Token string `hcl:"token" json:"token"`
}

type RegionAPI struct {
	Address string `hcl:"address" json:"address"`
	Default bool   `hcl:"default,optional" json:"default"`
}

type RegionTLS struct {
	CACert     string `hcl:"ca_cert" json:"ca_cert"`
	ClientCert string `hcl:"client_cert" json:"client_cert"`
	ClientKey  string `hcl:"client_key" json:"client_key"`
	ServerName string `hcl:"server_name" json:"server_name"`
	Insecure   bool   `hcl:"insecure" json:"insecure"`
}

// DefaultOrFirstAddress returns the default API endpoint address if one has
// been configured, or the first one in the array.
func (a *Region) DefaultOrFirstAddress() string {

	for _, apiEndpoint := range a.API {
		if apiEndpoint.Default {
			return apiEndpoint.Address
		}
	}

	// Validation within the Attila server means there should always be one API
	// endpoint detailed in a region object. If this array is empty, the program
	// will panic, which is OK as this is very unexpected behaviour.
	return a.API[0].Address
}

type RegionStub struct {
	Name       string   `json:"name"`
	Group      string   `json:"group"`
	Addresses  []string `json:"addresses"`
	TLSEnabled bool     `json:"tls_enabled"`
}

type RegionCreateReq struct {
	Region *Region `json:"region"`
}

type RegionCreateResp struct {
	Region *Region `json:"region"`
}

type RegionListResp struct {
	Regions []*RegionStub `json:"regions"`
}

type RegionGetResp struct {
	Region *Region `json:"region"`
}

type Regions struct {
	client *Client
}

func (c *Client) Regions() *Regions {
	return &Regions{client: c}
}

func (a *Regions) Create(ctx context.Context, req *RegionCreateReq) (*RegionCreateResp, *Response, error) {

	var regionCreateResp RegionCreateResp

	httpReq, err := a.client.NewRequest(http.MethodPost, "/v1alpha1/regions", req)
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.client.Do(ctx, httpReq, &regionCreateResp)
	if err != nil {
		return nil, nil, err
	}

	return &regionCreateResp, resp, nil
}

func (a *Regions) Delete(ctx context.Context, name string) (*Response, error) {

	req, err := a.client.NewRequest(http.MethodDelete, "/v1alpha1/regions/"+name, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (a *Regions) Get(ctx context.Context, name string) (*RegionGetResp, *Response, error) {

	var regionGetResp RegionGetResp

	req, err := a.client.NewRequest(http.MethodGet, "/v1alpha1/regions/"+name, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.client.Do(ctx, req, &regionGetResp)
	if err != nil {
		return nil, resp, err
	}

	return &regionGetResp, resp, nil
}

func (a *Regions) List(ctx context.Context) (*RegionListResp, *Response, error) {

	var regionListResp RegionListResp

	req, err := a.client.NewRequest(http.MethodGet, "/v1alpha1/regions", nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.client.Do(ctx, req, &regionListResp)
	if err != nil {
		return nil, resp, err
	}

	return &regionListResp, resp, nil
}
