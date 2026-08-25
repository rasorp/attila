// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package nomad

import (
	"github.com/hashicorp/nomad/api"
	"go.uber.org/zap"

	"github.com/rasorp/attila/internal/domain"
	"github.com/rasorp/attila/internal/nomad/client"
	"github.com/rasorp/attila/internal/nomad/job"
	"github.com/rasorp/attila/internal/nomad/topology"
	"github.com/rasorp/attila/internal/server/nomad"
	"github.com/rasorp/attila/internal/store"
)

type Controller struct {
	logger   *zap.Logger
	clients  *client.Clients
	topology nomad.TopologyController
}

func NewController(logger *zap.Logger) nomad.Controller {
	clientStore := client.New(logger)
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

// JobRegistrationPlanCreate analyzes the incoming Nomad job against current
// cluster state and existing registrations to produce a registration plan.
func (c *Controller) JobRegistrationPlanCreate(apiJob *api.Job, state store.State) (*domain.JobRegisterPlan, error) {
	plan, err := job.NewPlanner(c.logger, &job.PlannerReq{
		Clients: c.clients,
		Job:     apiJob,
		State:   state,
	}).Run()
	return plan, err
}

// JobRegistrationPlanRun executes the given job registration plan by applying
// each planned deployment action against the target Nomad regions. It returns
// a result describing what was actually done.
func (c *Controller) JobRegistrationPlanRun(
	req *nomad.JobRegistrationPlanRunReq,
) (*nomad.JobRegistrationPlanRunResp, error) {

	result, err := job.NewRegister(c.logger, &job.RegisterReq{Clients: c.clients, Plan: req.Plan}).Run()
	if err != nil {
		return nil, err
	}
	return &nomad.JobRegistrationPlanRunResp{Run: result}, nil
}

func (c *Controller) GetTopologies() []*nomad.Overview {
	return c.topology.GetTopologies()
}

func (c *Controller) GetTopology(name string) *nomad.Topology {
	return c.topology.GetTopology(name)
}
