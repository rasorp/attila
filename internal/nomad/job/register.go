// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package job

import (
	"github.com/hashicorp/nomad/api"
	"go.uber.org/zap"

	"github.com/rasorp/attila/internal/domain"
	"github.com/rasorp/attila/internal/nomad/client"
)

type Register struct {
	logger    *zap.Logger
	clients   *client.Clients
	plan      *domain.JobRegisterPlan
	runResult *domain.JobRegisterPlanRun
}

type RegisterReq struct {
	Clients *client.Clients
	Plan    *domain.JobRegisterPlan
}

func NewRegister(logger *zap.Logger, req *RegisterReq) *Register {
	return &Register{
		clients: req.Clients,
		logger: logger.With(
			zap.String("job_id", *req.Plan.Job.ID),
			zap.String("job_namespace", *req.Plan.Job.Namespace),
			zap.String("plan_id", req.Plan.ID.String()),
		).Named("job_register"),
		plan:      req.Plan,
		runResult: domain.NewJobRegisterPlanRun(req.Plan.Job),
	}
}

func (r *Register) Run() (*domain.JobRegisterPlanRun, error) {
	for _, plannedRegion := range r.plan.Regions {
		if err := r.runPlannedRegion(plannedRegion, r.plan.Job); err != nil {
			return nil, err
		}
	}
	return r.runResult, nil
}

func (r *Register) runPlannedRegion(regionPlan *domain.JobRegisterRegionPlan, apiJob *api.Job) error {
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

	registerResp, _, err := apiClient.Jobs().RegisterOpts(apiJob, &registerOpts, nil)
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
