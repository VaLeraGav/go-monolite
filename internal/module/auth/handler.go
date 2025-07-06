package auth

import (
	"go-monolite/internal/module/user"
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
	userTokensRepo := NewUserTokensRepository(store)
	codesRepo := NewAuthCodeRepository(store)
	userRepo := user.NewRepository(store)
	service := NewService(userTokensRepo, codesRepo, userRepo)
	return &Handler{service: service}
}

func (h *Handler) Init(r chi.Router) {
	r.Post("/sendCode", h.SendCode)
	// r.Post("/login", h.Login)
	// r.Post("/againSendCode", h.AgainSendCode)

	// r.Post("/logout", h.Logout)
	// r.Post("/refresh", h.Refresh)
}

// принимает email/phone, отправляет код.
func (h *Handler) SendCode(w http.ResponseWriter, r *http.Request) {
	body := respond.ParseBody(w, r)

	var sendCodeRequest SendCodeRequest

	mess, err := helper.Unmarshal(body, &sendCodeRequest)
	if err != nil {
		respond.ErrorHandler(w, r, http.StatusInternalServerError, err, mess)
		return
	}

	err = h.service.SendCode(r.Context(), sendCodeRequest)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationError); ok {
			respond.ErrorHandler(w, r, http.StatusBadRequest, validationErrors.Fields, validator.ErrorValidation)
			return
		}
		logger.ErrorCtx(r.Context(), err, mess)
		respond.ErrorHandler(w, r, http.StatusInternalServerError, err, mess)
		return
	}

	respond.SuccessHandler(w, r, http.StatusCreated, "", "")
}

// принимает email/phone + код, создаёт пользователя.
// func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
// 	body := respond.ParseBody(w, r)

// 	var registrationRequests RegistrationRequest

// 	mess, err := helper.Unmarshal(body, &registrationRequests)
// 	if err != nil {
// 		respond.ErrorHandler(w, r, http.StatusInternalServerError, err, mess)
// 		return
// 	}

// 	registrationResponse, mess, err := h.service.Register(r.Context(), registrationRequests)
// 	if err != nil {
// 		if validationErrors, ok := err.(validator.ValidationError); ok {
// 			respond.ErrorHandler(w, r, http.StatusBadRequest, validationErrors.Fields, validator.ErrorValidation)
// 			return
// 		}
// 		logger.ErrorCtx(r.Context(), err, mess)
// 		respond.ErrorHandler(w, r, http.StatusInternalServerError, err, mess)
// 		return
// 	}

// 	respond.SuccessHandler(w, r, http.StatusCreated, mess, registrationResponse)
// }

// func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
// 	body := respond.ParseBody(w, r)

// 	var categoryRequests LoginRequest

//		mess, err := helper.Unmarshal(body, &categoryRequests)
//		if err != nil {
//			respond.ErrorHandler(w, r, http.StatusInternalServerError, err, mess)
//			return
//		}
//		respond.SuccessHandler(w, r, http.StatusCreated, mess, "")
//	}
//
//	func (h *Handler) AgainSendCode(w http.ResponseWriter, r *http.Request) {
//		return
//	}
//
//	func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
//		return
//	}
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	// // 1. Извлечь refresh токен (например, из cookie)
	// cookie, err := r.Cookie("refresh_token")
	// if err != nil {
	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	// refreshToken := cookie.Value

	// // 2. Распарсить и проверить токен
	// claims, err := token.ParseRefreshToken(refreshToken)
	// if err != nil {
	// 	http.Error(w, "invalid refresh token", http.StatusUnauthorized)
	// 	return
	// }

	// // 3. Найти сессию в БД по user_id и device_id
	// session, err := h.sessionRepo.FindByUserAndDevice(ctx, claims.UserID, claims.DeviceID)
	// if err != nil || session.RefreshToken != refreshToken {
	// 	http.Error(w, "invalid session", http.StatusUnauthorized)
	// 	return
	// }

	// // 4. Проверка срока жизни и token_version
	// if time.Now().After(session.ExpiresAt) || claims.TokenVersion != session.TokenVersion {
	// 	http.Error(w, "token expired or invalid", http.StatusUnauthorized)
	// 	return
	// }

	// // 5. Создать новый access + refresh токены
	// newAccessToken, err := token.GenerateAccessToken(claims.UserID)
	// if err != nil {
	// 	http.Error(w, "failed to create access token", http.StatusInternalServerError)
	// 	return
	// }

	// newRefreshToken, err := token.GenerateRefreshToken(claims.UserID, claims.DeviceID, claims.TokenVersion)
	// if err != nil {
	// 	http.Error(w, "failed to create refresh token", http.StatusInternalServerError)
	// 	return
	// }

	// // 6. Обновить refresh токен в БД (если хранишь как opaque/hard-match)
	// err = h.sessionRepo.UpdateRefreshToken(ctx, claims.UserID, claims.DeviceID, newRefreshToken)
	// if err != nil {
	// 	http.Error(w, "failed to update session", http.StatusInternalServerError)
	// 	return
	// }

	// // 7. Установить новый refresh в cookie
	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "refresh_token",
	// 	Value:    newRefreshToken,
	// 	HttpOnly: true,
	// 	Path:     "/",
	// 	MaxAge:   60 * 60 * 24 * 30, // 30 дней
	// })

	// // 8. Вернуть access токен
	// resp := map[string]interface{}{
	// 	"access_token": newAccessToken,
	// 	"expires_in":   3600,
	// }
	// respond.SuccessHandler(w, r, http.StatusCreated, "", resp)
}
