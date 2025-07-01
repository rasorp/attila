// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package state

type JobRegisterState interface {
	Plan() JobRegisterPlanState
	Method() JobRegisterMethodState
	Rule() JobRegisterRuleState
}
