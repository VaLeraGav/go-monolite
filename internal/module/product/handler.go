package product

import (
	"go-monolite/internal/store"

	"github.com/go-chi/chi"
)

var (
	MessNotFound    = "Товар не найден"
	MessGetProduct  = "Произошла ошибка при получении товара"
	MessInvalidJSON = "Получен некорректный формат JSON"
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
	// r.Get("/", h.GetList)
	// r.Get("/{uuid}", h.GetByUUID)
	// r.Post("/create", h.Create)
	// r.Put("/update/{uuid}", h.Update)
	// r.Delete("/delete/{uuid}", h.Delete)
}

// func (h *Handler) GetByUUID(w http.ResponseWriter, r *http.Request) {
// 	uuidStr := chi.URLParam(r, "uuid")

// 	err := validator.ParseUUID(uuidStr)
// 	if err != nil {
// 		respond.ErrorHandler(w, r, http.StatusBadRequest, err)
// 		return
// 	}

// 	product, err := h.service.GetByUUID(r.Context(), uuidStr)
// 	if err != nil {
// 		logger.ErrorCtx(r.Context(), err, MessGetProduct)
// 		respond.ErrorHandler(w, r, http.StatusInternalServerError, err, MessGetProduct)
// 		return
// 	}
// 	if product == nil {
// 		respond.ErrorHandler(w, r, http.StatusNotFound, nil, MessNotFound)
// 		return
// 	}

// 	respond.SuccessHandler(w, r, http.StatusOK, "", product)
// }

// // @Summary Create new product
// // @Description Create a new product with the provided details
// // @Tags products
// // @Accept json
// // @Produce json
// // @Param product body ProductDto true "Product object"
// // @Success 200 {array} ProductDto
// // @Failure 400 {array} respond.ErrorResponse
// // @Failure 500 {array} respond.ErrorResponse
// // @Router /product [post]
// func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
// 	body := respond.ParseBody(w, r)

// 	var request ProductDto

// 	if err := json.Unmarshal(body, &request); err != nil {
// 		respond.ErrorHandler(w, r, http.StatusBadRequest, err, MessInvalidJSON)
// 		return
// 	}

// 	product, err := h.service.GetByUUID(r.Context(), request.UUID.String())
// 	if err != nil {
// 		logger.ErrorCtx(r.Context(), err, "Произошла ошибка при проверке существования товара")
// 		respond.ErrorHandler(w, r, http.StatusInternalServerError, err, "Произошла ошибка при проверке существования товара")
// 		return
// 	}

// 	if product != nil {
// 		respond.ErrorHandler(w, r, http.StatusConflict, nil, "Товар с таким UUID уже существует")
// 		return
// 	}

// 	resultId, err := h.service.Create(r.Context(), &request)
// 	if err != nil || resultId == nil {
// 		if validationErrors, ok := err.(validator.ValidationError); ok {
// 			respond.ErrorHandler(w, r, http.StatusBadRequest, validationErrors.Fields, validator.ErrorValidation)
// 			return
// 		}
// 		logger.ErrorCtx(r.Context(), err, "Произошла ошибка при создании товара")
// 		respond.ErrorHandler(w, r, http.StatusBadRequest, err, "Произошла ошибка при создании товара")
// 		return
// 	}

// 	respond.SuccessHandler(w, r, http.StatusCreated, "", map[string]uint{"Id": *resultId})
// }

// func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
// 	uuidStr := chi.URLParam(r, "uuid")

// 	err := validator.ParseUUID(uuidStr)
// 	if err != nil {
// 		respond.ErrorHandler(w, r, http.StatusBadRequest, err)
// 		return
// 	}

// 	existingProduct, err := h.service.GetByUUID(r.Context(), uuidStr)
// 	if err != nil {
// 		logger.ErrorCtx(r.Context(), err, MessGetProduct)
// 		respond.ErrorHandler(w, r, http.StatusInternalServerError, err, MessGetProduct)
// 		return
// 	}
// 	if existingProduct == nil {
// 		respond.ErrorHandler(w, r, http.StatusNotFound, nil, MessNotFound)
// 		return
// 	}

// 	body := respond.ParseBody(w, r)

// 	var request ProductDto
// 	if err := json.Unmarshal(body, &request); err != nil {
// 		respond.ErrorHandler(w, r, http.StatusBadRequest, err, MessInvalidJSON)
// 		return
// 	}

// 	request = existingProduct.PatchDto(request)

// 	err = h.service.Update(r.Context(), &request)
// 	if err != nil {
// 		if errors.Is(err, store.ErrNotFound) {
// 			respond.ErrorHandler(w, r, http.StatusNotFound, nil, MessNotFound)
// 			return
// 		}
// 		if validationErrors, ok := err.(validator.ValidationError); ok {
// 			respond.ErrorHandler(w, r, http.StatusBadRequest, validator.ErrorValidation, validationErrors.Fields)
// 			return
// 		}
// 		logger.ErrorCtx(r.Context(), err, "Произошла ошибка при обновлении товара")
// 		respond.ErrorHandler(w, r, http.StatusInternalServerError, err, "Произошла ошибка при обновлении товара")
// 		return
// 	}

// 	respond.SuccessHandler(w, r, http.StatusOK, "Товар успешно обновлен", nil)
// }

// func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
// 	uuidStr := chi.URLParam(r, "uuid")

// 	err := validator.ParseUUID(uuidStr)
// 	if err != nil {
// 		respond.ErrorHandler(w, r, http.StatusBadRequest, err)
// 		return
// 	}

// 	err = h.service.Delete(r.Context(), uuidStr)
// 	if err != nil {
// 		if errors.Is(err, store.ErrNotFound) {
// 			respond.ErrorHandler(w, r, http.StatusNotFound, err, MessNotFound)
// 			return
// 		}
// 		logger.ErrorCtx(r.Context(), err, "Произошла ошибка при удалении товара")
// 		respond.ErrorHandler(w, r, http.StatusInternalServerError, err, "Произошла ошибка при удалении товара")
// 		return
// 	}

// 	respond.SuccessHandler(w, r, http.StatusOK, "Товар успешно удален", nil)
// }

// func (h *Handler) GetList(w http.ResponseWriter, r *http.Request) {
// 	products, err := h.service.GetList(r.Context())
// 	if err != nil {
// 		logger.ErrorCtx(r.Context(), err, "Произошла ошибка при получении списка товаров")
// 		respond.ErrorHandler(w, r, http.StatusInternalServerError, err, "Произошла ошибка при получении списка товаров")
// 		return
// 	}
// 	if products == nil {
// 		respond.ErrorHandler(w, r, http.StatusNotFound, err, MessNotFound)
// 		return
// 	}
// 	// TODO: перевести а Response
// 	// productResponses := ToProductResponseList(products)
// 	respond.SuccessHandler(w, r, http.StatusOK, "", products)
// }
