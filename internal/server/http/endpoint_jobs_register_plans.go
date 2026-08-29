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
	"go.uber.org/zap"

	"github.com/rasorp/attila/internal/domain"
	"github.com/rasorp/attila/internal/server/nomad"
)

type JobsRegisterPlansCreateReq struct {
	Job *api.Job `json:"job"`
}

type JobsRegisterPlansCreateResp struct {
	Plan                 *domain.JobRegisterPlan `json:"plan"`
	internalResponseMeta `json:"-"`
}

type JobsRegisterPlansDeleteResp struct {
	internalResponseMeta `json:"-"`
}

type JobsRegisterPlansGetResp struct {
	Plan                 *domain.JobRegisterPlan `json:"plan"`
	internalResponseMeta `json:"-"`
}

type JobsRegisterPlansListResp struct {
	Plans                []*domain.JobRegisterPlan `json:"plans"`
	internalResponseMeta `json:"-"`
}

type JobsRegisterPlansRunResp struct {
	Run                  *domain.JobRegisterPlanRun `json:"run"`
	PatrialFailureError  error                      `json:"partial_failure_error"`
	internalResponseMeta `json:"-"`
}

type jobsRegisterPlansEndpoint struct {
	logger          *zap.Logger
	nomadController nomad.Controller
	state           domain.State
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
		httpWriteResponseError(
			w,
			NewResponseError(fmt.Errorf("failed to decode object: %w", err), http.StatusBadRequest))
		return
	}

	controllerResp, err := j.nomadController.JobRegistrationPlanCreate(
		&nomad.JobRegistrationPlanCreateReq{
			Job:       req.Job,
			Namespace: reqNamespace(r),
		},
	)
	if err != nil {
		httpWriteResponseError(w, NewResponseError(err, http.StatusInternalServerError))
		return
	}

	stateResp, stateErr := j.state.JobRegister().Plan().Create(&domain.JobRegisterPlanCreateReq{Plan: controllerResp.Plan})
	if stateErr != nil {
		httpWriteResponseError(w, NewResponseError(stateErr.Err(), stateErr.StatusCode()))
		return
	}

	httpWriteResponse(w, &JobsRegisterPlansCreateResp{
		Plan:                 stateResp.Plan,
		internalResponseMeta: newInternalResponseMeta(http.StatusCreated),
	})
}

func (j jobsRegisterPlansEndpoint) delete(w http.ResponseWriter, r *http.Request) {
	planID := r.Context().Value("id").(ulid.ULID)

	stateReq := domain.JobRegisterPlanDeleteReq{
		ID:        planID,
		Namespace: reqNamespace(r),
	}

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

	stateReq := domain.JobRegisterPlanGetReq{
		ID:        planID,
		Namespace: reqNamespace(r),
	}

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

	stateResp, err := j.state.JobRegister().Plan().List(&domain.JobRegisterPlanListReq{Namespace: reqNamespace(r)})
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

	planID := r.Context().Value("id").(ulid.ULID)
	requestNS := reqNamespace(r)

	planResp, err := j.state.JobRegister().Plan().Get(
		&domain.JobRegisterPlanGetReq{
			ID:        planID,
			Namespace: requestNS,
		},
	)
	if err != nil {
		httpWriteResponseError(w, NewResponseError(err.Err(), err.StatusCode()))
		return
	}

	controllerReq := nomad.JobRegistrationPlanRunReq{Plan: planResp.Plan}

	result, runErr := j.nomadController.JobRegistrationPlanRun(&controllerReq)
	if runErr != nil && result == nil {
		httpWriteResponseError(w, NewResponseError(runErr, http.StatusInternalServerError))
		return
	}

	responseCode := http.StatusCreated
	if runErr != nil {
		responseCode = http.StatusInternalServerError
	}

	stateReq := domain.JobRegisterPlanDeleteReq{ID: planID, Namespace: requestNS}

	if _, err := j.state.JobRegister().Plan().Delete(&stateReq); err != nil {
		j.logger.Error("failed to delete job register plan", zap.Error(err))
	}

	httpWriteResponse(w, &JobsRegisterPlansRunResp{
		Run:                  result.Run,
		PatrialFailureError:  runErr,
		internalResponseMeta: newInternalResponseMeta(responseCode),
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
