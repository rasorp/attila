// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"net/http"

	"github.com/hashicorp/nomad/api"
	"github.com/oklog/ulid/v2"
)

type JobRegisterMethod struct {
	Name      string                       `hcl:"name" json:"name"`
	Namespace string                       `hcl:"namespace,optional" json:"namespace"`
	Selectors []*JobRegisterMethodSelector `hcl:"selector,block" json:"selectors"`
	Rules     []*JobRegisterMethodRuleLink `hcl:"rule,block" json:"rules"`
	Metadata  *Metadata                    `hcl:"metadata" json:"metadata"`
}

// JobRegisterMethodSelector contains all the configuration required to run the
// method selector process when selecting what rules apply to an incoming job
// plan or registration.
type JobRegisterMethodSelector struct {

	// Name provides a human friendly name to the strategy being used. This
	// allows the same type to be used multiple times and gives a useful value
	// in logs.
	Name string `hcl:",label" json:"name"`

	// Provider is the strategy provider that will be called to execute the
	// calculation.
	Provider string `hcl:"provider" json:"provider"`

	// Config is the JSON blob that will be passed to the strategy when executed
	// that controls its behaviour. The contents of this are specific to the
	// strategy implementation.
	Config map[string]any `hcl:"config,block" json:"config,omitempty"`
}

type JobRegisterMethodStub struct {
	Name      string                           `json:"name"`
	Namespace string                           `json:"namespace"`
	Selectors []*JobRegisterMethodSelectorStub `json:"selectors"`
}

type JobRegisterMethodSelectorStub struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
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
	ctx context.Context,
	method *JobRegisterMethod,
	writeOpts *WriteOpts,
) (*JobRegisterMethodCreateResp, *Response, error) {

	var regionCreateResp JobRegisterMethodCreateResp

	req, err := a.client.NewRequest(
		http.MethodPost,
		"/v1alpha1/jobs/register/methods",
		method,
		writeOpts.SetOpts(),
	)
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.client.Do(ctx, req, &regionCreateResp)
	if err != nil {
		return nil, nil, err
	}

	return &regionCreateResp, resp, nil
}

func (a *JobRegisterMethods) Delete(
	ctx context.Context,
	name string,
	writeOpts *WriteOpts,
) (*Response, error) {

	req, err := a.client.NewRequest(
		http.MethodDelete,
		"/v1alpha1/jobs/register/methods/"+name,
		nil,
		writeOpts.SetOpts(),
	)
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
	ctx context.Context,
	name string,
	queryOpts *QueryOpts,
) (*JobRegisterMethodGetResp, *Response, error) {

	var methodGetResp JobRegisterMethodGetResp

	req, err := a.client.NewRequest(
		http.MethodGet,
		"/v1alpha1/jobs/register/methods/"+name,
		nil,
		queryOpts.SetOpts(),
	)
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.client.Do(ctx, req, &methodGetResp)
	if err != nil {
		return nil, resp, err
	}

	return &methodGetResp, resp, nil
}

func (a *JobRegisterMethods) List(
	ctx context.Context,
	queryOpts *QueryOpts,
) (*JobRegisterMethodListResp, *Response, error) {

	var methodListResp JobRegisterMethodListResp

	req, err := a.client.NewRequest(
		http.MethodGet,
		"/v1alpha1/jobs/register/methods",
		nil,
		queryOpts.SetOpts(),
	)
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
	Name           string                          `hcl:"name" json:"name"`
	Namespace      string                          `hcl:"namespace,optional" json:"namespace"`
	RegionContexts []*JobRegisterRuleRegionContext `hcl:"region_context,block" json:"region_contexts"`
	RegionPickers  []*JobRegisterRegionPicker      `hcl:"region_picker,block" json:"region_pickers"`
	Metadata       *Metadata                       `hcl:"metadata" json:"metadata"`
}

// JobRegisterRuleRegionContext is a way to make additional resource information
// about the region being available to the picker.
type JobRegisterRuleRegionContext struct {

	// Kind is the context that will be made available to the rule picker. It
	// currently supports namespace and node-pool.
	Kind string `hcl:"kind" json:"kind"`
}

const (
	// JobRegisterRuleContextKindNamespace is the context kind that supplies
	// information from Nomad's "v1/namespaces" endpoint.
	JobRegisterRuleContextKindNamespace = "namespace"

	// JobRegisterRuleContextKindNodepool is the context kind that supplies
	// information from Nomad's "v1/node/pools" endpoint.
	JobRegisterRuleContextKindNodepool = "node-pool"
)

// JobRegisterRegionPicker contains all the configuration required to run the
// region picker process when selecting what regions to register a job into.
type JobRegisterRegionPicker struct {

	// Name provides a human friendly name to the strategy being used. This
	// allows the same type to be used multiple times and gives a useful value
	// in logs.
	Name string `hcl:",label" json:"name"`

	// Provider is the strategy provider that will be called to execute the
	// calculation.
	Provider string `hcl:"provider" json:"provider"`

	// Config is the JSON blob that will be passed to the strategy when executed
	// that controls its behaviour. The contents of this are specific to the
	// strategy implementation.
	Config map[string]any `hcl:"config,block" json:"config,omitempty"`
}

type JobRegisterRuleStub struct {
	Name           string                         `json:"name"`
	Namespace      string                         `json:"namespace"`
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
	ctx context.Context,
	rule *JobRegisterRule,
	writeOpts *WriteOpts,
) (*JobRegisterRuleCreateResp, *Response, error) {

	var ruleCreateResp JobRegisterRuleCreateResp

	req, err := a.client.NewRequest(
		http.MethodPost,
		"/v1alpha1/jobs/register/rules",
		rule,
		writeOpts.SetOpts(),
	)
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.client.Do(ctx, req, &ruleCreateResp)
	if err != nil {
		return nil, nil, err
	}

	return &ruleCreateResp, resp, nil
}

func (a *JobRegisterRules) Delete(
	ctx context.Context,
	name string,
	writeOpts *WriteOpts,
) (*Response, error) {

	req, err := a.client.NewRequest(
		http.MethodDelete,
		"/v1alpha1/jobs/register/rules/"+name,
		nil,
		writeOpts.SetOpts(),
	)
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
	ctx context.Context,
	name string,
	queryOpts *QueryOpts,
) (*JobRegisterRuleGetResp, *Response, error) {

	var ruleGetResp JobRegisterRuleGetResp

	req, err := a.client.NewRequest(
		http.MethodGet,
		"/v1alpha1/jobs/register/rules/"+name,
		nil,
		queryOpts.SetOpts(),
	)
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.client.Do(ctx, req, &ruleGetResp)
	if err != nil {
		return nil, resp, err
	}

	return &ruleGetResp, resp, nil
}

func (a *JobRegisterRules) List(
	ctx context.Context,
	queryOpts *QueryOpts,
) (*JobRegisterRuleListResp, *Response, error) {

	var ruleListResp JobRegisterRuleListResp

	req, err := a.client.NewRequest(
		http.MethodGet,
		"/v1alpha1/jobs/register/rules",
		nil,
		queryOpts.SetOpts(),
	)
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
	ID        ulid.ULID                         `json:"id"`
	Namespace string                            `json:"namespace"`
	Job       *api.Job                          `json:"job"`
	Regions   map[string]*JobRegisterRegionPlan `json:"regions"`
}

type JobRegisterRegionPlan struct {
	Region string               `json:"region"`
	Plan   *api.JobPlanResponse `json:"plan"`
}

type JobRegisterRegionPlanRun struct {
	Region       string                   `json:"region"`
	RegisterResp *api.JobRegisterResponse `json:"register_response"`
	Error        error                    `json:"error"`
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
	ID ulid.ULID `json:"id"`
}

type JobsRegisterPlanRunResp struct {
	Run                 *JobRegisterRun `json:"run"`
	PatrialFailureError error           `json:"partial_failure_error"`
}

type JobRegisterPlans struct {
	client *Client
}

func (c *Client) JobRegisterPlans() *JobRegisterPlans {
	return &JobRegisterPlans{client: c}
}

func (j *JobRegisterPlans) Create(
	ctx context.Context,
	req *JobRegisterPlanCreateReq,
	writeOpts *WriteOpts,
) (*JobRegisterPlanCreateResp, *Response, error) {

	var resp JobRegisterPlanCreateResp

	httpReq, err := j.client.NewRequest(
		http.MethodPost,
		"/v1alpha1/jobs/register/plans",
		req,
		writeOpts.SetOpts(),
	)
	if err != nil {
		return nil, nil, err
	}

	httpResp, err := j.client.Do(ctx, httpReq, &resp)
	if err != nil {
		return nil, nil, err
	}

	return &resp, httpResp, nil
}

func (j *JobRegisterPlans) Delete(
	ctx context.Context,
	req *JobRegisterPlanDeleteReq,
	writeOpts *WriteOpts,
) (*Response, error) {

	httpReq, err := j.client.NewRequest(
		http.MethodDelete,
		"/v1alpha1/jobs/register/plans/"+req.ID.String(),
		nil,
		writeOpts.SetOpts(),
	)
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
	ctx context.Context,
	req *JobRegisterPlanGetReq,
	queryOpts *QueryOpts,
) (*JobRegisterPlanGetResp, *Response, error) {

	var resp JobRegisterPlanGetResp

	httpReq, err := j.client.NewRequest(
		http.MethodGet,
		"/v1alpha1/jobs/register/plans/"+req.ID.String(),
		nil,
		queryOpts.SetOpts(),
	)
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
	ctx context.Context,
	req *JobRegisterPlanListReq,
	queryOpts *QueryOpts,
) (*JobRegisterPlanListResp, *Response, error) {

	var resp JobRegisterPlanListResp

	httpReq, err := j.client.NewRequest(
		http.MethodGet,
		"/v1alpha1/jobs/register/plans",
		req,
		queryOpts.SetOpts(),
	)
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
	ctx context.Context,
	req *JobsRegisterPlanRunReq,
	queryOpts *QueryOpts,
) (*JobsRegisterPlanRunResp, *Response, error) {

	var resp JobsRegisterPlanRunResp

	httpReq, err := j.client.NewRequest(
		http.MethodPost,
		"/v1alpha1/jobs/register/plans/"+req.ID.String()+"/run",
		req,
		queryOpts.SetOpts(),
	)
	if err != nil {
		return nil, nil, err
	}

	httpResp, err := j.client.Do(ctx, httpReq, &resp)
	if err != nil {
		return nil, nil, err
	}

	return &resp, httpResp, nil
}
