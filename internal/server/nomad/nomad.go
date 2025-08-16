// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package nomad

import (
	"github.com/hashicorp/nomad/api"
	"github.com/oklog/ulid/v2"

	"github.com/rasorp/attila/internal/server/state"
)

// Controller
type Controller interface {
	ClientController
	JobController
	TopologyController
}

// JobController is the interface that must be satisfied in order to implement
// Attila's backend that performs Nomad job planning and registration.
type JobController interface {

	// JobRegistrationPlan generates a job registration plan by processing the
	// submitted job against Attila's stored registration methods and rules.
	JobRegistrationPlan(job *api.Job) (*state.JobRegisterPlan, error)

	// JobRegistrationRun executes the job registration plan as specified by the
	// passed plan ID.
	JobRegistrationRun(planID ulid.ULID) (*state.JobRegisterPlanRun, error)
}

type ClientController interface {

	// RegionDelete
	RegionDelete(name string)

	RegionGet(name string) (*api.Client, error)

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
