// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"net/http"

	"github.com/hashicorp/nomad/api"
	"github.com/oklog/ulid/v2"
)

type JobRegisterMethod struct {
	Name     string                       `hcl:"name" json:"name"`
	Selector string                       `hcl:"selector" json:"selector"`
	Rules    []*JobRegisterMethodRuleLink `hcl:"rule,block" json:"rule"`
	Metadata *Metadata                    `hcl:"metadata" json:"metadata"`
}

type JobRegisterMethodStub struct {
	Name     string `hcl:"name" json:"name"`
	Selector string `hcl:"selector" json:"selector"`
}

type JobRegisterMethodRuleLink struct {
	Name string `hcl:"name" json:"name"`
}

type JobRegisterMethodCreateResp struct {
	Method *JobRegisterMethod `json:"method"`
}

type JobRegisterMethodListResp struct {
	Methods []*JobRegisterMethodStub `json:"methods"`
}

type JobRegisterMethodGetResp struct {
	Method *JobRegisterMethod `json:"method"`
}

type JobRegisterMethods struct {
	client *Client
}

func (c *Client) JobRegisterMethods() *JobRegisterMethods {
	return &JobRegisterMethods{client: c}
}

func (a *JobRegisterMethods) Create(
	ctx context.Context, method *JobRegisterMethod) (*JobRegisterMethodCreateResp, *Response, error) {

	var regionCreateResp JobRegisterMethodCreateResp

	req, err := a.client.NewRequest(http.MethodPost, "/v1alpha1/jobs/register/methods", method)
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.client.Do(ctx, req, &regionCreateResp)
	if err != nil {
		return nil, nil, err
	}

	return &regionCreateResp, resp, nil
}

func (a *JobRegisterMethods) Delete(ctx context.Context, name string) (*Response, error) {

	req, err := a.client.NewRequest(http.MethodDelete, "/v1alpha1/jobs/register/methods/"+name, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (a *JobRegisterMethods) Get(
	ctx context.Context, name string) (*JobRegisterMethodGetResp, *Response, error) {

	var methodGetResp JobRegisterMethodGetResp

	req, err := a.client.NewRequest(http.MethodGet, "/v1alpha1/jobs/register/methods/"+name, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.client.Do(ctx, req, &methodGetResp)
	if err != nil {
		return nil, resp, err
	}

	return &methodGetResp, resp, nil
}

func (a *JobRegisterMethods) List(ctx context.Context) (*JobRegisterMethodListResp, *Response, error) {

	var methodListResp JobRegisterMethodListResp

	req, err := a.client.NewRequest(http.MethodGet, "/v1alpha1/jobs/register/methods", nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.client.Do(ctx, req, &methodListResp)
	if err != nil {
		return nil, resp, err
	}

	return &methodListResp, resp, nil
}

type JobRegisterRule struct {
	Name           string                         `hcl:"name" json:"name"`
	RegionContexts []JobRegisterRuleRegionContext `hcl:"region_contexts,optional" json:"region_contexts"`
	RegionFilter   *JobRegisterRuleFilter         `hcl:"region_filter,block" json:"region_filter"`
	RegionPicker   *JobRegisterRulePicker         `hcl:"region_picker,block" json:"region_picker"`
	Metadata       *Metadata                      `hcl:"metadata" json:"metadata"`
}

type JobRegisterRuleFilter struct {
	Expression *JobRegisterRuleFilterExpression `hcl:"expression,block" json:"expression"`
}

type JobRegisterRulePicker struct {
	Expression *JobRegisterRuleFilterExpression `hcl:"expression,block" json:"expression"`
}

type JobRegisterRuleFilterExpression struct {
	Selector string `hcl:"selector" json:"selector"`
}

type JobRegisterRuleRegionContext string

const (
	JobRegisterRuleContextNamespace JobRegisterRuleRegionContext = "namespace"
	JobRegisterRuleContextNodePool  JobRegisterRuleRegionContext = "node-pool"
)

type JobRegisterRuleStub struct {
	Name           string                         `json:"name"`
	RegionContexts []JobRegisterRuleRegionContext `json:"region_contexts"`
}

type JobRegisterRuleCreateResp struct {
	Rule *JobRegisterRule `json:"rule"`
}

type JobRegisterRuleListResp struct {
	Rules []*JobRegisterRuleStub `json:"rules"`
}

type JobRegisterRuleGetResp struct {
	Rule *JobRegisterRule `json:"rule"`
}

type JobRegisterRules struct {
	client *Client
}

func (c *Client) JobRegisterRules() *JobRegisterRules {
	return &JobRegisterRules{client: c}
}

func (a *JobRegisterRules) Create(
	ctx context.Context, rule *JobRegisterRule) (*JobRegisterRuleCreateResp, *Response, error) {

	var ruleCreateResp JobRegisterRuleCreateResp

	req, err := a.client.NewRequest(http.MethodPost, "/v1alpha1/jobs/register/rules", rule)
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.client.Do(ctx, req, &ruleCreateResp)
	if err != nil {
		return nil, nil, err
	}

	return &ruleCreateResp, resp, nil
}

func (a *JobRegisterRules) Delete(ctx context.Context, name string) (*Response, error) {

	req, err := a.client.NewRequest(http.MethodDelete, "/v1alpha1/jobs/register/rules/"+name, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (a *JobRegisterRules) Get(
	ctx context.Context, name string) (*JobRegisterRuleGetResp, *Response, error) {

	var ruleGetResp JobRegisterRuleGetResp

	req, err := a.client.NewRequest(http.MethodGet, "/v1alpha1/jobs/register/rules/"+name, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.client.Do(ctx, req, &ruleGetResp)
	if err != nil {
		return nil, resp, err
	}

	return &ruleGetResp, resp, nil
}

func (a *JobRegisterRules) List(ctx context.Context) (*JobRegisterRuleListResp, *Response, error) {

	var ruleListResp JobRegisterRuleListResp

	req, err := a.client.NewRequest(http.MethodGet, "/v1alpha1/jobs/register/rules", nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.client.Do(ctx, req, &ruleListResp)
	if err != nil {
		return nil, resp, err
	}

	return &ruleListResp, resp, nil
}

type JobRegisterPlan struct {
	ID           ulid.ULID                         `json:"id"`
	JobID        string                            `json:"job_id"`
	JobNamespace string                            `json:"job_namespace"`
	Regions      map[string]*JobRegisterRegionPlan `json:"regions"`
}

type JobRegisterRegionPlan struct {
	Region string               `json:"region"`
	Plan   *api.JobPlanResponse `json:"plan"`
}

type JobRegisterPlanRun struct {
	ID           ulid.ULID                            `json:"id"`
	JobID        string                               `json:"job_id"`
	JobNamespace string                               `json:"job_namespace"`
	Regions      map[string]*JobRegisterRegionPlanRun `json:"regions"`
}

type JobRegisterRegionPlanRun struct {
	Region string                   `json:"region"`
	Run    *api.JobRegisterResponse `json:"run"`
	Error  error                    `json:"error"`
}

type JobRegisterPlanCreateReq struct {
	Job *api.Job `json:"job"`
}

type JobRegisterPlanCreateResp struct {
	Plan *JobRegisterPlan `json:"plan"`
}

type JobRegisterPlanDeleteReq struct {
	ID ulid.ULID `json:"id"`
}

type JobRegisterPlanDeleteResp struct{}

type JobRegisterPlanGetReq struct {
	ID ulid.ULID `json:"id"`
}

type JobRegisterPlanGetResp struct {
	Plan *JobRegisterPlan `json:"plan"`
}

type JobRegisterPlanListReq struct{}

type JobRegisterPlanListResp struct {
	Plans []*JobRegisterPlan `json:"plans"`
}

type JobsRegisterPlanRunReq struct {
	ID  ulid.ULID `json:"id"`
	Job *api.Job  `json:"job"`
}

type JobsRegisterPlanRunResp struct {
	Run                 *JobRegisterPlanRun `json:"run"`
	PatrialFailureError error               `json:"partial_failure_error"`
}

type JobRegisterPlans struct {
	client *Client
}

func (c *Client) JobRegisterPlans() *JobRegisterPlans {
	return &JobRegisterPlans{client: c}
}

func (j *JobRegisterPlans) Create(
	ctx context.Context, req *JobRegisterPlanCreateReq) (*JobRegisterPlanCreateResp, *Response, error) {

	var resp JobRegisterPlanCreateResp

	httpReq, err := j.client.NewRequest(http.MethodPost, "/v1alpha1/jobs/register/plans", req)
	if err != nil {
		return nil, nil, err
	}

	httpResp, err := j.client.Do(ctx, httpReq, &resp)
	if err != nil {
		return nil, nil, err
	}

	return &resp, httpResp, nil
}

func (j *JobRegisterPlans) Delete(ctx context.Context, req *JobRegisterPlanDeleteReq) (*Response, error) {

	httpReq, err := j.client.NewRequest(http.MethodDelete, "/v1alpha1/jobs/register/plans/"+req.ID.String(), nil)
	if err != nil {
		return nil, err
	}

	httpResp, err := j.client.Do(ctx, httpReq, nil)
	if err != nil {
		return nil, err
	}

	return httpResp, nil
}

func (j *JobRegisterPlans) Get(
	ctx context.Context, req *JobRegisterPlanGetReq) (*JobRegisterPlanGetResp, *Response, error) {

	var resp JobRegisterPlanGetResp

	httpReq, err := j.client.NewRequest(http.MethodGet, "/v1alpha1/jobs/register/plans/"+req.ID.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	httpResp, err := j.client.Do(ctx, httpReq, &resp)
	if err != nil {
		return nil, nil, err
	}

	return &resp, httpResp, nil
}

func (j *JobRegisterPlans) List(
	ctx context.Context, req *JobRegisterPlanListReq) (*JobRegisterPlanListResp, *Response, error) {

	var resp JobRegisterPlanListResp

	httpReq, err := j.client.NewRequest(http.MethodGet, "/v1alpha1/jobs/register/plans", req)
	if err != nil {
		return nil, nil, err
	}

	httpResp, err := j.client.Do(ctx, httpReq, &resp)
	if err != nil {
		return nil, nil, err
	}

	return &resp, httpResp, nil
}

func (j *JobRegisterPlans) Run(
	ctx context.Context, req *JobsRegisterPlanRunReq) (*JobsRegisterPlanRunResp, *Response, error) {

	var resp JobsRegisterPlanRunResp

	path := "/v1alpha1/jobs/register/plans/" + req.ID.String() + "/run"

	httpReq, err := j.client.NewRequest(http.MethodPost, path, req)
	if err != nil {
		return nil, nil, err
	}

	httpResp, err := j.client.Do(ctx, httpReq, &resp)
	if err != nil {
		return nil, nil, err
	}

	return &resp, httpResp, nil
}
