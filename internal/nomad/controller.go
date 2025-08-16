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
	clients  nomad.ClientController
	job      nomad.JobController
	topology nomad.TopologyController
}

func NewController(logger *zerolog.Logger, store state.State) nomad.Controller {

	clientStore := client.New(*logger)
	topologyController := topology.New(logger, clientStore)
	jobController := job.New(clientStore, logger, store)

	return &Controller{
		logger:   logger,
		clients:  clientStore,
		job:      jobController,
		topology: topologyController,
	}
}

func (c *Controller) RegionDelete(name string) {
	c.topology.RegionDelete(name)
	c.clients.RegionDelete(name)
}

func (c *Controller) RegionGet(name string) (*api.Client, error) {
	return c.clients.RegionGet(name)
}

func (c *Controller) RegionSet(name string, client *api.Client) {
	c.clients.RegionSet(name, client)
	c.topology.RegionSet(name, nil)
}

func (c *Controller) RegionNum() int { return c.clients.RegionNum() }

func (c *Controller) JobRegistrationPlan(job *api.Job) (*state.JobRegisterPlan, error) {
	return c.job.JobRegistrationPlan(job)
}

func (c *Controller) JobRegistrationRun(planID ulid.ULID) (*state.JobRegisterPlanRun, error) {
	return c.job.JobRegistrationRun(planID)
}

func (c *Controller) GetTopologies() []*nomad.Overview {
	return c.topology.GetTopologies()
}

func (c *Controller) GetTopology(name string) *nomad.Topology {
	return c.topology.GetTopology(name)
}
