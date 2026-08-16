// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package store

import "github.com/rasorp/attila/internal/domain"

type RegionState interface {
	Create(*RegionCreateReq) (*RegionCreateResp, *ErrorResp)
	Delete(*RegionDeleteReq) (*RegionDeleteResp, *ErrorResp)
	Get(*RegionGetReq) (*RegionGetResp, *ErrorResp)
	List(*RegionListReq) (*RegionListResp, *ErrorResp)
}

type RegionCreateReq struct {
	Region *domain.Region
}

type RegionCreateResp struct {
	Region *domain.Region `json:"region"`
}

type RegionDeleteReq struct {
	RegionName string
}

type RegionDeleteResp struct{}

type RegionGetReq struct {
	RegionName string
}

type RegionGetResp struct {
	Region *domain.Region `json:"region"`
}

type RegionListReq struct{}

type RegionListResp struct {
	Regions []*domain.Region `json:"regions"`
}
