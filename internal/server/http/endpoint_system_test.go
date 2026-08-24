// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shoenig/test/must"
)

func Test_healthEndpoint_get(t *testing.T) {
	t.Run("healthy status", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		systemEndpoint{}.routes().ServeHTTP(w, req)

		resp := w.Result()
		must.Eq(t, http.StatusOK, resp.StatusCode)
		must.StrEqFold(t, "application/json", w.Result().Header.Get("Content-Type"))

		var body SystemHealthResp
		must.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		must.Eq(t, SystemHealthStatusHealthy, body.Status)
	})

	t.Run("get only", func(t *testing.T) {
		for _, httpMethod := range []string{
			http.MethodHead,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodConnect,
			http.MethodOptions,
			http.MethodTrace,
		} {
			req := httptest.NewRequest(httpMethod, "/health", nil)
			w := httptest.NewRecorder()

			systemEndpoint{}.routes().ServeHTTP(w, req)
			must.Eq(t, http.StatusMethodNotAllowed, w.Result().StatusCode)
		}
	})
}
