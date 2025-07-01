// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package nomad

import (
	"github.com/hashicorp/nomad/api"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"

	"github.com/rasorp/attila/internal/nomad/client"
	"github.com/rasorp/attila/internal/nomad/job"
	"github.com/rasorp/attila/internal/server/nomad"
	"github.com/rasorp/attila/internal/server/state"
)

type Controller struct {
	logger  *zerolog.Logger
	clients *client.Clients
}

func NewController(logger *zerolog.Logger) nomad.Controller {
	return &Controller{
		logger:  logger,
		clients: client.New(*logger),
	}
}

func (c *Controller) RegionClientDelete(name string) { c.clients.Delete(name) }

func (c *Controller) RegionClientSet(name string, client *api.Client) { c.clients.Set(name, client) }

func (c *Controller) JobRegistrationPlanCreate(apiJob *api.Job, state state.State) (*state.JobRegisterPlan, error) {
	return job.NewPlanner(*c.logger, &job.PlannerReq{
		Clients: c.clients,
		Job:     apiJob,
		State:   state,
	}).Run()
}

func (c *Controller) JobRegistrationRun(planID ulid.ULID, apiJob *api.Job, state state.State) (*state.JobRegisterPlanRun, error) {
	return job.NewRegister(*c.logger, &job.RegisterReq{
		Clients: c.clients,
		Job:     apiJob,
		PlanID:  planID,
		State:   state,
	}).Run()
}
