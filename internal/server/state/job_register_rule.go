// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/expr-lang/expr"
	"github.com/hashicorp/go-multierror"
)

type JobRegisterRuleState interface {
	Create(*JobRegisterRuleCreateReq) (*JobRegisterRuleCreateResp, *ErrorResp)
	Delete(*JobRegisterRuleDeleteReq) (*JobRegisterRuleDeleteResp, *ErrorResp)
	Get(*JobRegisterRuleGetReq) (*JobRegisterRuleGetResp, *ErrorResp)
	List(*JobRegisterRuleListReq) (*JobRegisterRuleListResp, *ErrorResp)
}

type JobRegisterRuleCreateReq struct {
	Rule *JobRegisterRule `json:"rule"`
}

type JobRegisterRuleCreateResp struct {
	Rule *JobRegisterRule `json:"rule"`
}

type JobRegisterRuleDeleteReq struct {
	Name string `json:"name"`
}

type JobRegisterRuleDeleteResp struct{}

type JobRegisterRuleGetReq struct {
	Name string `json:"name"`
}

type JobRegisterRuleGetResp struct {
	Rule *JobRegisterRule `json:"rule"`
}

type JobRegisterRuleListReq struct{}

type JobRegisterRuleListResp struct {
	Rules []*JobRegisterRule `json:"rules"`
}

type JobRegisterRule struct {
	Name           string                         `json:"name"`
	RegionContexts []JobRegisterRuleRegionContext `json:"region_contexts"`
	RegionFilter   *JobRegisterRuleFilter         `json:"region_filter"`
	RegionPicker   *JobRegisterRulePicker         `json:"region_picker"`
	Metadata       *Metadata                      `json:"metadata"`
}

// Validate performs validation of the job registration rule. It is safe
// to call without checking whether the rule object is nil, although this would
// indicate a serious error in the functionality of the caller.
func (a *JobRegisterRule) Validate() error {

	// Protect against complete incorrect use which would cause the server to
	// panic. This does not use the multierror because it will be the only error
	// to occur.
	if a == nil {
		return errors.New("job register rule is empty")
	}

	var mErr *multierror.Error

	if err := a.RegionPicker.Validate(); err != nil {
		mErr = multierror.Append(mErr, err)
	}

	return mErr.ErrorOrNil()
}

func (a *JobRegisterRule) Stub() *JobRegisterRuleStub {
	return &JobRegisterRuleStub{
		Name:           a.Name,
		RegionContexts: a.RegionContexts,
	}
}

type JobRegisterRuleStub struct {
	Name           string                         `json:"name"`
	RegionContexts []JobRegisterRuleRegionContext `json:"region_contexts"`
}

type JobRegisterRuleFilter struct {
	Expression *JobRegisterRuleFilterExpression `json:"expression"`
}

type JobRegisterRulePicker struct {
	Expression *JobRegisterRuleFilterExpression `json:"expression"`
}

type JobRegisterRuleFilterExpression struct {
	Selector string `json:"selector"`
}

// Validate performs validation of the job registration rule picker. It is safe
// to call without checking whether the picker object is nil.
func (j *JobRegisterRulePicker) Validate() error {

	// Protect against the rule picker being nil, so callers do not have to do
	// this.
	if j == nil {
		return nil
	}

	// If the expression is nil, we have nothing to valid. This will change if
	// alternate picker functionality is added.
	if j.Expression == nil {
		return nil
	}

	// Ensure the expression compiles with an expected slice result. This is as
	// close as we can get to validation without running the expression against
	// some data.
	if _, err := expr.Compile(j.Expression.Selector, expr.AsKind(reflect.Slice)); err != nil {
		return fmt.Errorf("failed to compile rule picker selector: %w", err)
	}

	return nil
}

type JobRegisterRuleRegionContext string

const (
	JobRegisterRuleContextNamespace JobRegisterRuleRegionContext = "namespace"
	JobRegisterRuleContextNodePool  JobRegisterRuleRegionContext = "node-pool"
)
