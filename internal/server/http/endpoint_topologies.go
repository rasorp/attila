// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/rasorp/attila/internal/server/nomad"
)

type TopologyListReq struct{}

type TopologyListResp struct {
	Topologies           []*nomad.Overview
	internalResponseMeta `json:"-"`
}

type TopologyGetReq struct{}

type TopologyGetResp struct {
	Topology             *nomad.Topology
	internalResponseMeta `json:"-"`
}

type topologiesEndpoint struct {
	nomadController nomad.Controller
}

func (t topologiesEndpoint) routes() chi.Router {
	r := chi.NewRouter()

	// Add the root endpoints which do not include a region name within the URI.
	r.Route("/", func(r chi.Router) {
		r.Get("/", t.list)
	})

	// Add the endpoints which are specific to a named region.
	r.Route("/{regionName}", func(r chi.Router) {
		r.Use(t.context)
		r.Get("/", t.get)
	})

	return r
}

func (t topologiesEndpoint) get(w http.ResponseWriter, r *http.Request) {
	regionName := r.Context().Value("region-name").(string)

	if topology := t.nomadController.GetTopology(regionName); topology == nil {
		respErr := NewResponseError(errors.New("region not found"), http.StatusNotFound)
		httpWriteResponseError(w, respErr)
	} else {
		resp := TopologyGetResp{
			Topology:             topology,
			internalResponseMeta: newInternalResponseMeta(http.StatusOK),
		}
		httpWriteResponse(w, &resp)
	}
}

func (t topologiesEndpoint) list(w http.ResponseWriter, r *http.Request) {
	resp := TopologyListResp{
		Topologies:           t.nomadController.GetTopologies(),
		internalResponseMeta: newInternalResponseMeta(http.StatusOK),
	}
	httpWriteResponse(w, &resp)
}

func (t topologiesEndpoint) context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var regionName string

		if regionName = chi.URLParam(r, "regionName"); regionName == "" {
			httpWriteResponseError(w, errors.New("region not found"))
			return
		}

		ctx := context.WithValue(r.Context(), "region-name", regionName) //nolint:staticcheck
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
