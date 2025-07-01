// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"errors"
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/hashicorp/go-multierror"
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

	var mErr *multierror.Error

	if len(am.Rules) < 1 {
		mErr = multierror.Append(mErr, errors.New("at least one rule required"))
	} else {
		for i, rule := range am.Rules {
			if err := rule.Validate(); err != nil {
				mErr = multierror.Append(mErr, fmt.Errorf("rule %v, %w", i, err))
			}
		}
	}

	if am.Selector != "" {
		if _, err := expr.Compile(am.Selector, expr.AsBool()); err != nil {
			mErr = multierror.Append(mErr, fmt.Errorf("failed to compile expression: %w", err))
		}
	}

	return mErr.ErrorOrNil()
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

	var mErr *multierror.Error

	if arl.Name == "" {
		mErr = multierror.Append(mErr, errors.New("name required"))
	}

	return mErr.ErrorOrNil()
}
