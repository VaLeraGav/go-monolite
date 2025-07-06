package price

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
	typePriceRepo := NewTypePriceRepository(store)
	productPriceRepo := NewProductPricesRepository(store)
	service := NewService(typePriceRepo, productPriceRepo)
	return &Handler{service: service}
}

func (h *Handler) Init(r chi.Router) {
	r.Post("/upsert", h.Upsert)
	r.Get("/type-price", h.GetTypePrice)
}

// @Summary Upsert price information
// @Description Insert or update price data
// @Tags prices
// @Accept json
// @Produce json
// @Param price body PriceRequest true "Price data to upsert"
// @Success 201 {object} respond.SuccessResponse{data=PriceResponse}
// @Failure 400 {object} respond.ErrorResponse
// @Failure 500 {object} respond.ErrorResponse
// @Router /upsert [post]
func (h *Handler) Upsert(w http.ResponseWriter, r *http.Request) {
	body := respond.ParseBody(w, r)

	var priceRequest UpsertRequest
	mess, err := helper.Unmarshal(body, &priceRequest)
	if err != nil {
		respond.ErrorHandler(w, r, http.StatusInternalServerError, err, mess)
		return
	}

	resp, mess, err := h.service.Upsert(r.Context(), priceRequest)
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

// @Summary Get price types
// @Description Get list of price types for provided request
// @Tags prices
// @Accept json
// @Produce json
// @Param name query string false "Price name"
// @Success 201 {object} respond.SuccessResponse{data=[]PriceTypeResponse}
// @Failure 400 {object} respond.ErrorResponse
// @Failure 500 {object} respond.ErrorResponse
// @Router /type-price [get]
func (h *Handler) GetTypePrice(w http.ResponseWriter, r *http.Request) {
	resp, mess, err := h.service.GetTypePrice(r.Context())
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
