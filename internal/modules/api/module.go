package api

import (
	"net/http"
)

// RegisterRoutes регистрирует эндпоинты модуля API.
// Пример использования:
//
//	mux := http.NewServeMux()
//	api.RegisterRoutes(mux)
func RegisterRoutes(mux *http.ServeMux) {
	svc := NewService(
	// передайте зависимости при необходимости, например: db.Client
	)
	h := NewHandler(svc)

	mux.HandleFunc("GET /api/health", h.Health)
	mux.HandleFunc("GET /api/info", h.Info)
}
