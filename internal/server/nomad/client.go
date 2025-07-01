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

	// RegionClientDelete
	RegionClientDelete(name string)

	// RegionClientSet
	RegionClientSet(name string, client *api.Client)

	// JobRegistrationPlanCreate
	JobRegistrationPlanCreate(job *api.Job, state state.State) (*state.JobRegisterPlan, error)

	// JobRegistrationRun
	JobRegistrationRun(planID ulid.ULID, job *api.Job, state state.State) (*state.JobRegisterPlanRun, error)
}
