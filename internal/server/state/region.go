// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-set/v3"
	"github.com/hashicorp/nomad/api"
)

type RegionState interface {
	Create(*RegionCreateReq) (*RegionCreateResp, *ErrorResp)
	Delete(*RegionDeleteReq) (*RegionDeleteResp, *ErrorResp)
	Get(*RegionGetReq) (*RegionGetResp, *ErrorResp)
	List(*RegionListReq) (*RegionListResp, *ErrorResp)
}

type RegionCreateReq struct {
	Region *Region
}

type RegionCreateResp struct {
	Region *Region `json:"region"`
}

type RegionDeleteReq struct {
	RegionName string
}

type RegionDeleteResp struct{}

type RegionGetReq struct {
	RegionName string
}

type RegionGetResp struct {
	Region *Region `json:"region"`
}

type RegionListReq struct{}

type RegionListResp struct {
	Regions []*Region `json:"regions"`
}

type Region struct {
	Name     string       `json:"name"`
	Group    string       `json:"group"`
	Auth     *RegionAuth  `json:"auth,omitempty"`
	API      []*RegionAPI `json:"api,omitempty"`
	TLS      *RegionTLS   `json:"tls,omitempty"`
	Metadata *Metadata    `json:"metadata"`
}

type RegionAuth struct {
	Token string `hcl:"token" json:"token"`
}

type RegionAPI struct {
	Address string `json:"address"`
	Default bool   `json:"default"`
}

type RegionTLS struct {
	CACert     string `json:"ca_cert"`
	ClientCert string `json:"client_cert"`
	ClientKey  string `json:"client_key"`
	ServerName string `json:"server_name"`
	Insecure   bool   `json:"insecure"`
}

// DefaultOrFirstAddress returns the default API endpoint address if one has
// been configured, or the first one in the array.
func (r *Region) DefaultOrFirstAddress() string {

	for _, apiEndpoint := range r.API {
		if apiEndpoint.Default {
			return apiEndpoint.Address
		}
	}

	// Validation within the Attila server means there should always be one API
	// endpoint detailed in a region object. If this array is empty, the program
	// will panic, which is OK as this is very unexpected behaviour.
	return r.API[0].Address
}

func (r *Region) SetDefaults() {
	if r.Group == "" {
		r.Group = "default"
	}
}

func (r *Region) Validate() error {

	var (
		numDefault int
		mErr       *multierror.Error
	)

	addrSet := set.New[string](0)

	for _, apiEndpoint := range r.API {

		// If the API endpoint is set to the default, increment our counter, so
		// we can indentify how many defaults have been set.
		if apiEndpoint.Default {
			numDefault++
		}

		// Inset the address into our set. A false response indicates the
		// address is already held in the set and therefore is a duplicate.
		if !addrSet.Insert(apiEndpoint.Address) {
			mErr = multierror.Append(mErr,
				fmt.Errorf("duplicate API address found: %q", apiEndpoint.Address))
		}
	}

	// The operator must set at least one Nomad API endpoint to talk, otherwise
	// there is no way to talk to the reg
	if addrSet.Size() < 1 {
		mErr = multierror.Append(mErr, errors.New("API list must have at least one entry"))
	}

	// Add an error if the region has more than one API endpoint configured to
	// be the default.
	if numDefault > 1 {
		mErr = multierror.Append(mErr, errors.New("API list can only have one default"))
	}

	if r.Group == "" {
		mErr = multierror.Append(mErr, errors.New("group cannot be empty"))
	}

	return mErr.ErrorOrNil()
}

// GenerateNomadClient is used to generate a new Nomad API client. The Nomad API
// performs its own validation which we do not duplicate in Attila which should
// be taken into account and used appropriately and in conjunction with the
// Validate function.
func (r *Region) GenerateNomadClient() (*api.Client, error) {

	cfg := api.Config{
		Address: r.DefaultOrFirstAddress(),
	}

	// If the region has authentication enabled and the token is not an empty
	// string, we assume it's an ACL token secret ID to use.
	//
	// In the future, we will want to offer the use of identities, so these
	// tokens are not hard-coded.
	if r.Auth != nil && r.Auth.Token != "" {
		cfg.SecretID = r.Auth.Token
	}

	return api.NewClient(&cfg)
}

func (r *Region) Stub() *RegionStub {

	addrs := make([]string, len(r.API))
	for i, addr := range r.API {
		addrs[i] = addr.Address
	}

	return &RegionStub{
		Name:       r.Name,
		Group:      r.Group,
		Addresses:  addrs,
		TLSEnabled: r.TLS != nil,
	}
}

type RegionStub struct {
	Name       string   `json:"name"`
	Group      string   `json:"group"`
	Addresses  []string `json:"addresses"`
	TLSEnabled bool     `json:"tls_enabled"`
}
