// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type SystemHealthStatus string

const (
	SystemHealthStatusHealthy   SystemHealthStatus = "healthy"
	SystemHealthStatusUnhealthy SystemHealthStatus = "unhealthy"
)

// SystemHealthResp is the response body returned by the health GET endpoint.
type SystemHealthResp struct {
	Status               SystemHealthStatus `json:"status"`
	internalResponseMeta `json:"-"`
}

// healthEndpoint serves the system/health route. It is intentionally kept
// minimal right now.
type systemEndpoint struct{}

func (s systemEndpoint) routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/health", s.getHealth)
	return r
}

func (s systemEndpoint) getHealth(w http.ResponseWriter, r *http.Request) {
	httpWriteResponse(
		w,
		&SystemHealthResp{
			Status:               SystemHealthStatusHealthy,
			internalResponseMeta: newInternalResponseMeta(http.StatusOK),
		},
	)
}
