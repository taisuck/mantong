package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/taisuck/mantong/src/mongodb"
)

const stockQueryTimeout = 5 * time.Second

// stockResponse 는 GET /stock 응답 항목이다. (src/api/openapi.yaml 의 Stock 스키마)
type stockResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type stockHandler struct {
	repo *mongodb.StockRepository
}

// ServeHTTP 는 현재 수집 중인 주식의 목록을 반환한다.
func (h *stockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), stockQueryTimeout)
	defer cancel()

	stocks, err := h.repo.FindAll(ctx)
	if err != nil {
		slog.Error("stock 목록 조회 실패", "method", r.Method, "path", r.URL.Path, "err", err)
		writeError(w, r, http.StatusInternalServerError, "stock 목록을 조회하지 못했습니다")
		return
	}

	resp := make([]stockResponse, 0, len(stocks))
	for _, s := range stocks {
		resp = append(resp, stockResponse{ID: s.ID.Hex(), Name: s.Name, Code: s.Code})
	}
	writeJSON(w, r, http.StatusOK, resp)
}
