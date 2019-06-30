package routing

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func New(sources, feeds http.Handler) *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(middleware.AllowContentType("application/json"))

	mux.Get("/sources", sources.ServeHTTP)
	mux.Get("/feeds", feeds.ServeHTTP)
	return mux
}
