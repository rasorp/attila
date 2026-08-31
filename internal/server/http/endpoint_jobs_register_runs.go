// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hashicorp/nomad/api"
	"github.com/oklog/ulid/v2"

	"github.com/rasorp/attila/internal/domain"
	"github.com/rasorp/attila/internal/server/nomad"
)

// JobRegisterPlanRunCreateResp returns the created plan run record.
type JobRegisterRunDeleteResp struct {
	internalResponseMeta `json:"-"`
}

type JobsRegisterRunsCreateReq struct {
	PlanID ulid.ULID `json:"plan_id"`
	Job    *api.Job  `json:"job"`
}

type JobsRegisterRunsGetResp struct {
	Run                  *domain.JobRegisterRun `json:"run"`
	internalResponseMeta `json:"-"`
}

type JobsRegisterRunsCreateResp struct {
	Run                  *domain.JobRegisterRun `json:"run"`
	internalResponseMeta `json:"-"`
}

type JobRegisterRunsListResp struct {
	Runs                 []*domain.JobRegisterRunStub `json:"runs"`
	internalResponseMeta `json:"-"`
}

type jobsRegisterRunsEndpoint struct {
	nomadController nomad.Controller
	state           domain.State
}

func (j jobsRegisterRunsEndpoint) routes() chi.Router {
	r := chi.NewRouter()

	// List all runs for a namespace.
	r.Get("/", j.list)
	r.Post("/", j.create)

	// Actions specific to a run identified by it's ID in the URL.
	r.Route("/{id}", func(r chi.Router) {
		r.Use(j.context)
		r.Delete("/", j.delete)
		r.Get("/", j.get)
	})

	return r
}

func (j jobsRegisterRunsEndpoint) create(w http.ResponseWriter, r *http.Request) {

	var req JobsRegisterRunsCreateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpWriteResponseError(
			w,
			NewResponseError(fmt.Errorf("failed to decode object: %w", err), http.StatusBadRequest))
		return
	}

	requestNS := reqNamespace(r)

	var plan *domain.JobRegisterPlan
	var runErr error

	switch {
	case !req.PlanID.IsZero():
		planResp, err := j.state.JobRegister().Plan().Get(
			&domain.JobRegisterPlanGetReq{
				ID:        req.PlanID,
				Namespace: requestNS,
			},
		)
		if err != nil {
			httpWriteResponseError(w, NewResponseError(err.Err(), err.StatusCode()))
			return
		}
		plan = planResp.Plan

	case req.Job != nil:
		planResp, err := j.nomadController.JobRegistrationPlanCreate(
			&nomad.JobRegistrationPlanCreateReq{
				Job:       req.Job,
				Namespace: requestNS,
			},
		)
		if err != nil {
			httpWriteResponseError(w, NewResponseError(err, http.StatusInternalServerError))
			return
		}

		stateResp, stateErr := j.state.JobRegister().Plan().Create(
			&domain.JobRegisterPlanCreateReq{
				Plan: planResp.Plan,
			},
		)
		if stateErr != nil {
			httpWriteResponseError(w, NewResponseError(stateErr.Err(), stateErr.StatusCode()))
			return
		}
		plan = stateResp.Plan

	default:
		httpWriteResponseError(
			w,
			NewResponseError(errors.New("either plan-id or jobspec must be provided"), http.StatusBadRequest))
		return
	}

	// Execute the run.
	controllerReq := nomad.JobRegistrationPlanRunReq{Plan: plan}
	result, runErr := j.nomadController.JobRegistrationPlanRun(&controllerReq)
	if runErr != nil && result == nil {
		httpWriteResponseError(w, NewResponseError(runErr, http.StatusInternalServerError))
		return
	}

	responseCode := http.StatusCreated
	if runErr != nil {
		responseCode = http.StatusInternalServerError
	}

	if _, err := j.state.JobRegister().Run().Create(
		&domain.JobRegisterRunCreateReq{
			Run: result.Run,
		},
	); err != nil {
		httpWriteResponseError(w, NewResponseError(err.Err(), err.StatusCode()))
		return
	}

	httpWriteResponse(w, &JobsRegisterRunsCreateResp{
		Run:                  result.Run,
		internalResponseMeta: newInternalResponseMeta(responseCode),
	})
}

func (j jobsRegisterRunsEndpoint) get(w http.ResponseWriter, r *http.Request) {

	stateReq := domain.JobRegisterRunGetReq{
		ID:        r.Context().Value("id").(ulid.ULID),
		Namespace: reqNamespace(r),
	}

	stateResp, err := j.state.JobRegister().Run().Get(&stateReq)
	if err != nil {
		httpWriteResponseError(w, NewResponseError(err.Err(), err.StatusCode()))
		return
	}

	httpWriteResponse(w, &JobsRegisterRunsGetResp{
		Run:                  stateResp.Run,
		internalResponseMeta: newInternalResponseMeta(http.StatusOK),
	})
}

func (j jobsRegisterRunsEndpoint) delete(w http.ResponseWriter, r *http.Request) {

	stateReq := domain.JobRegisterRunDeleteReq{
		ID:        r.Context().Value("id").(ulid.ULID),
		Namespace: reqNamespace(r),
	}

	_, err := j.state.JobRegister().Run().Delete(&stateReq)
	if err != nil {
		httpWriteResponseError(w, NewResponseError(err.Err(), err.StatusCode()))
	} else {
		resp := JobRegisterRunDeleteResp{
			internalResponseMeta: newInternalResponseMeta(http.StatusNoContent),
		}
		httpWriteResponse(w, &resp)
	}
}

func (j jobsRegisterRunsEndpoint) list(w http.ResponseWriter, r *http.Request) {
	stateResp, err := j.state.JobRegister().Run().List(
		&domain.JobRegisterRunListReq{Namespace: reqNamespace(r)},
	)
	if err != nil {
		httpWriteResponseError(w, NewResponseError(err.Err(), err.StatusCode()))
		return
	}

	resp := JobRegisterRunsListResp{
		Runs:                 make([]*domain.JobRegisterRunStub, len(stateResp.Runs)),
		internalResponseMeta: newInternalResponseMeta(http.StatusOK),
	}

	for i, run := range stateResp.Runs {
		resp.Runs[i] = run.Stub()
	}

	httpWriteResponse(w, &resp)
}

func (j jobsRegisterRunsEndpoint) context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var idString string

		if idString = chi.URLParam(r, "id"); idString == "" {
			httpWriteResponseError(w, errors.New("id not found"))
			return
		}

		if idULID, err := ulid.Parse(idString); err != nil {
			httpWriteResponseError(w, fmt.Errorf("failed to parse ID: %w", err))
		} else {
			ctx := context.WithValue(r.Context(), "id", idULID) //nolint:staticcheck
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
