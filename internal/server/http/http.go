// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"net/http"

	"github.com/rasorp/attila/internal/domain"
)

// requestNamespace pulls the namespace query parameter from the HTTP request
// object. If no parameter is present, we assume the target namespace is the
// default one.
func reqNamespace(r *http.Request) string {
	if ns := r.URL.Query().Get("namespace"); ns != "" {
		return ns
	} else {
		return domain.NamespaceDefaultName
	}
}

// namespacesMatch returns whether the two namespace identifiers match. This is
// used to ensure an objects namespace matches that request namespace.
func namespacesMatch(obj, req string) bool { return obj == req }
