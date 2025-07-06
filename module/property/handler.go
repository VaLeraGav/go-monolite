package property

import (
	"go-monolite/internal/store"
	"go-monolite/pkg/logger"
	"go-monolite/pkg/respond"
	"go-monolite/pkg/validator"
	"net/http"

	"github.com/go-chi/chi"
)

var (
	MessInvalidJSON = "Получен некорректный формат JSON"
)

type Handler struct {
	service *Service
}

func NewHandler(store *store.Store) *Handler {
	propertyRepo := NewPropertyRepository(store)
	propertyValuesRepo := NewPropertyValuesRepository(store)
	service := NewService(propertyRepo, propertyValuesRepo)
	return &Handler{service: service}
}

func (h *Handler) Init(r chi.Router) {
	r.Post("/upsert", h.Upsert)
}

// @Summary Upsert property list
// @Description Create or update a list of properties
// @Tags properties
// @Accept json
// @Produce json
// @Param properties body []PropertyDto true "Array of property objects"
// @Success 201 {object} respond.SuccessResponse{data=[]PropertyResponse}
// @Failure 400 {object} respond.ErrorResponse
// @Failure 500 {object} respond.ErrorResponse
// @Router /upsert [post]
func (h *Handler) Upsert(w http.ResponseWriter, r *http.Request) {
	body := respond.ParseBody(w, r)

	propertyList, err := ParsePropertyDto(body)
	if err != nil {
		respond.ErrorHandler(w, r, http.StatusBadRequest, err, MessInvalidJSON)
		return
	}

	resp, err := h.service.Upsert(r.Context(), propertyList)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationError); ok {
			respond.ErrorHandler(w, r, http.StatusBadRequest, validationErrors.Fields, validator.ErrorValidation)
			return
		}
		logger.ErrorCtx(r.Context(), err, "Произошла ошибка при создании свойства")
		respond.ErrorHandler(w, r, http.StatusBadRequest, err, "Произошла ошибка при создании свойства")
		return
	}

	respond.SuccessHandler(w, r, http.StatusCreated, "", resp)
}
