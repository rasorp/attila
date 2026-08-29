// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package store

import "github.com/rasorp/attila/internal/domain"

type JobRegisterMethodState interface {
	Create(*JobRegisterMethodCreateReq) (*JobRegisterMethodCreateResp, *ErrorResp)
	Delete(*JobRegisterMethodDeleteReq) (*JobRegisterMethodDeleteResp, *ErrorResp)
	Get(*JobRegisterMethodGetReq) (*JobRegisterMethodGetResp, *ErrorResp)
	List(*JobRegisterMethodListReq) (*JobRegisterMethodListResp, *ErrorResp)
}

type JobRegisterMethodCreateReq struct {
	Method *domain.JobRegisterMethod `json:"method"`
}

type JobRegisterMethodCreateResp struct {
	Method *domain.JobRegisterMethod `json:"method"`
}

type JobRegisterMethodDeleteReq struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type JobRegisterMethodDeleteResp struct{}

type JobRegisterMethodGetReq struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type JobRegisterMethodGetResp struct {
	Method *domain.JobRegisterMethod `json:"method"`
}

type JobRegisterMethodListReq struct {
	Namespace string `json:"namespace"`
}

type JobRegisterMethodListResp struct {
	Methods []*domain.JobRegisterMethod `json:"methods"`
}
