// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"github.com/oklog/ulid/v2"

	"github.com/rasorp/attila/internal/domain"
)

type JobRegisterPlanState interface {
	Create(*JobRegisterPlanCreateReq) (*JobRegisterPlanCreateResp, *ErrorResp)
	Delete(*JobRegisterPlanDeleteReq) (*JobRegisterPlanDeleteResp, *ErrorResp)
	Get(*JobRegisterPlanGetReq) (*JobRegisterPlanGetResp, *ErrorResp)
	List(*JobRegisterPlanListReq) (*JobRegisterPlanListResp, *ErrorResp)
}

type JobRegisterPlanCreateReq struct {
	Plan *domain.JobRegisterPlan
}

type JobRegisterPlanCreateResp struct {
	Plan *domain.JobRegisterPlan `json:"plan"`
}

type JobRegisterPlanDeleteReq struct {
	ID ulid.ULID `json:"id"`
}

type JobRegisterPlanDeleteResp struct{}

type JobRegisterPlanGetReq struct {
	ID ulid.ULID `json:"id"`
}

type JobRegisterPlanGetResp struct {
	Plan *domain.JobRegisterPlan `json:"plan"`
}

type JobRegisterPlanListReq struct{}

type JobRegisterPlanListResp struct {
	Plans []*domain.JobRegisterPlan `json:"plans"`
}
