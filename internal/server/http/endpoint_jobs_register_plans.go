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
	"github.com/hashicorp/nomad/api"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"

	"github.com/rasorp/attila/internal/server/nomad"
	"github.com/rasorp/attila/internal/server/state"
)

type JobsRegisterPlansCreateReq struct {
	Job *api.Job `json:"job"`
}

type JobsRegisterPlansCreateResp struct {
	Plan                 *state.JobRegisterPlan `json:"plan"`
	internalResponseMeta `json:"-"`
}

type JobsRegisterPlansDeleteResp struct {
	internalResponseMeta `json:"-"`
}

type JobsRegisterPlansGetResp struct {
	Plan                 *state.JobRegisterPlan `json:"plan"`
	internalResponseMeta `json:"-"`
}

type JobsRegisterPlansListResp struct {
	Plans                []*state.JobRegisterPlan `json:"plans"`
	internalResponseMeta `json:"-"`
}

type JobsRegisterPlansRunReq struct {
	Job *api.Job `json:"job"`
}

type JobsRegisterPlansRunResp struct {
	Run                  *state.JobRegisterPlanRun `json:"run"`
	PatrialFailureError  error                     `json:"partial_failure_error"`
	internalResponseMeta `json:"-"`
}

type jobsRegisterPlansEndpoint struct {
	logger          zerolog.Logger
	nomadController nomad.Controller
	state           state.State
}

func (j jobsRegisterPlansEndpoint) routes() chi.Router {
	r := chi.NewRouter()

	// Add the root endpoints which do not include a plan ID within the URI.
	r.Route("/", func(r chi.Router) {
		r.Get("/", j.list)
		r.Post("/", j.create)
	})

	// Add the endpoints which are specific to a job register plan using the ID.
	r.Route("/{id}", func(r chi.Router) {
		r.Use(j.context)
		r.Delete("/", j.delete)
		r.Get("/", j.get)
		r.Post("/run", j.run)
	})

	return r
}

func (j jobsRegisterPlansEndpoint) create(w http.ResponseWriter, r *http.Request) {

	var req JobsRegisterPlansCreateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpWriteResponseError(w,
			NewResponseError(fmt.Errorf("failed to decode object: %w", err), http.StatusBadRequest))
		return
	}

	controllerResp, err := j.nomadController.JobRegistrationPlanCreate(req.Job, j.state)
	if err != nil {
		httpWriteResponseError(w, NewResponseError(err, http.StatusInternalServerError))
		return
	}

	stateResp, errr := j.state.JobRegister().Plan().Create(&state.JobRegisterPlanCreateReq{Plan: controllerResp})
	if errr != nil {
		httpWriteResponseError(w, NewResponseError(err, http.StatusInternalServerError))
		return
	}

	httpWriteResponse(w, &JobsRegisterPlansCreateResp{
		Plan:                 stateResp.Plan,
		internalResponseMeta: newInternalResponseMeta(http.StatusCreated),
	})
}

func (j jobsRegisterPlansEndpoint) delete(w http.ResponseWriter, r *http.Request) {
	planID := r.Context().Value("id").(ulid.ULID)

	stateReq := state.JobRegisterPlanDeleteReq{ID: planID}

	_, err := j.state.JobRegister().Plan().Delete(&stateReq)
	if err != nil {
		httpWriteResponseError(w, NewResponseError(err.Err(), err.StatusCode()))
	} else {
		httpWriteResponse(w, &JobsRegisterPlansDeleteResp{
			internalResponseMeta: newInternalResponseMeta(http.StatusNoContent),
		})
	}
}

func (j jobsRegisterPlansEndpoint) get(w http.ResponseWriter, r *http.Request) {
	planID := r.Context().Value("id").(ulid.ULID)

	stateReq := state.JobRegisterPlanGetReq{ID: planID}

	stateResp, err := j.state.JobRegister().Plan().Get(&stateReq)
	if err != nil {
		httpWriteResponseError(w, NewResponseError(err.Err(), err.StatusCode()))
	} else {
		httpWriteResponse(w, &JobsRegisterPlansGetResp{
			Plan:                 stateResp.Plan,
			internalResponseMeta: newInternalResponseMeta(http.StatusOK),
		})
	}
}

func (j jobsRegisterPlansEndpoint) list(w http.ResponseWriter, r *http.Request) {
	stateResp, err := j.state.JobRegister().Plan().List(&state.JobRegisterPlanListReq{})
	if err != nil {
		respErr := NewResponseError(err.Err(), err.StatusCode())
		httpWriteResponseError(w, respErr)
	} else {
		resp := JobsRegisterPlansListResp{
			Plans:                stateResp.Plans,
			internalResponseMeta: newInternalResponseMeta(http.StatusOK),
		}
		httpWriteResponse(w, &resp)
	}
}

func (j jobsRegisterPlansEndpoint) run(w http.ResponseWriter, r *http.Request) {

	var httpReq JobsRegisterPlansRunReq

	if err := json.NewDecoder(r.Body).Decode(&httpReq); err != nil {
		httpWriteResponseError(w,
			NewResponseError(fmt.Errorf("failed to decode object: %w", err), http.StatusBadRequest))
		return
	}

	planID := r.Context().Value("id").(ulid.ULID)

	result, err := j.nomadController.JobRegistrationRun(planID, httpReq.Job, j.state)
	if err != nil && result == nil {
		httpWriteResponseError(w, NewResponseError(err, http.StatusInternalServerError))
		return
	}

	respnseCode := http.StatusCreated
	if err != nil {
		respnseCode = http.StatusInternalServerError
	}

	stateReq := state.JobRegisterPlanDeleteReq{ID: planID}

	if _, err := j.state.JobRegister().Plan().Delete(&stateReq); err != nil {
		j.logger.Err(err).Msg("failed to delete job register plan")
	}

	httpWriteResponse(w, &JobsRegisterPlansRunResp{
		Run:                  result,
		PatrialFailureError:  err,
		internalResponseMeta: newInternalResponseMeta(respnseCode),
	})
}

func (j jobsRegisterPlansEndpoint) context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var planIDString string

		if planIDString = chi.URLParam(r, "id"); planIDString == "" {
			httpWriteResponseError(w, errors.New("id not found"))
			return
		}

		if planULID, err := ulid.Parse(planIDString); err != nil {
			httpWriteResponseError(w, fmt.Errorf("failed to parse ID: %w", err))
		} else {
			ctx := context.WithValue(r.Context(), "id", planULID) //nolint:staticcheck
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
