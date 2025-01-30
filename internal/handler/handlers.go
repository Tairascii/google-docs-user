package handler

import (
	"encoding/json"
	"errors"
	app "github.com/Tairascii/google-docs-user/internal"
	"github.com/Tairascii/google-docs-user/internal/usecase"
	"github.com/Tairascii/google-docs-user/pkg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

var (
	ErrBadCredentials    = errors.New("bad credentials")
	ErrUserAlreadyExists = errors.New("user already exists")
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
	rg.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))
	rg.Group(func(r chi.Router) {
		r.Post("/sign-in", h.SignIn)
		r.Post("/sign-up", h.SignUp)
	})

	return rg
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var payload SignInPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		pkg.JSONErrorResponseWriter(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ctx := r.Context()
	tokens, err := h.DI.UseCase.Auth.SignIn(ctx, payload.Email, payload.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			pkg.JSONErrorResponseWriter(w, ErrBadCredentials, http.StatusUnauthorized)
			return
		}
		pkg.JSONErrorResponseWriter(w, err, http.StatusInternalServerError)
		return
	}
	pkg.JSONResponseWriter[Tokens](w, toTokens(tokens), http.StatusOK)
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var payload SignUpPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		pkg.JSONErrorResponseWriter(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ctx := r.Context()
	tokens, err := h.DI.UseCase.Auth.SignUp(ctx, usecase.SignUpData{
		Name:          payload.Name,
		Email:         payload.Email,
		Password:      payload.Password,
		ProfilePicUrl: payload.ProfilePicUrl,
	})
	if err != nil {
		if errors.Is(err, usecase.ErrUserAlreadyExists) {
			pkg.JSONErrorResponseWriter(w, ErrUserAlreadyExists, http.StatusBadRequest)
			return
		}
		pkg.JSONErrorResponseWriter(w, err, http.StatusInternalServerError)
		return
	}

	pkg.JSONResponseWriter[Tokens](w, toTokens(tokens), http.StatusOK)
}
