package routing

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func New(sources http.Handler) *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(middleware.AllowContentType("application/json"))

	mux.Get("/sources", sources.ServeHTTP)
	return mux
}
