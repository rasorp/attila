// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func httpWriteResponse(w http.ResponseWriter, obj any) {

	code := http.StatusInternalServerError

	if respMeta, ok := obj.(internalResponseMeta); ok {
		code = respMeta.StatusCode()
	}

	switch code {
	case http.StatusNoContent:
		w.WriteHeader(code)
		return
	case http.StatusOK:
	default:
		w.WriteHeader(code)
	}

	// If we have a response object, encode it.
	if obj != nil {
		w.Header().Set("Content-Type", "application/json")

		objBytes, err := json.Marshal(obj)
		if err != nil {
			httpWriteResponseError(w, fmt.Errorf("failed to marshal JSON response: %w", err))
			return
		}

		if _, err := w.Write(objBytes); err != nil {
			httpWriteResponseError(w, fmt.Errorf("failed to write JSON response: %w", err))
			return
		}
	}
}

func httpWriteResponseError(w http.ResponseWriter, err error) {
	var (
		code int
		resp []byte
	)

	codedErr, ok := err.(*ResponseError)
	if !ok {
		code = http.StatusInternalServerError
		resp = []byte(err.Error())
	} else {
		code = codedErr.StatusCode()

		objBytes, err := json.Marshal(codedErr)
		if err != nil {
			return
		}
		resp = objBytes
		w.Header().Set("Content-Type", "application/json")
	}

	// Write the status code header.
	w.WriteHeader(code)
	_, _ = w.Write(resp)
}
