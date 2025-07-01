// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func loggerMiddleware(logger zerolog.Logger, accessLevel string) func(next http.Handler) http.Handler {

	var accessLoggerFn func() *zerolog.Event

	switch accessLevel {
	case zerolog.LevelTraceValue:
		accessLoggerFn = logger.Trace
	case zerolog.LevelDebugValue:
		accessLoggerFn = logger.Debug
	case zerolog.LevelInfoValue:
		accessLoggerFn = logger.Info
	default:
		panic(fmt.Sprintf("unsupported access log level: %q", accessLevel))
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			startTime := time.Now()

			defer func() {

				// Recover any panicking handler and log the stack trace, so
				// this is available for debugging. Although the output is a
				// little tricky to grok, it is very useful.
				if rec := recover(); rec != nil {
					logger.Error().
						Interface("recover_info", rec).
						Bytes("debug_stack", debug.Stack()).
						Msg("panic during handling of HTTP request")

					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				accessLoggerFn().
					Str("remote_address", r.RemoteAddr).
					Str("path", r.URL.Path).
					Str("proto", r.Proto).
					Str("method", r.Method).
					Str("user_agent", r.Header.Get("User-Agent")).
					Int("status", ww.Status()).
					Int64("latency_ns", int64(time.Since(startTime).Nanoseconds())).
					Int("content_in_bytes", contentInBytes(r.Header)).
					Int("content_out_bytes", ww.BytesWritten()).
					Msg("successfully handled HTTP request")
			}()

			next.ServeHTTP(ww, r)

		}
		return http.HandlerFunc(fn)
	}
}

func contentInBytes(header http.Header) int {
	if i, err := strconv.Atoi(header.Get("Content-Length")); err != nil {
		return 0
	} else {
		return i
	}
}
