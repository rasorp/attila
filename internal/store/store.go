// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package store

import "github.com/rasorp/attila/internal/server/state"

type State interface {
	JobRegister() JobRegisterState
	Namespace() state.Namespace
	Region() RegionState
	Name() string
}
