// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package domain

import (
	"github.com/hashicorp/nomad/api"
	"github.com/oklog/ulid/v2"
)

// JobRegisterPlanRunState provides storage access for plan execution results.
type JobRegisterRunState interface {
	Create(*JobRegisterRunCreateReq) (*JobRegisterRunCreateResp, StateError)
	Delete(*JobRegisterRunDeleteReq) (*JobRegisterRunDeleteResp, StateError)
	Get(*JobRegisterRunGetReq) (*JobRegisterRunGetResp, StateError)
	List(*JobRegisterRunListReq) (*JobRegisterRunListResp, StateError)
}

// JobRegisterPlanRunCreateReq wraps a plan run for creation.
type JobRegisterRunCreateReq struct {
	Run *JobRegisterRun
}

// JobRegisterPlanRunCreateResp returns the created plan run record.
type JobRegisterRunCreateResp struct {
	Run *JobRegisterRun `json:"run"`
}

// JobRegisterPlanRunCreateReq wraps a plan run for creation.
type JobRegisterRunDeleteReq struct {
	ID        ulid.ULID
	Namespace string
}

// JobRegisterPlanRunCreateResp returns the created plan run record.
type JobRegisterRunDeleteResp struct{}

// JobRegisterPlanRunGetReq identifies a plan run by ID and namespace.
type JobRegisterRunGetReq struct {
	ID        ulid.ULID
	Namespace string
}

// JobRegisterPlanRunGetResp returns the plan run record.
type JobRegisterRunGetResp struct {
	Run *JobRegisterRun
}

type JobRegisterRunListReq struct {
	Namespace string
}

type JobRegisterRunListResp struct {
	Runs []*JobRegisterRun
}

type JobRegisterRun struct {
	ID           ulid.ULID                            `json:"id"`
	Namespace    string                               `json:"namespace"`
	JobID        string                               `json:"job_id"`
	JobNamespace string                               `json:"job_namespace"`
	Regions      map[string]*JobRegisterRegionPlanRun `json:"regions"`
}

type JobRegisterRegionPlanRun struct {
	Region       string                   `json:"region"`
	RegisterResp *api.JobRegisterResponse `json:"register_response"`
	Error        error                    `json:"error"`
}

func NewJobRegisterRun(plan *JobRegisterPlan) *JobRegisterRun {
	return &JobRegisterRun{
		ID:           ulid.Make(),
		Namespace:    plan.Namespace,
		JobID:        *plan.Job.ID,
		JobNamespace: *plan.Job.Namespace,
		Regions:      make(map[string]*JobRegisterRegionPlanRun),
	}
}

func (j *JobRegisterRun) AddRegion(regionName string, regResp *api.JobRegisterResponse, err error) {
	runResp := JobRegisterRegionPlanRun{Region: regionName}

	if err != nil {
		runResp.Error = err
	} else {
		runResp.RegisterResp = regResp
	}

	j.Regions[regionName] = &runResp
}

func (j *JobRegisterRun) Stub() *JobRegisterRunStub {
	return &JobRegisterRunStub{
		ID:           j.ID,
		Namespace:    j.Namespace,
		JobID:        j.JobID,
		JobNamespace: j.JobNamespace,
	}
}

type JobRegisterRunStub struct {
	ID           ulid.ULID `json:"id"`
	Namespace    string    `json:"namespace"`
	JobID        string    `json:"job_id"`
	JobNamespace string    `json:"job_namespace"`
}
