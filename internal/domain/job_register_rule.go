// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package domain

import (
	"errors"

	jobsdk "github.com/rasorp/attila/pkg/job"
)

type JobRegisterRuleState interface {
	Create(*JobRegisterRuleCreateReq) (*JobRegisterRuleCreateResp, StateError)
	Delete(*JobRegisterRuleDeleteReq) (*JobRegisterRuleDeleteResp, StateError)
	Get(*JobRegisterRuleGetReq) (*JobRegisterRuleGetResp, StateError)
	List(*JobRegisterRuleListReq) (*JobRegisterRuleListResp, StateError)
}

type JobRegisterRuleCreateReq struct {
	Rule *JobRegisterRule `json:"rule"`
}

type JobRegisterRuleCreateResp struct {
	Rule *JobRegisterRule `json:"rule"`
}

type JobRegisterRuleDeleteReq struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type JobRegisterRuleDeleteResp struct{}

type JobRegisterRuleGetReq struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type JobRegisterRuleGetResp struct {
	Rule *JobRegisterRule `json:"rule"`
}

type JobRegisterRuleListReq struct {
	Namespace string `json:"namespace"`
}

type JobRegisterRuleListResp struct {
	Rules []*JobRegisterRule `json:"rules"`
}

type JobRegisterRule struct {
	Name           string                         `json:"name"`
	Namespace      string                         `json:"namespace"`
	RegionContexts []JobRegisterRuleRegionContext `json:"region_contexts"`
	RegionPickers  []*jobsdk.RegionPickerConfig   `json:"region_pickers"`
	Metadata       *Metadata                      `json:"metadata"`
}

func (r *JobRegisterRule) SetDefaults(reqNamespace string) {
	if r.Namespace == "" {
		r.Namespace = reqNamespace
	}
	r.Metadata = NewMetadata()
}

// Validate performs validation of the job registration rule. It is safe
// to call without checking whether the rule object is nil, although this would
// indicate a serious error in the functionality of the caller.
func (r *JobRegisterRule) Validate() error {
	if r == nil {
		return errors.New("job register rule is empty")
	}

	var errs []error

	if r.Namespace == "" {
		errs = append(errs, errors.New("register rule \"namespace\" cannot be empty"))
	}

	for _, picker := range r.RegionPickers {
		if err := picker.Validate(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (r *JobRegisterRule) Stub() *JobRegisterRuleStub {
	return &JobRegisterRuleStub{
		Name:           r.Name,
		Namespace:      r.Namespace,
		RegionContexts: r.RegionContexts,
	}
}

type JobRegisterRuleStub struct {
	Name           string                         `json:"name"`
	Namespace      string                         `json:"namespace"`
	RegionContexts []JobRegisterRuleRegionContext `json:"region_contexts"`
}

// JobRegisterRuleRegionContext is a way to make additional resource information
// about the region being available to the picker.
type JobRegisterRuleRegionContext struct {

	// Kind is the context that will be made available to the rule picker. It
	// currently supports namespace and node-pool.
	Kind string `json:"kind"`
}

const (
	// JobRegisterRuleContextKindNamespace is the context kind that supplies
	// information from Nomad's "v1/namespaces" endpoint.
	JobRegisterRuleContextKindNamespace = "namespace"

	// JobRegisterRuleContextKindNodepool is the context kind that supplies
	// information from Nomad's "v1/node/pools" endpoint.
	JobRegisterRuleContextKindNodepool = "node-pool"
)
