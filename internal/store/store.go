// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package store

type State interface {
	JobRegister() JobRegisterState
	Region() RegionState
	Name() string
}
