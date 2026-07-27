// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package api

import "time"

type Metadata struct {
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}
