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

	"github.com/rasorp/attila/internal/server/state"
)

type JobRegisterRuleCreateResp struct {
	Rule                 *state.JobRegisterRule `json:"rule"`
	internalResponseMeta `json:"-"`
}

type JobRegisterRuleDeleteResp struct {
	internalResponseMeta `json:"-"`
}

type JobRegisterRuleGetResp struct {
	Rule                 *state.JobRegisterRule `json:"rule"`
	internalResponseMeta `json:"-"`
}

type JobRegisterRuleListResp struct {
	Rules                []*state.JobRegisterRuleStub `json:"rules"`
	internalResponseMeta `json:"-"`
}

type jobsRegisterRulesEndpoint struct {
	state state.State
}

func (j jobsRegisterRulesEndpoint) routes() chi.Router {
	r := chi.NewRouter()

	// Add the root endpoints which do not include a rule name within the
	// URI.
	r.Route("/", func(r chi.Router) {
		r.Get("/", j.list)
		r.Post("/", j.create)
	})

	r.Route("/{ruleName}", func(r chi.Router) {
		r.Use(j.context)
		r.Delete("/", j.delete)
		r.Get("/", j.get)
	})

	return r
}

func (j jobsRegisterRulesEndpoint) create(w http.ResponseWriter, r *http.Request) {

	var ruleObj state.JobRegisterRule

	if err := json.NewDecoder(r.Body).Decode(&ruleObj); err != nil {
		httpWriteResponseError(w, NewResponseError(fmt.Errorf("failed to decode object: %w", err), 400))
		return
	}

	if err := ruleObj.Validate(); err != nil {
		respErr := NewResponseError(err, http.StatusBadRequest)
		httpWriteResponseError(w, respErr)
		return
	}

	ruleObj.Metadata = state.NewMetadata()

	stateReq := state.JobRegisterRuleCreateReq{Rule: &ruleObj}

	ruleCreateResp, err := j.state.JobRegister().Rule().Create(&stateReq)
	if err != nil {
		respErr := NewResponseError(err.Err(), err.StatusCode())
		httpWriteResponseError(w, respErr)
	} else {
		resp := JobRegisterRuleCreateResp{
			Rule:                 ruleCreateResp.Rule,
			internalResponseMeta: newInternalResponseMeta(http.StatusCreated),
		}
		httpWriteResponse(w, &resp)
	}
}

func (j jobsRegisterRulesEndpoint) delete(w http.ResponseWriter, r *http.Request) {
	ruleName := r.Context().Value("rule-name").(string)

	stateReq := state.JobRegisterRuleDeleteReq{Name: ruleName}

	_, err := j.state.JobRegister().Rule().Delete(&stateReq)
	if err != nil {
		respErr := NewResponseError(err.Err(), err.StatusCode())
		httpWriteResponseError(w, respErr)
	} else {
		resp := JobRegisterMethodDeleteResp{
			internalResponseMeta: newInternalResponseMeta(http.StatusNoContent),
		}
		httpWriteResponse(w, &resp)
	}
}

func (j jobsRegisterRulesEndpoint) get(w http.ResponseWriter, r *http.Request) {
	ruleName := r.Context().Value("rule-name").(string)

	stateReq := state.JobRegisterRuleGetReq{Name: ruleName}

	ruleGetResp, err := j.state.JobRegister().Rule().Get(&stateReq)
	if err != nil {
		respErr := NewResponseError(err.Err(), err.StatusCode())
		httpWriteResponseError(w, respErr)
	} else {
		resp := JobRegisterRuleGetResp{
			Rule:                 ruleGetResp.Rule,
			internalResponseMeta: newInternalResponseMeta(http.StatusOK),
		}
		httpWriteResponse(w, &resp)
	}
}

func (j jobsRegisterRulesEndpoint) list(w http.ResponseWriter, r *http.Request) {
	ruleListResp, err := j.state.JobRegister().Rule().List(&state.JobRegisterRuleListReq{})
	if err != nil {
		respErr := NewResponseError(err.Err(), err.StatusCode())
		httpWriteResponseError(w, respErr)
	} else {
		resp := JobRegisterRuleListResp{
			Rules:                make([]*state.JobRegisterRuleStub, len(ruleListResp.Rules)),
			internalResponseMeta: newInternalResponseMeta(http.StatusOK),
		}

		for i, rule := range ruleListResp.Rules {
			resp.Rules[i] = rule.Stub()
		}

		httpWriteResponse(w, &resp)
	}
}

func (j jobsRegisterRulesEndpoint) context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var ruleName string

		if ruleName = chi.URLParam(r, "ruleName"); ruleName == "" {
			httpWriteResponseError(w, errors.New("job register rule not found"))
			return
		}

		ctx := context.WithValue(r.Context(), "rule-name", ruleName) //nolint:staticcheck
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
