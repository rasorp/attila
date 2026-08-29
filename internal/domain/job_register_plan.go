// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package domain

import (
	"github.com/hashicorp/nomad/api"
	"github.com/oklog/ulid/v2"
)

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

type JobRegisterPlanRun struct {
	ID           ulid.ULID                            `json:"id"`
	JobID        string                               `json:"job_id"`
	JobNamespace string                               `json:"job_namespace"`
	Regions      map[string]*JobRegisterRegionPlanRun `json:"regions"`
}

type JobRegisterRegionPlanRun struct {
	Region       string                   `json:"region"`
	RegisterResp *api.JobRegisterResponse `json:"register_response"`
	Error        error                    `json:"error"`
}

func NewJobRegisterPlanRun(job *api.Job) *JobRegisterPlanRun {
	return &JobRegisterPlanRun{
		ID:           ulid.Make(),
		JobID:        *job.ID,
		JobNamespace: *job.Namespace,
		Regions:      make(map[string]*JobRegisterRegionPlanRun),
	}
}

func (j *JobRegisterPlanRun) AddRegion(regionName string, regResp *api.JobRegisterResponse, err error) {
	runResp := JobRegisterRegionPlanRun{Region: regionName}

	if err != nil {
		runResp.Error = err
	} else {
		runResp.RegisterResp = regResp
	}

	j.Regions[regionName] = &runResp
}
