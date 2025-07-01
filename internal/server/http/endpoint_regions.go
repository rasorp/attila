// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/rasorp/attila/internal/server/nomad"
	"github.com/rasorp/attila/internal/server/state"
)

type RegionCreateReq struct {
	Region *state.Region `json:"region"`
}

type RegionCreateResp struct {
	Region               *state.Region `json:"region"`
	internalResponseMeta `json:"-"`
}

type RegionDeleteResp struct {
	internalResponseMeta `json:"-"`
}

type RegionGetResp struct {
	Region               *state.Region `json:"region"`
	internalResponseMeta `json:"-"`
}

type RegionListResp struct {
	Regions              []*state.RegionStub `json:"regions"`
	internalResponseMeta `json:"-"`
}

type regionsEndpoint struct {
	state           state.State
	nomadController nomad.Controller
}

func (a regionsEndpoint) routes() chi.Router {
	r := chi.NewRouter()

	// Add the root endpoints which do not include a region name within the URI.
	r.Route("/", func(r chi.Router) {
		r.Get("/", a.list)
		r.Post("/", a.create)
	})

	// Add the endpoints which are specific to a named region.
	r.Route("/{regionName}", func(r chi.Router) {
		r.Use(a.context)
		r.Delete("/", a.delete)
		r.Get("/", a.get)
	})

	return r
}

func (a regionsEndpoint) create(w http.ResponseWriter, r *http.Request) {

	var req state.RegionCreateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpWriteResponseError(w, NewResponseError(fmt.Errorf("failed to decode object: %w", err), 400))
		return
	}

	req.Region.SetDefaults()

	if err := req.Region.Validate(); err != nil {
		respErr := NewResponseError(err, http.StatusBadRequest)
		httpWriteResponseError(w, respErr)
		return
	}

	nomadClient, clientErr := req.Region.GenerateNomadClient()
	if clientErr != nil {
		respErr := NewResponseError(clientErr, http.StatusBadRequest)
		httpWriteResponseError(w, respErr)
	}

	req.Region.Metadata = state.NewMetadata()

	stateReq := state.RegionCreateReq{Region: req.Region}

	stateResp, err := a.state.Region().Create(&stateReq)
	if err != nil {
		respErr := NewResponseError(err.Err(), err.StatusCode())
		httpWriteResponseError(w, respErr)
	} else {
		a.nomadController.RegionClientSet(stateResp.Region.Name, nomadClient)
		resp := RegionCreateResp{
			Region:               stateResp.Region,
			internalResponseMeta: newInternalResponseMeta(http.StatusCreated),
		}
		httpWriteResponse(w, &resp)
	}
}

func (a regionsEndpoint) delete(w http.ResponseWriter, r *http.Request) {
	regionName := r.Context().Value("region-name").(string)

	stateReq := state.RegionDeleteReq{RegionName: regionName}

	_, err := a.state.Region().Delete(&stateReq)
	if err != nil {
		respErr := NewResponseError(err.Err(), err.StatusCode())
		httpWriteResponseError(w, respErr)
	} else {
		a.nomadController.RegionClientDelete(regionName)
		resp := RegionDeleteResp{
			internalResponseMeta: newInternalResponseMeta(http.StatusNoContent),
		}
		httpWriteResponse(w, &resp)
	}
}

func (a regionsEndpoint) get(w http.ResponseWriter, r *http.Request) {
	regionName := r.Context().Value("region-name").(string)

	stateReq := state.RegionGetReq{RegionName: regionName}

	regionGetResp, err := a.state.Region().Get(&stateReq)
	if err != nil {
		respErr := NewResponseError(err.Err(), err.StatusCode())
		httpWriteResponseError(w, respErr)
	} else {
		resp := RegionGetResp{
			Region:               regionGetResp.Region,
			internalResponseMeta: newInternalResponseMeta(http.StatusOK),
		}
		httpWriteResponse(w, &resp)
	}
}

func (a regionsEndpoint) list(w http.ResponseWriter, r *http.Request) {
	regionListResp, err := a.state.Region().List(&state.RegionListReq{})
	if err != nil {
		respErr := NewResponseError(err.Err(), err.StatusCode())
		httpWriteResponseError(w, respErr)
	} else {
		resp := RegionListResp{
			Regions:              make([]*state.RegionStub, len(regionListResp.Regions)),
			internalResponseMeta: newInternalResponseMeta(http.StatusOK),
		}

		for i, region := range regionListResp.Regions {
			resp.Regions[i] = region.Stub()
		}

		httpWriteResponse(w, &resp)
	}
}

func (a regionsEndpoint) context(next http.Handler) http.Handler {
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
