// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package state

import "time"

type State interface {
	JobRegister() JobRegisterState

	Region() RegionState

	Name() string
}

type Metadata struct {
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func NewMetadata() *Metadata {
	t := time.Now()
	return &Metadata{
		CreateTime: t,
		UpdateTime: t,
	}
}
