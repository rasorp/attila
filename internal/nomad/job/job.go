// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package job

import (
	"github.com/hashicorp/nomad/api"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"

	"github.com/rasorp/attila/internal/server/nomad"
	"github.com/rasorp/attila/internal/server/state"
)

type Job struct {
	clients nomad.ClientController
	logger  zerolog.Logger
	state   state.State
}

func New(clients nomad.ClientController, logger *zerolog.Logger, state state.State) nomad.JobController {
	return &Job{
		clients: clients,
		logger:  logger.With().Str("component", "job_controller").Logger(),
		state:   state,
	}
}

func (j *Job) JobRegistrationPlan(job *api.Job) (*state.JobRegisterPlan, error) {
	planRunner := newPlanner(j.clients, j.logger, j.state, job)
	return planRunner.run()
}

func (j *Job) JobRegistrationRun(planID ulid.ULID) (*state.JobRegisterPlanRun, error) {
	registerRunner, err := newRegister(
		j.clients,
		j.logger,
		planID,
		j.state,
	)
	if err != nil {
		return nil, err
	}
	return registerRunner.run()
}
