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

	"github.com/rasorp/attila/internal/domain"
	"github.com/rasorp/attila/internal/store"
)

type JobRegisterMethodCreateResp struct {
	Method               *domain.JobRegisterMethod `json:"method"`
	internalResponseMeta `json:"-"`
}

type JobRegisterMethodDeleteResp struct {
	internalResponseMeta `json:"-"`
}

type JobRegisterMethodGetResp struct {
	Method               *domain.JobRegisterMethod `json:"method"`
	internalResponseMeta `json:"-"`
}

type JobRegisterMethodListResp struct {
	Methods              []*domain.JobRegisterMethodStub `json:"methods"`
	internalResponseMeta `json:"-"`
}

type jobsRegisterMethodsEndpoint struct {
	state store.State
}

func (j jobsRegisterMethodsEndpoint) routes() chi.Router {
	r := chi.NewRouter()

	// Add the root endpoints which do not include a method name within the
	// URI.
	r.Route("/", func(r chi.Router) {
		r.Get("/", j.list)
		r.Post("/", j.create)
	})

	r.Route("/{methodName}", func(r chi.Router) {
		r.Use(j.context)
		r.Delete("/", j.delete)
		r.Get("/", j.get)
	})

	return r
}

func (j jobsRegisterMethodsEndpoint) create(w http.ResponseWriter, r *http.Request) {
	var methodObj domain.JobRegisterMethod

	if err := json.NewDecoder(r.Body).Decode(&methodObj); err != nil {
		httpWriteResponseError(w, NewResponseError(fmt.Errorf("failed to decode object: %w", err), http.StatusBadRequest))
		return
	}

	if err := methodObj.Validate(); err != nil {
		respErr := NewResponseError(err, http.StatusBadRequest)
		httpWriteResponseError(w, respErr)
		return
	}

	methodObj.Metadata = domain.NewMetadata()

	stateReq := store.JobRegisterMethodCreateReq{Method: &methodObj}

	methodCreateResp, err := j.state.JobRegister().Method().Create(&stateReq)
	if err != nil {
		respErr := NewResponseError(err.Err(), err.StatusCode())
		httpWriteResponseError(w, respErr)
	} else {
		resp := JobRegisterMethodCreateResp{
			Method:               methodCreateResp.Method,
			internalResponseMeta: newInternalResponseMeta(http.StatusCreated),
		}
		httpWriteResponse(w, &resp)
	}
}

func (j jobsRegisterMethodsEndpoint) delete(w http.ResponseWriter, r *http.Request) {
	methodName := r.Context().Value("method-name").(string)

	stateReq := store.JobRegisterMethodDeleteReq{Name: methodName}

	_, err := j.state.JobRegister().Method().Delete(&stateReq)
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

func (a jobsRegisterMethodsEndpoint) get(w http.ResponseWriter, r *http.Request) {
	methodName := r.Context().Value("method-name").(string)

	stateReq := store.JobRegisterMethodGetReq{Name: methodName}

	methodGetResp, err := a.state.JobRegister().Method().Get(&stateReq)
	if err != nil {
		respErr := NewResponseError(err.Err(), err.StatusCode())
		httpWriteResponseError(w, respErr)
	} else {
		resp := JobRegisterMethodGetResp{
			Method:               methodGetResp.Method,
			internalResponseMeta: newInternalResponseMeta(http.StatusOK),
		}
		httpWriteResponse(w, &resp)
	}
}

func (j jobsRegisterMethodsEndpoint) list(w http.ResponseWriter, r *http.Request) {
	stateResp, err := j.state.JobRegister().Method().List(&store.JobRegisterMethodListReq{})
	if err != nil {
		respErr := NewResponseError(err.Err(), err.StatusCode())
		httpWriteResponseError(w, respErr)
	} else {
		resp := JobRegisterMethodListResp{
			Methods:              make([]*domain.JobRegisterMethodStub, len(stateResp.Methods)),
			internalResponseMeta: newInternalResponseMeta(http.StatusOK),
		}

		for i, method := range stateResp.Methods {
			resp.Methods[i] = method.Stub()
		}

		httpWriteResponse(w, &resp)
	}
}

func (j jobsRegisterMethodsEndpoint) context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var methodName string

		if methodName = chi.URLParam(r, "methodName"); methodName == "" {
			httpWriteResponseError(w, errors.New("job register method not found"))
			return
		}

		ctx := context.WithValue(r.Context(), "method-name", methodName) //nolint:staticcheck
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
