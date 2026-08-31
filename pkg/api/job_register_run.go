// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"net/http"

	"github.com/hashicorp/nomad/api"
	"github.com/oklog/ulid/v2"
)

type JobRegisterRun struct {
	ID           ulid.ULID                            `json:"id"`
	Namespace    string                               `json:"namespace"`
	JobID        string                               `json:"job_id"`
	JobNamespace string                               `json:"job_namespace"`
	Regions      map[string]*JobRegisterRegionPlanRun `json:"regions"`
}

type JobRegisterRunStub struct {
	ID           ulid.ULID `json:"id"`
	Namespace    string    `json:"namespace"`
	JobID        string    `json:"job_id"`
	JobNamespace string    `json:"job_namespace"`
}

type JobRegisterRuns struct {
	client *Client
}

func (c *Client) JobRegisterRuns() *JobRegisterRuns {
	return &JobRegisterRuns{client: c}
}

type JobsRegisterRunsCreateReq struct {
	PlanID ulid.ULID `json:"plan_id,omitempty"`
	Job    *api.Job  `json:"job,omitempty"`
}

type JobsRegisterRunsCreateResp struct {
	Run *JobRegisterRun `json:"run"`
}

type JobsRegisterRunsDeleteReq struct {
	ID ulid.ULID
}

// JobsRegisterRunsListResp is the response from a List call.
type JobsRegisterRunsListResp struct {
	Runs []*JobRegisterRunStub `json:"runs"`
}

// JobRegisterRunsGetResp is the response from a Get call.
type JobRegisterRunsGetResp struct {
	Run *JobRegisterRun `json:"run"`
}

func (j *JobRegisterRuns) Create(
	ctx context.Context,
	req *JobsRegisterRunsCreateReq,
	writeOpts *WriteOpts,
) (*JobsRegisterRunsCreateResp, *Response, error) {

	var resp JobsRegisterRunsCreateResp

	httpReq, err := j.client.NewRequest(
		http.MethodPost,
		"/v1alpha1/jobs/register/runs",
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

func (j *JobRegisterRuns) Delete(
	ctx context.Context,
	req *JobsRegisterRunsDeleteReq,
	writeOpts *WriteOpts,
) (*Response, error) {

	httpReq, err := j.client.NewRequest(
		http.MethodDelete,
		"/v1alpha1/jobs/register/runs/"+req.ID.String(),
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

func (j *JobRegisterRuns) Get(
	ctx context.Context,
	id ulid.ULID,
	queryOpts *QueryOpts,
) (*JobRegisterRunsGetResp, *Response, error) {

	var resp JobRegisterRunsGetResp

	httpReq, err := j.client.NewRequest(
		http.MethodGet,
		"/v1alpha1/jobs/register/runs/"+id.String(),
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

func (j *JobRegisterRuns) List(
	ctx context.Context,
	queryOpts *QueryOpts,
) (*JobsRegisterRunsListResp, *Response, error) {

	var resp JobsRegisterRunsListResp

	httpReq, err := j.client.NewRequest(
		http.MethodGet,
		"/v1alpha1/jobs/register/runs",
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
