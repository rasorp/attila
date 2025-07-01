// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package job

import (
	"github.com/hashicorp/nomad/api"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"

	"github.com/rasorp/attila/internal/nomad/client"
	"github.com/rasorp/attila/internal/server/state"
)

type Register struct {
	logger zerolog.Logger

	clients *client.Clients
	job     *api.Job
	state   state.State
	planID  ulid.ULID

	runResult *state.JobRegisterPlanRun
}

type RegisterReq struct {
	Clients *client.Clients
	Job     *api.Job
	PlanID  ulid.ULID
	State   state.State
}

func NewRegister(logger zerolog.Logger, req *RegisterReq) *Register {
	return &Register{
		clients: req.Clients,
		job:     req.Job,
		logger: logger.With().
			Str("job_id", *req.Job.ID).
			Str("job_namespace", *req.Job.Namespace).
			Str("plan_id", req.PlanID.String()).
			Str("component", "job_register_runner").Logger(),
		planID:    req.PlanID,
		runResult: state.NewJobRegisterPlanRun(*req.Job.ID, *req.Job.Namespace),
		state:     req.State,
	}
}

func (r *Register) Run() (*state.JobRegisterPlanRun, error) {

	planResp, err := r.state.JobRegister().Plan().Get(&state.JobRegisterPlanGetReq{ID: r.planID})
	if err != nil {
		return nil, err
	}

	for _, plannedRegion := range planResp.Plan.Regions {
		if err := r.runPlannedRegion(plannedRegion); err != nil {
			return nil, err
		}
	}

	return r.runResult, nil
}

func (r *Register) runPlannedRegion(regionPlan *state.JobRegisterRegionPlan) error {

	apiClient, err := r.clients.Get(regionPlan.Region)
	if err != nil {
		return err
	}

	//
	registerOpts := api.RegisterOptions{
		EnforceIndex: true,
		ModifyIndex:  regionPlan.Plan.JobModifyIndex,
	}

	r.logger.Info().
		Str("region_name", regionPlan.Region).
		Uint64("job_modify_index", registerOpts.ModifyIndex).
		Msg("regional job register started")

	//
	registerResp, _, err := apiClient.Jobs().RegisterOpts(r.job, &registerOpts, nil)
	r.runResult.AddRegion(regionPlan.Region, registerResp, err)

	if err != nil {
		r.logger.Error().
			Err(err).
			Str("region_name", regionPlan.Region).
			Uint64("job_modify_index", registerOpts.ModifyIndex).
			Msg("regional job register failed")
		return err
	}

	r.logger.Info().
		Str("region_name", regionPlan.Region).
		Uint64("job_modify_index", registerOpts.ModifyIndex).
		Msg("regional job register successful")

	return nil
}
