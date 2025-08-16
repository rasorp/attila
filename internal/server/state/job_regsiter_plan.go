// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"github.com/hashicorp/nomad/api"
	"github.com/oklog/ulid/v2"
)

type JobRegisterPlanState interface {

	// Create
	Create(*JobRegisterPlanCreateReq) (*JobRegisterPlanCreateResp, *ErrorResp)

	// Delete
	Delete(*JobRegisterPlanDeleteReq) (*JobRegisterPlanDeleteResp, *ErrorResp)

	// Get
	Get(*JobRegisterPlanGetReq) (*JobRegisterPlanGetResp, *ErrorResp)

	// List
	List(*JobRegisterPlanListReq) (*JobRegisterPlanListResp, *ErrorResp)
}

type JobRegisterPlanCreateReq struct {
	Plan *JobRegisterPlan
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

type JobRegisterPlan struct {
	ID      ulid.ULID                         `json:"id"`
	Job     *api.Job                          `json:"job"`
	Regions map[string]*JobRegisterRegionPlan `json:"regions"`
}

type JobRegisterRegionPlan struct {
	Region string               `json:"region"`
	Plan   *api.JobPlanResponse `json:"plan"`
}

func NewJobRegisterPlan(job *api.Job) *JobRegisterPlan {
	return &JobRegisterPlan{
		ID:      ulid.Make(),
		Job:     job,
		Regions: make(map[string]*JobRegisterRegionPlan),
	}
}

func (j *JobRegisterPlan) AddRegion(region *Region, nomadPlan *api.JobPlanResponse) {
	j.Regions[region.Name] = &JobRegisterRegionPlan{
		Region: region.Name,
		Plan:   nomadPlan,
	}
}

type JobRegisterPlanRun struct {
	ID      ulid.ULID                            `json:"id"`
	Job     *api.Job                             `json:"job"`
	Regions map[string]*JobRegisterRegionPlanRun `json:"regions"`
}

type JobRegisterRegionPlanRun struct {
	Region string                   `json:"region"`
	Run    *api.JobRegisterResponse `json:"run"`
	Error  error                    `json:"error"`
}

func NewJobRegisterPlanRun(job *api.Job) *JobRegisterPlanRun {
	return &JobRegisterPlanRun{
		ID:      ulid.Make(),
		Job:     job,
		Regions: make(map[string]*JobRegisterRegionPlanRun),
	}
}

func (j *JobRegisterPlanRun) AddRegion(regionName string, regResp *api.JobRegisterResponse, err error) {

	runResp := JobRegisterRegionPlanRun{Region: regionName}

	if err != nil {
		runResp.Error = err
	} else {
		runResp.Run = regResp
	}

	j.Regions[regionName] = &runResp
}
