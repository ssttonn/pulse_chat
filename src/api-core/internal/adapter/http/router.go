package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter khởi tạo và cấu hình các endpoints cho API
func NewRouter(userHandler *UserHandler) *chi.Mux {
	r := chi.NewRouter()

	// Cài đặt các Middleware cơ bản
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Định nghĩa API Healthcheck
	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Post("/v1/users/register", userHandler.Register)

	return r
}
