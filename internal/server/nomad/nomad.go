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

	// JobRegistrationPlanCreate
	JobRegistrationPlanCreate(job *api.Job, state state.State) (*state.JobRegisterPlan, error)

	// JobRegistrationRun
	JobRegistrationRun(planID ulid.ULID, job *api.Job, state state.State) (*state.JobRegisterPlanRun, error)

	TopologyController
}

type ClientController interface {

	// RegionDelete
	RegionDelete(name string)

	// RegionSet
	RegionSet(name string, client *api.Client)
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
