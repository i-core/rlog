/*
Copyright (c) JSC iCore.

This source code is licensed under the MIT license found in the
LICENSE file in the root directory of this source tree.
*/

package rlog_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/i-core/rlog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestFromContextWithMiddleware(t *testing.T) {
	var (
		out          bytes.Buffer
		originCalled bool
		encoderCfg   = zapcore.EncoderConfig{
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		}
		core  = zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.AddSync(&out), zap.DebugLevel)
		logmw = rlog.NewMiddleware(zap.New(core))
		rr    = httptest.NewRecorder()
		r     = httptest.NewRequest(http.MethodGet, "http://example.org", nil)
	)

	handler := logmw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := rlog.FromContext(r.Context())
		log.Info("handle request")
		originCalled = true
	}))
	handler.ServeHTTP(rr, r)

	var records []map[string]interface{}
	for _, raw := range strings.Split(strings.Trim(out.String(), "\n"), "\n") {
		var rec map[string]interface{}
		if err := json.Unmarshal([]byte(raw), &rec); err != nil {
			t.Fatalf("failed to unmarshal log records: %s", err)
		}
		records = append(records, rec)
	}

	if !originCalled {
		t.Errorf("origin HTTP handler is not called, want origin HTTP handler to be called")
	}
	if cnt := len(records); cnt != 3 {
		t.Errorf("got %d record(s), want 3 records", cnt)
	}
	for _, rec := range records {
		if i, ok := rec["requestID"]; !ok {
			t.Errorf("record %d does not have a requestID", i)
		}
	}
}

func TestFromContextWithoutMiddleware(t *testing.T) {
	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://example.org", nil)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if log := rlog.FromContext(r.Context()); log == nil {
			t.Errorf("got no logger, want logger")
		}
	})
	handler.ServeHTTP(rr, r)
}
