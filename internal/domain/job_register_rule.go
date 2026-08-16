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

type JobRegisterRuleRegionContext string

const (
	JobRegisterRuleContextNamespace JobRegisterRuleRegionContext = "namespace"
	JobRegisterRuleContextNodePool  JobRegisterRuleRegionContext = "node-pool"
)
