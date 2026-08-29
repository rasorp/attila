// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"
)

// queryParamNamespace is the query parameter added to requests, so we specify
// the Attila namespace to query or write to.
const queryParamNamespace = "namespace"

type QueryOpts struct {
	Namespace string
}

func (q *QueryOpts) SetOpts() func(req *http.Request) {
	return func(req *http.Request) {
		query := req.URL.Query()
		query.Set(queryParamNamespace, defaultNS(q.Namespace))
		req.URL.RawQuery = query.Encode()
	}
}

type WriteOpts struct {
	Namespace string
}

func (w *WriteOpts) SetOpts() func(req *http.Request) {
	return func(req *http.Request) {
		query := req.URL.Query()
		query.Set(queryParamNamespace, defaultNS(w.Namespace))
		req.URL.RawQuery = query.Encode()
	}
}

func defaultNS(ns string) string {
	if ns == "" {
		return NamespaceDefaultName
	}
	return ns
}
