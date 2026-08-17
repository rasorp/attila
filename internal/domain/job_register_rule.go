// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package domain

import (
	"errors"

	jobsdk "github.com/rasorp/attila/pkg/job"
)

type JobRegisterRule struct {
	Name           string                         `json:"name"`
	RegionContexts []JobRegisterRuleRegionContext `json:"region_contexts"`
	RegionPickers  []*jobsdk.RegionPickerConfig   `json:"region_pickers"`
	Metadata       *Metadata                      `json:"metadata"`
}

// Validate performs validation of the job registration rule. It is safe
// to call without checking whether the rule object is nil, although this would
// indicate a serious error in the functionality of the caller.
func (r *JobRegisterRule) Validate() error {
	if r == nil {
		return errors.New("job register rule is empty")
	}

	var errs []error

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
		RegionContexts: r.RegionContexts,
	}
}

type JobRegisterRuleStub struct {
	Name           string                         `json:"name"`
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
