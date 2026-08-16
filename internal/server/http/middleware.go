// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func loggerMiddleware(logger *zap.Logger, accessLevel string) func(next http.Handler) http.Handler {

	var level zapcore.Level
	switch accessLevel {
	case "trace":
		level = zapcore.DebugLevel
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
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
					logger.Error(
						"panic during handling of HTTP request",
						zap.Any("recover_info", rec),
						zap.ByteString("debug_stack", debug.Stack()),
					)
					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				// Check the log message would be written before writing out all the data.
				if entry := logger.Check(level, "successfully handled HTTP request"); entry != nil {
					entry.Write(
						zap.String("remote_address", r.RemoteAddr),
						zap.String("path", r.URL.Path),
						zap.String("proto", r.Proto),
						zap.String("method", r.Method),
						zap.String("user_agent", r.Header.Get("User-Agent")),
						zap.Int("status", ww.Status()),
						zap.Int64("latency_ns", int64(time.Since(startTime).Nanoseconds())),
						zap.Int("content_in_bytes", contentInBytes(r.Header)),
						zap.Int("content_out_bytes", ww.BytesWritten()),
					)
				}
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
