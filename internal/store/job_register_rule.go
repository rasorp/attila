// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package store

import "github.com/rasorp/attila/internal/domain"

type JobRegisterRuleState interface {
	Create(*JobRegisterRuleCreateReq) (*JobRegisterRuleCreateResp, *ErrorResp)
	Delete(*JobRegisterRuleDeleteReq) (*JobRegisterRuleDeleteResp, *ErrorResp)
	Get(*JobRegisterRuleGetReq) (*JobRegisterRuleGetResp, *ErrorResp)
	List(*JobRegisterRuleListReq) (*JobRegisterRuleListResp, *ErrorResp)
}

type JobRegisterRuleCreateReq struct {
	Rule *domain.JobRegisterRule `json:"rule"`
}

type JobRegisterRuleCreateResp struct {
	Rule *domain.JobRegisterRule `json:"rule"`
}

type JobRegisterRuleDeleteReq struct {
	Name string `json:"name"`
}

type JobRegisterRuleDeleteResp struct{}

type JobRegisterRuleGetReq struct {
	Name string `json:"name"`
}

type JobRegisterRuleGetResp struct {
	Rule *domain.JobRegisterRule `json:"rule"`
}

type JobRegisterRuleListReq struct{}

type JobRegisterRuleListResp struct {
	Rules []*domain.JobRegisterRule `json:"rules"`
}
