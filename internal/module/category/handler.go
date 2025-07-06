package category

import (
	"errors"
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
	repo := NewRepository(store)
	service := NewService(repo)
	return &Handler{service: service}
}

func (h *Handler) Init(r chi.Router) {
	r.Get("/{uuid}", h.GetByUUID)
	r.Post("/create", h.Create)
	r.Put("/update/{uuid}", h.Update)
	r.Delete("/delete/{uuid}", h.Delete)
	r.Get("/tree", h.GetTree)
	r.Get("/tree/{uuid}", h.GetTree)
}

// @Summary Get category by UUID
// @Description Get a single category by its UUID
// @Tags categories
// @Accept json
// @Produce json
// @Param uuid path string true "Category UUID"
// @Success 200 {object} respond.SuccessResponse{data=CategoryResponse}
// @Failure 400 {object} respond.ErrorResponse
// @Failure 404 {object} respond.ErrorResponse
// @Failure 500 {object} respond.ErrorResponse
// @Router /{uuid} [get]
func (h *Handler) GetByUUID(w http.ResponseWriter, r *http.Request) {
	uuidStr := chi.URLParam(r, "uuid")

	_, err := validator.ParseUUID(uuidStr)
	if err != nil {
		respond.ErrorHandler(w, r, http.StatusInternalServerError, err)
		return
	}

	category, mess, err := h.service.GetByUUID(r.Context(), uuidStr)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			respond.ErrorHandler(w, r, http.StatusNotFound, mess, mess)
			return
		}
		logger.ErrorCtx(r.Context(), err, mess)
		respond.ErrorHandler(w, r, http.StatusInternalServerError, err, mess)
		return
	}
	if category == nil {
		respond.ErrorHandler(w, r, http.StatusInternalServerError, nil, mess)
		return
	}

	respond.SuccessHandler(w, r, http.StatusOK, "", category)
}

// @Summary Create new categories
// @Description Create one or more new categories
// @Tags categories
// @Accept json
// @Produce json
// @Param categories body []CategoryDto true "Array of category objects"
// @Success 201 {array} respond.SuccessResponse{data=[]CategoryResponse}
// @Failure 400 {object} respond.ErrorResponse
// @Failure 500 {object} respond.ErrorResponse
// @Router /create [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	body := respond.ParseBody(w, r)

	var categoryRequests []CategoryRequest

	mess, err := helper.Unmarshal(body, &categoryRequests)
	if err != nil {
		respond.ErrorHandler(w, r, http.StatusInternalServerError, err, mess)
		return
	}

	categoryResponse, mess, err := h.service.Create(r.Context(), categoryRequests)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationError); ok {
			respond.ErrorHandler(w, r, http.StatusBadRequest, validationErrors.Fields, validator.ErrorValidation)
			return
		}
		logger.ErrorCtx(r.Context(), err, mess)
		respond.ErrorHandler(w, r, http.StatusInternalServerError, err, mess)
		return
	}

	respond.SuccessHandler(w, r, http.StatusCreated, mess, categoryResponse)
}

// @Summary Update category
// @Description Update an existing category by UUID
// @Tags categories
// @Accept json
// @Produce json
// @Param uuid path string true "Category UUID"
// @Param category body CategoryDto true "Updated category object"
// @Success 200 {object} respond.SuccessResponse{data=CategoryResponse}
// @Failure 400 {object} respond.ErrorResponse
// @Failure 404 {object} respond.ErrorResponse
// @Failure 500 {object} respond.ErrorResponse
// @Router /update/{uuid} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	uuidStr := chi.URLParam(r, "uuid")

	uuid, err := validator.ParseUUID(uuidStr)
	if err != nil {
		respond.ErrorHandler(w, r, http.StatusInternalServerError, err)
		return
	}

	body := respond.ParseBody(w, r)

	var request CategoryRequest
	mess, err := helper.Unmarshal(body, &request)
	if err != nil {
		respond.ErrorHandler(w, r, http.StatusInternalServerError, err, mess)
		return
	}

	request.UUID = uuid

	categoryResponse, mess, err := h.service.Update(r.Context(), &request)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationError); ok {
			respond.ErrorHandler(w, r, http.StatusBadRequest, validationErrors.Fields, validator.ErrorValidation)
			return
		}
		if errors.Is(err, store.ErrNotFound) {
			respond.ErrorHandler(w, r, http.StatusNotFound, mess, mess)
			return
		}
		logger.ErrorCtx(r.Context(), err, mess)
		respond.ErrorHandler(w, r, http.StatusInternalServerError, err, mess)
		return
	}

	respond.SuccessHandler(w, r, http.StatusOK, mess, categoryResponse)
}

// @Summary Delete category
// @Description Delete a category by UUID
// @Tags categories
// @Accept json
// @Produce json
// @Param uuid path string true "Category UUID"
// @Success 200 {object} respond.SuccessResponse
// @Failure 404 {object} respond.ErrorResponse
// @Failure 500 {object} respond.ErrorResponse
// @Router /delete/{uuid} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	uuidStr := chi.URLParam(r, "uuid")

	_, err := validator.ParseUUID(uuidStr)
	if err != nil {
		respond.ErrorHandler(w, r, http.StatusInternalServerError, err)
		return
	}

	mess, err := h.service.Delete(r.Context(), uuidStr)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			respond.ErrorHandler(w, r, http.StatusNotFound, nil, mess)
			return
		}
		logger.ErrorCtx(r.Context(), err, mess)
		respond.ErrorHandler(w, r, http.StatusInternalServerError, err, mess)
		return
	}

	respond.SuccessHandler(w, r, http.StatusOK, mess)
}

// @Summary Get category tree
// @Description Get category tree (optionally from a specific UUID node)
// @Tags categories
// @Accept json
// @Produce json
// @Param uuid path string false "Category UUID (optional root)"
// @Success 200 {array} respond.SuccessResponse{data=[]CategoryTreeResponse}
// @Failure 500 {object} respond.ErrorResponse
// @Router /tree [get]
// @Router /tree/{uuid} [get]
func (h *Handler) GetTree(w http.ResponseWriter, r *http.Request) {
	uuidStr := chi.URLParam(r, "uuid")

	if uuidStr != "" {
		_, err := validator.ParseUUID(uuidStr)
		if err != nil {
			respond.ErrorHandler(w, r, http.StatusInternalServerError, err)
			return
		}
	}

	categories, mess, err := h.service.GetTree(r.Context(), uuidStr)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			respond.ErrorHandler(w, r, http.StatusNotFound, mess, mess)
			return
		}
		logger.ErrorCtx(r.Context(), err, mess)
		respond.ErrorHandler(w, r, http.StatusInternalServerError, err, mess)
		return
	}
	if categories == nil {
		respond.ErrorHandler(w, r, http.StatusInternalServerError, nil, mess)
		return
	}

	respond.SuccessHandler(w, r, http.StatusOK, "", categories)
}
