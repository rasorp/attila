// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package job

import (
	"go.uber.org/zap"

	"github.com/hashicorp/nomad/api"
	"github.com/oklog/ulid/v2"

	"github.com/rasorp/attila/internal/domain"
	"github.com/rasorp/attila/internal/nomad/client"
	"github.com/rasorp/attila/internal/store"
)

type Register struct {
	logger *zap.Logger

	clients *client.Clients
	job     *api.Job
	state   store.State
	planID  ulid.ULID

	runResult *domain.JobRegisterPlanRun
}

type RegisterReq struct {
	Clients *client.Clients
	Job     *api.Job
	PlanID  ulid.ULID
	State   store.State
}

func NewRegister(logger *zap.Logger, req *RegisterReq) *Register {
	return &Register{
		clients: req.Clients,
		job:     req.Job,
		logger: logger.With(
			zap.String("job_id", *req.Job.ID),
			zap.String("job_namespace", *req.Job.Namespace),
			zap.String("plan_id", req.PlanID.String()),
		).Named("job_register"),
		planID:    req.PlanID,
		runResult: domain.NewJobRegisterPlanRun(*req.Job.ID, *req.Job.Namespace),
		state:     req.State,
	}
}

func (r *Register) Run() (*domain.JobRegisterPlanRun, error) {
	planResp, err := r.state.JobRegister().Plan().Get(&store.JobRegisterPlanGetReq{ID: r.planID})
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

func (r *Register) runPlannedRegion(regionPlan *domain.JobRegisterRegionPlan) error {
	apiClient, err := r.clients.Get(regionPlan.Region)
	if err != nil {
		return err
	}

	registerOpts := api.RegisterOptions{
		EnforceIndex: true,
		ModifyIndex:  regionPlan.Plan.JobModifyIndex,
	}

	r.logger.Info(
		"regional job register started",
		zap.String("region_name", regionPlan.Region),
		zap.Uint64("job_modify_index", registerOpts.ModifyIndex),
	)

	registerResp, _, err := apiClient.Jobs().RegisterOpts(r.job, &registerOpts, nil)
	r.runResult.AddRegion(regionPlan.Region, registerResp, err)

	if err != nil {
		r.logger.Error(
			"regional job register failed",
			zap.String("region_name", regionPlan.Region),
			zap.Uint64("job_modify_index", registerOpts.ModifyIndex),
			zap.Error(err),
		)
		return err
	}

	r.logger.Info(
		"regional job register successful",
		zap.String("region_name", regionPlan.Region),
		zap.Uint64("job_modify_index", registerOpts.ModifyIndex),
	)

	return nil
}
