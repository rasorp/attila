// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package file

import "github.com/rasorp/attila/internal/server/state"

type JobRegister struct {
	store *Store
}

func (j *JobRegister) Method() state.JobRegisterMethodState {
	return &JobRegisterMethod{store: j.store}
}

func (j *JobRegister) Plan() state.JobRegisterPlanState { return &JobRegisterPlan{store: j.store} }

func (j *JobRegister) Rule() state.JobRegisterRuleState { return &JobRegisterRule{store: j.store} }
