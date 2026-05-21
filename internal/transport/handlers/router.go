package handlers

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) NewRouter() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger, middleware.Recoverer)

	router.Route("/api/", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", h.CreateUser)
			r.Get("/{id}", h.GetUser)
			r.Patch("/{id}", h.UpdateUser)
			r.Delete("/{id}", h.DeleteUser)
		})
	})
	return router
}
