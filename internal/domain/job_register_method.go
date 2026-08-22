// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package domain

import (
	"errors"
	"fmt"

	jobsdk "github.com/rasorp/attila/pkg/job"
)

type JobRegisterMethod struct {
	Name      string                         `json:"name"`
	Selectors []*jobsdk.MethodSelectorConfig `json:"selectors"`
	Rules     []*JobRegisterMethodRuleLink   `json:"rules"`
	Metadata  *Metadata                      `json:"metadata"`
}

func (m *JobRegisterMethod) Validate() error {
	var errs []error

	if len(m.Rules) < 1 {
		errs = append(errs, errors.New("at least one rule required"))
	} else {
		for i, rule := range m.Rules {
			if err := rule.Validate(); err != nil {
				errs = append(errs, fmt.Errorf("rule %v; %w", i, err))
			}
		}
	}

	// The method needs at least one selector. If there are entries, validate
	// them all, so the operator can make any fix in a single pass.
	if len(m.Selectors) < 1 {
		errs = append(errs, errors.New("at least one selector required"))
	} else {
		for i, selector := range m.Selectors {
			if err := selector.Validate(); err != nil {
				errs = append(errs, fmt.Errorf("selector %v; %w", i, err))
			}
		}
	}

	return errors.Join(errs...)
}

func (m *JobRegisterMethod) Stub() *JobRegisterMethodStub {

	selectors := make([]*JobRegisterMethodSelectorStub, len(m.Selectors))

	for i, selector := range m.Selectors {
		selectors[i] = &JobRegisterMethodSelectorStub{
			Name:     selector.Name,
			Provider: selector.Provider,
		}
	}

	return &JobRegisterMethodStub{
		Name:      m.Name,
		Selectors: selectors,
	}
}

type JobRegisterMethodStub struct {
	Name      string                           `json:"name"`
	Selectors []*JobRegisterMethodSelectorStub `json:"selectors"`
}

type JobRegisterMethodSelectorStub struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
}

type JobRegisterMethodRuleLink struct {
	Name string `json:"name"`
}

func (r *JobRegisterMethodRuleLink) Validate() error {
	var errs []error
	if r.Name == "" {
		errs = append(errs, errors.New("method rule \"name\" cannot be empty"))
	}
	return errors.Join(errs...)
}
