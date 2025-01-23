package handler

import (
	app "github.com/Tairascii/google-docs-user/internal"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Handler struct {
	DI *app.DI
}

func NewHandler(DI *app.DI) *Handler {
	return &Handler{DI: DI}
}

func (h *Handler) InitHandlers() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/api", func(api chi.Router) {
		api.Route("/v1", func(v1 chi.Router) {
			v1.Mount("/user", handlers(h))
		})
	})
	return r
}

func handlers(h *Handler) http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Post("/sign-in", h.SignIn)
	})

	return rg
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {

}
