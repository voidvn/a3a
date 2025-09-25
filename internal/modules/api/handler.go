package api

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.Health(r.Context())
	if err != nil {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.Info(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// Вспомогательная функция ответа JSON.
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
