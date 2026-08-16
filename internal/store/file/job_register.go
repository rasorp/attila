// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package file

import "github.com/rasorp/attila/internal/store"

type JobRegister struct {
	store *Store
}

func (j *JobRegister) Method() store.JobRegisterMethodState {
	return &JobRegisterMethod{store: j.store}
}
func (j *JobRegister) Plan() store.JobRegisterPlanState { return &JobRegisterPlan{store: j.store} }
func (j *JobRegister) Rule() store.JobRegisterRuleState { return &JobRegisterRule{store: j.store} }
