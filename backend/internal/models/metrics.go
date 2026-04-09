package models

import (
	"net/http"
	"time"
)

const (
	CPU uint8 = iota
	Memory
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func ObserverRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		_ = duration
	})
}

type OptimizedMetric struct {
	ID    uint32
	Value float64
	Type  uint8
}
