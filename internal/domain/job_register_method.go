// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package domain

import (
	"errors"
	"fmt"

	"github.com/expr-lang/expr"
)

type JobRegisterMethod struct {
	Name     string                       `json:"name"`
	Selector string                       `json:"selector"`
	Rules    []*JobRegisterMethodRuleLink `json:"rule"`
	Metadata *Metadata                    `json:"metadata"`
}

func (m *JobRegisterMethod) Validate() error {
	var errs []error

	if len(m.Rules) < 1 {
		errs = append(errs, errors.New("at least one rule required"))
	} else {
		for i, rule := range m.Rules {
			if err := rule.Validate(); err != nil {
				errs = append(errs, fmt.Errorf("rule %v, %w", i, err))
			}
		}
	}

	if m.Selector != "" {
		if _, err := expr.Compile(m.Selector, expr.AsBool()); err != nil {
			errs = append(errs, fmt.Errorf("failed to compile expression: %w", err))
		}
	}

	return errors.Join(errs...)
}

func (m *JobRegisterMethod) Stub() *JobRegisterMethodStub {
	return &JobRegisterMethodStub{
		Name:     m.Name,
		Selector: m.Selector,
	}
}

type JobRegisterMethodStub struct {
	Name     string `json:"name"`
	Selector string `json:"selector"`
}

type JobRegisterMethodRuleLink struct {
	Name string `json:"name"`
}

func (r *JobRegisterMethodRuleLink) Validate() error {
	var errs []error
	if r.Name == "" {
		errs = append(errs, errors.New("name required"))
	}
	return errors.Join(errs...)
}
