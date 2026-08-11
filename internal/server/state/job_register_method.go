// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"errors"
	"fmt"

	"github.com/expr-lang/expr"
)

type JobRegisterMethodState interface {
	Create(*JobRegisterMethodCreateReq) (*JobRegisterMethodCreateResp, *ErrorResp)
	Delete(*JobRegisterMethodDeleteReq) (*JobRegisterMethodDeleteResp, *ErrorResp)
	Get(*JobRegisterMethodGetReq) (*JobRegisterMethodGetResp, *ErrorResp)
	List(*JobRegisterMethodListReq) (*JobRegisterMethodListResp, *ErrorResp)
}

type JobRegisterMethodCreateReq struct {
	Method *JobRegisterMethod `json:"method"`
}

type JobRegisterMethodCreateResp struct {
	Method *JobRegisterMethod `json:"method"`
}

type JobRegisterMethodDeleteReq struct {
	Name string `json:"name"`
}

type JobRegisterMethodDeleteResp struct{}

type JobRegisterMethodGetReq struct {
	Name string `json:"name"`
}

type JobRegisterMethodGetResp struct {
	Method *JobRegisterMethod `json:"method"`
}

type JobRegisterMethodListReq struct{}

type JobRegisterMethodListResp struct {
	Methods []*JobRegisterMethod `json:"methods"`
}

type JobRegisterMethod struct {
	Name     string                       `json:"name"`
	Selector string                       `json:"selector"`
	Rules    []*JobRegisterMethodRuleLink `json:"rule"`
	Metadata *Metadata                    `json:"metadata"`
}

func (am *JobRegisterMethod) Validate() error {

	var errs []error

	if len(am.Rules) < 1 {
		errs = append(errs, errors.New("at least one rule required"))
	} else {
		for i, rule := range am.Rules {
			if err := rule.Validate(); err != nil {
				errs = append(errs, fmt.Errorf("rule %v, %w", i, err))
			}
		}
	}

	if am.Selector != "" {
		if _, err := expr.Compile(am.Selector, expr.AsBool()); err != nil {
			errs = append(errs, fmt.Errorf("failed to compile expression: %w", err))
		}
	}

	return errors.Join(errs...)
}

func (am *JobRegisterMethod) Stub() *JobRegisterMethodStub {
	return &JobRegisterMethodStub{
		Name:     am.Name,
		Selector: am.Selector,
	}
}

type JobRegisterMethodStub struct {
	Name     string `json:"name"`
	Selector string `json:"selector"`
}

type JobRegisterMethodRuleLink struct {
	Name string `json:"name"`
}

func (arl *JobRegisterMethodRuleLink) Validate() error {

	var errs []error

	if arl.Name == "" {
		errs = append(errs, errors.New("name required"))
	}

	return errors.Join(errs...)
}
