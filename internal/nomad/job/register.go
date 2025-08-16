// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package job

import (
	"github.com/hashicorp/nomad/api"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"

	"github.com/rasorp/attila/internal/server/nomad"
	"github.com/rasorp/attila/internal/server/state"
)

type register struct {
	clients nomad.ClientController
	logger  zerolog.Logger
	state   state.State

	plan      *state.JobRegisterPlan
	runResult *state.JobRegisterPlanRun
}

func newRegister(
	clients nomad.ClientController,
	logger zerolog.Logger,
	planID ulid.ULID,
	store state.State,
) (*register, error) {

	planResp, err := store.JobRegister().Plan().Get(&state.JobRegisterPlanGetReq{ID: planID})
	if err != nil {
		return nil, err
	}

	return &register{
		clients: clients,
		logger: logger.With().
			Str("job_id", *planResp.Plan.Job.ID).
			Str("job_namespace", *planResp.Plan.Job.Namespace).
			Str("plan_id", planID.String()).
			Str("component", "job_register_runner").Logger(),
		plan:      planResp.Plan,
		runResult: state.NewJobRegisterPlanRun(planResp.Plan.Job),
		state:     store,
	}, nil
}

func (r *register) run() (*state.JobRegisterPlanRun, error) {

	for _, plannedRegion := range r.plan.Regions {
		if err := r.runPlannedRegion(plannedRegion); err != nil {
			return nil, err
		}
	}

	return r.runResult, nil
}

func (r *register) runPlannedRegion(regionPlan *state.JobRegisterRegionPlan) error {

	apiClient, err := r.clients.RegionGet(regionPlan.Region)
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
	registerResp, _, err := apiClient.Jobs().RegisterOpts(r.plan.Job, &registerOpts, nil)
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
