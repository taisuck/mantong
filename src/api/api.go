// Package api 는 웹 api request handler 를 제공한다.
package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/taisuck/mantong/src/mongodb"
)

// NewRouter 는 api 라우팅이 등록된 http.Handler 를 반환한다.
func NewRouter(stocks *mongodb.StockRepository) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /stock", &stockHandler{repo: stocks})
	return withLogging(mux)
}

type errorResponse struct {
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("응답 인코딩 실패", "method", r.Method, "path", r.URL.Path, "err", err)
	}
}

func writeError(w http.ResponseWriter, r *http.Request, status int, msg string) {
	writeJSON(w, r, status, errorResponse{Message: msg})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// withLogging 은 요청을 로깅하고, handler 에서 panic 이 발생하면 로그를 남기고 500 을 응답한다.
func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		defer func() {
			if p := recover(); p != nil {
				slog.Error("handler panic", "method", r.Method, "path", r.URL.Path, "panic", p)
				writeError(rec, r, http.StatusInternalServerError, "internal server error")
			}

			attrs := []any{"method", r.Method, "path", r.URL.Path, "status", rec.status, "elapsed", time.Since(start)}
			if rec.status >= http.StatusInternalServerError {
				slog.Error("request", attrs...)
			} else {
				slog.Info("request", attrs...)
			}
		}()

		next.ServeHTTP(rec, r)
	})
}
