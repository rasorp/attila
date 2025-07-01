// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package state

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

func (a *JobRegisterRule) Validate() error { return nil }

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
	Selelctor string `json:"selector"`
}

type JobRegisterRuleRegionContext string

const (
	JobRegisterRuleContextNamespace JobRegisterRuleRegionContext = "namespace"
	JobRegisterRuleContextNodePool  JobRegisterRuleRegionContext = "node-pool"
)
