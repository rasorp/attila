// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package domain

import (
	"github.com/hashicorp/nomad/api"
	"github.com/oklog/ulid/v2"
)

// JobRegisterPlanState provides access to job registration plan resources.
type JobRegisterPlanState interface {
	Create(*JobRegisterPlanCreateReq) (*JobRegisterPlanCreateResp, StateError)
	Delete(*JobRegisterPlanDeleteReq) (*JobRegisterPlanDeleteResp, StateError)
	Get(*JobRegisterPlanGetReq) (*JobRegisterPlanGetResp, StateError)
	List(*JobRegisterPlanListReq) (*JobRegisterPlanListResp, StateError)
}

type JobRegisterPlanCreateReq struct {
	Plan *JobRegisterPlan
}

type JobRegisterPlanCreateResp struct {
	Plan *JobRegisterPlan `json:"plan"`
}

type JobRegisterPlanDeleteReq struct {
	ID        ulid.ULID `json:"id"`
	Namespace string    `json:"namespace"`
}

type JobRegisterPlanDeleteResp struct{}

type JobRegisterPlanGetReq struct {
	ID        ulid.ULID `json:"id"`
	Namespace string    `json:"namespace"`
}

type JobRegisterPlanGetResp struct {
	Plan *JobRegisterPlan `json:"plan"`
}

type JobRegisterPlanListReq struct {
	Namespace string `json:"namespace"`
}

type JobRegisterPlanListResp struct {
	Plans []*JobRegisterPlan `json:"plans"`
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

func NewJobRegisterPlan(namespace string, job *api.Job) *JobRegisterPlan {
	return &JobRegisterPlan{
		ID:        ulid.Make(),
		Namespace: namespace,
		Job:       job,
		Regions:   make(map[string]*JobRegisterRegionPlan),
	}
}

func (j *JobRegisterPlan) AddRegion(region *Region, nomadPlan *api.JobPlanResponse) {
	j.Regions[region.Name] = &JobRegisterRegionPlan{
		Region: region.Name,
		Plan:   nomadPlan,
	}
}
