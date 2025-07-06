package user

import (
	"go-monolite/internal/store"
	"go-monolite/pkg/respond"
	"net/http"

	"github.com/go-chi/chi"
)

type Handler struct {
	service *Service
}

func NewHandler(store *store.Store) *Handler {
	repo := NewRepository(store)
	service := NewService(repo)
	return &Handler{service: service}
}

func (h *Handler) Init(r chi.Router) {
	r.Get("/me", h.Me)
	// r.Put("/update/{uuid}", h.Update)
	// r.Get("/{id}", h.GetByUUID)        // для админа
	// r.Delete("/delete/{id}", h.Delete) // для админа
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	respond.SuccessHandler(w, r, http.StatusCreated, "", "")
}
