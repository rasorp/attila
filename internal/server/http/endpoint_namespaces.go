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
	"github.com/rasorp/attila/internal/server/state"
	"github.com/rasorp/attila/internal/store"
)

type NamespaceCreateReq struct {
	Namespace *domain.Namespace `json:"namespace"`
}

type NamespaceCreateResp struct {
	Namespace            *domain.Namespace `json:"namespace"`
	internalResponseMeta `json:"-"`
}

type NamespaceDeleteResp struct {
	internalResponseMeta `json:"-"`
}

type NamespaceGetResp struct {
	Namespace            *domain.Namespace `json:"namespace"`
	internalResponseMeta `json:"-"`
}

type NamespaceListResp struct {
	Namespaces           []*domain.NamespaceStub `json:"namespaces"`
	internalResponseMeta `json:"-"`
}

type namespacesEndpoint struct {
	state store.State
}

func (n namespacesEndpoint) routes() chi.Router {
	r := chi.NewRouter()

	// Add the root endpoints which do not include a namespace name within the URI.
	r.Route("/", func(r chi.Router) {
		r.Get("/", n.list)
		r.Post("/", n.create)
	})

	// Add the endpoints which are specific to a named namespace.
	r.Route("/{namespaceName}", func(r chi.Router) {
		r.Use(n.context)
		r.Delete("/", n.delete)
		r.Get("/", n.get)
	})

	return r
}

func (n namespacesEndpoint) create(w http.ResponseWriter, r *http.Request) {
	var req NamespaceCreateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpWriteResponseError(w, NewResponseError(fmt.Errorf("failed to decode object: %w", err), http.StatusBadRequest))
		return
	}

	if err := req.Namespace.Validate(); err != nil {
		respErr := NewResponseError(err, http.StatusBadRequest)
		httpWriteResponseError(w, respErr)
		return
	}

	// Validate region references against registered regions.
	regionListResp, err := n.state.Region().List(nil)
	if err != nil {
		httpWriteResponseError(w, NewResponseError(err.Err(), err.StatusCode()))
		return
	}

	regionNames := make(map[string]struct{}, len(regionListResp.Regions))
	for _, region := range regionListResp.Regions {
		regionNames[region.Name] = struct{}{}
	}
	for _, ref := range req.Namespace.Regions {
		if ref == domain.NamespaceRegionsWildcard {
			continue
		}
		if _, ok := regionNames[ref]; !ok {
			httpWriteResponseError(
				w,
				NewResponseError(fmt.Errorf("region %q not found", ref), http.StatusBadRequest),
			)
			return
		}
	}

	stateReq := state.NamespaceCreateReq{Namespace: req.Namespace}

	stateResp, stateErr := n.state.Namespace().Create(&stateReq)
	if stateErr != nil {
		respErr := NewResponseError(stateErr.Err(), stateErr.StatusCode())
		httpWriteResponseError(w, respErr)
	} else {
		resp := NamespaceCreateResp{
			Namespace:            stateResp.Namespace,
			internalResponseMeta: newInternalResponseMeta(http.StatusCreated),
		}
		httpWriteResponse(w, &resp)
	}
}

func (n namespacesEndpoint) delete(w http.ResponseWriter, r *http.Request) {
	namespaceName := r.Context().Value("namespace-name").(string)

	// Block any request to delete the default namespace.
	if namespaceName == domain.NamespaceDefaultName {
		httpWriteResponseError(
			w,
			NewResponseError(errors.New("cannot delete default namespace"), http.StatusBadRequest))
		return
	}

	_, err := n.state.Namespace().Delete(&state.NamespaceDeleteReq{Name: namespaceName})
	if err != nil {
		respErr := NewResponseError(err.Err(), err.StatusCode())
		httpWriteResponseError(w, respErr)
	} else {
		resp := NamespaceDeleteResp{
			internalResponseMeta: newInternalResponseMeta(http.StatusNoContent),
		}
		httpWriteResponse(w, &resp)
	}
}

func (n namespacesEndpoint) get(w http.ResponseWriter, r *http.Request) {
	namespaceName := r.Context().Value("namespace-name").(string)

	stateReq := state.NamespaceGetReq{Name: namespaceName}

	namespaceGetResp, err := n.state.Namespace().Get(&stateReq)
	if err != nil {
		respErr := NewResponseError(err.Err(), err.StatusCode())
		httpWriteResponseError(w, respErr)
	} else {
		resp := NamespaceGetResp{
			Namespace:            namespaceGetResp.Namespace,
			internalResponseMeta: newInternalResponseMeta(http.StatusOK),
		}
		httpWriteResponse(w, &resp)
	}
}

func (n namespacesEndpoint) list(w http.ResponseWriter, r *http.Request) {
	namespaceListResp, err := n.state.Namespace().List(&state.NamespaceListReq{})
	if err != nil {
		respErr := NewResponseError(err.Err(), err.StatusCode())
		httpWriteResponseError(w, respErr)
	} else {
		resp := NamespaceListResp{
			Namespaces:           make([]*domain.NamespaceStub, len(namespaceListResp.Namespaces)),
			internalResponseMeta: newInternalResponseMeta(http.StatusOK),
		}

		for i, ns := range namespaceListResp.Namespaces {
			resp.Namespaces[i] = ns.Stub()
		}

		httpWriteResponse(w, &resp)
	}
}

func (n namespacesEndpoint) context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var namespaceName string

		if namespaceName = chi.URLParam(r, "namespaceName"); namespaceName == "" {
			httpWriteResponseError(w, errors.New("namespace not found"))
			return
		}

		ctx := context.WithValue(r.Context(), "namespace-name", namespaceName) //nolint:staticcheck
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
