package storage

import (
	"go-monolite/internal/store"
	"go-monolite/pkg/helper"
	"go-monolite/pkg/logger"
	"go-monolite/pkg/respond"
	"go-monolite/pkg/validator"
	"net/http"

	"github.com/go-chi/chi"
)

type Handler struct {
	service *Service
}

func NewHandler(store *store.Store) *Handler {
	storageRepo := NewStorageRepository(store)
	productStoragesRepo := NewProductStoragesRepository(store)
	service := NewService(storageRepo, productStoragesRepo)
	return &Handler{service: service}
}

func (h *Handler) Init(r chi.Router) {
	r.Post("/upsert", h.Upsert)
	r.Get("/storages", h.GetStorage)
}

// @Summary Upsert storages
// @Description Create or update storage information
// @Tags storages
// @Accept json
// @Produce json
// @Param storage body StorageSyncRequest true "Storage sync request payload"
// @Success 201 {object} respond.SuccessResponse{data=StorageSyncResponse}
// @Failure 400 {object} respond.ErrorResponse
// @Failure 500 {object} respond.ErrorResponse
// @Router /upsert [post]
func (h *Handler) Upsert(w http.ResponseWriter, r *http.Request) {
	body := respond.ParseBody(w, r)

	var storageRequest UpsertRequest
	mess, err := helper.Unmarshal(body, &storageRequest)
	if err != nil {
		respond.ErrorHandler(w, r, http.StatusInternalServerError, err, mess)
		return
	}

	resp, mess, err := h.service.Upsert(r.Context(), storageRequest)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationError); ok {
			respond.ErrorHandler(w, r, http.StatusBadRequest, validationErrors.Fields, validator.ErrorValidation)
			return
		}
		logger.ErrorCtx(r.Context(), err, mess)
		respond.ErrorHandler(w, r, http.StatusBadRequest, err, mess)
		return
	}

	respond.SuccessHandler(w, r, http.StatusCreated, "", resp)
}

// @Summary Get storages
// @Description Get list of storages based on query parameters
// @Tags storages
// @Accept json
// @Produce json
// @Param name query string false "Storage name"
// @Param type query string false "Storage type"
// @Success 200 {object} respond.SuccessResponse{data=[]StorageResponse}
// @Failure 400 {object} respond.ErrorResponse
// @Failure 500 {object} respond.ErrorResponse
// @Router /storages [get]
func (h *Handler) GetStorage(w http.ResponseWriter, r *http.Request) {
	resp, mess, err := h.service.GetStorage(r.Context())
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationError); ok {
			respond.ErrorHandler(w, r, http.StatusBadRequest, validationErrors.Fields, validator.ErrorValidation)
			return
		}
		logger.ErrorCtx(r.Context(), err, mess)
		respond.ErrorHandler(w, r, http.StatusBadRequest, err, mess)
		return
	}

	respond.SuccessHandler(w, r, http.StatusCreated, "", resp)
}
