/*
Copyright (c) JSC iCore.

This source code is licensed under the MIT license found in the
LICENSE file in the root directory of this source tree.
*/

// Package rlog provides access to logging service from HTTP request's context.Context.
package rlog // import "github.com/i-core/rlog"

import (
	"context"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type requestLogCtxKey int

// requestLogKey is a context's key to store a request's logger.
const requestLogKey requestLogCtxKey = iota

type traceResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *traceResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// NewMiddleware returns a middleware that places a logger with request's ID to context, and logs the request.
func NewMiddleware(log *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				log = log.With(zap.String("requestID", uuid.Must(uuid.NewV4()).String()))
				ctx = context.WithValue(r.Context(), requestLogKey, log)
			)
			log.Info("New request",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()))
			start := time.Now()
			tw := &traceResponseWriter{w, http.StatusOK}
			next.ServeHTTP(tw, r.WithContext(ctx))
			log.Debug("The request is handled",
				zap.Int("httpStatus", tw.statusCode),
				zap.Duration("duration", time.Since(start)))
		})
	}
}

// FromContext returns a request's logger stored in a context.
func FromContext(ctx context.Context) *zap.Logger {
	v := ctx.Value(requestLogKey)
	if v == nil {
		return zap.NewNop()
	}
	return v.(*zap.Logger)
}
