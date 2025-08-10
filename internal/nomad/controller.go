// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package nomad

import (
	"github.com/hashicorp/nomad/api"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"

	"github.com/rasorp/attila/internal/nomad/client"
	"github.com/rasorp/attila/internal/nomad/job"
	"github.com/rasorp/attila/internal/nomad/topology"
	"github.com/rasorp/attila/internal/server/nomad"
	"github.com/rasorp/attila/internal/server/state"
)

type Controller struct {
	logger   *zerolog.Logger
	clients  *client.Clients
	topology nomad.TopologyController
}

func NewController(logger *zerolog.Logger) nomad.Controller {

	clientStore := client.New(*logger)
	topologyController := topology.New(logger, clientStore)

	return &Controller{
		logger:   logger,
		clients:  clientStore,
		topology: topologyController,
	}
}

func (c *Controller) RegionDelete(name string) {
	c.topology.RegionDelete(name)
	c.clients.Delete(name)
}

func (c *Controller) RegionSet(name string, client *api.Client) {
	c.clients.Set(name, client)
	c.topology.RegionSet(name, nil)
}

func (c *Controller) RegionNum() int { return c.clients.Num() }

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

func (c *Controller) GetTopologies() []*nomad.Overview {
	return c.topology.GetTopologies()
}

func (c *Controller) GetTopology(name string) *nomad.Topology {
	return c.topology.GetTopology(name)
}
