// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package nomad

import (
	"github.com/hashicorp/nomad/api"

	"github.com/rasorp/attila/internal/domain"
	"github.com/rasorp/attila/internal/store"
)

// Controller is the composite interface that combines all sub-controllers,
// providing a single entry point for clients to interact with Nomad region
// management, topology queries, and job registration planning/execution.
type Controller interface {
	ClientController
	JobRegistrationController
	TopologyController
}

// JobRegistrationPlanRunReq is the request payload passed to the
// controller's JobRegistrationPlanRun method to execute (apply) a
// previously computed job registration plan.
type JobRegistrationPlanRunReq struct {
	Plan *domain.JobRegisterPlan
}

// JobRegistrationPlanRunResp is the response object used when the controller
// has executed a run of a Nomad job registration plan.
type JobRegistrationPlanRunResp struct {
	Run *domain.JobRegisterPlanRun
}

// JobRegistrationController is the interface that defines how Attila performs
// job registration actions with the backend Nomad regions in mind.
type JobRegistrationController interface {

	// JobRegistrationPlanCreate analyzes the incoming Nomad job against current
	// cluster state and existing registrations to produce a registration plan.
	JobRegistrationPlanCreate(job *api.Job, store store.State) (*domain.JobRegisterPlan, error)

	// JobRegistrationPlanRun executes the given job registration plan by applying
	// each planned deployment action against the target Nomad regions. It returns
	// a result describing what was actually done.
	JobRegistrationPlanRun(*JobRegistrationPlanRunReq) (*JobRegistrationPlanRunResp, error)
}

type ClientController interface {

	// RegionDelete
	RegionDelete(name string)

	// RegionSet
	RegionSet(name string, client *api.Client)

	// RegionNum returns the number of regions being tracked within the
	// controller. This is a convenience method used within testing, logging,
	// and telemetry.
	RegionNum() int
}

// TopologyController is the interface that must be satisfied in order to
// implement Attila's backend topology controller.
type TopologyController interface {

	// GetTopologies returns a list of topology overviews and is used by the
	// list HTTP endpoint.
	GetTopologies() []*Overview

	// GetTopology returns the full topology object of the named Nomad region.
	// If the region is not being tracked, the implementation should return nil,
	// so the caller can check this and return a 404.
	GetTopology(name string) *Topology

	// ClientController ensures modifications to the tracked regions within
	// Attila state can be propagated to the topology controller.
	ClientController
}
