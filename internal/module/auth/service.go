package auth

import (
	"context"
	"errors"
	"fmt"
	"go-monolite/internal/infra/sender"
	"go-monolite/internal/module/user"
	"go-monolite/internal/store"
	"go-monolite/pkg/helper"
	"time"
)

const maxDailyAttempts = 5

type Service struct {
	userTokensRepo *UserTokensRepository
	codeRepo       *AuthCodeRepository
	userRepo       *user.Repository
	sender         *sender.Sender
}

func NewService(userTokensRepo *UserTokensRepository, codeRepo *AuthCodeRepository, userRepo *user.Repository) *Service {
	sender := sender.New()
	return &Service{
		userTokensRepo: userTokensRepo,
		userRepo:       userRepo,
		codeRepo:       codeRepo,
		sender:         sender,
	}
}

func (s *Service) SendCode(ctx context.Context, req SendCodeRequest) error {

	//  валидация

	// 2. Проверка лимита попыток
	count, err := s.codeRepo.CountCodesLast24Hours(ctx, req.Email, req.Phone)
	if err != nil {
		return fmt.Errorf("failed to count sent codes: %w", err)
	}
	if count >= maxDailyAttempts {
		return fmt.Errorf("maximum number of attempts reached, try again later")
	}

	// 3. Повторно отправляем или создаём новый код
	existingCode, err := s.codeRepo.GetActiveCode(ctx, req.Email, req.Phone)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return fmt.Errorf("get active code error: %w", err)
	}

	var code string
	if existingCode != nil {
		code = existingCode.Code
	} else {
		code = helper.GenerateCode(4)
		expiry := time.Now().Add(5 * time.Minute)
		if err := s.codeRepo.SaveCode(ctx, req.Email, req.Phone, code, expiry); err != nil {
			return fmt.Errorf("failed to save code: %w", err)
		}
	}

	// 4. Отправка
	if req.Email != "" {
		err = s.sender.SendEmailCode(req.Email, code)
	} else {
		err = s.sender.SendSmsCode(req.Phone, code)
	}
	if err != nil {
		return fmt.Errorf("failed to send code: %w", err)
	}

	return nil
}

// func (s *Service) Register(ctx context.Context, registrationResponse RegistrationRequest) (*AuthResponse, string, error) {

// // 1. Проверка: есть ли активный код с таким email/телефоном и purpose = 'register'
// codeValid, err := s.codeRepo.ValidateCode(ctx, req.Email, req.Phone, req.Code, "register")
// if err != nil {
// 	return nil, "", fmt.Errorf("failed to validate code: %w", err)
// }
// if !codeValid {
// 	return nil, "", fmt.Errorf("invalid or expired verification code")
// }

// // 2. Проверка, не существует ли уже такой пользователь
// exists, err := s.userRepo.ExistsByEmailOrPhone(ctx, req.Email, req.Phone)
// if err != nil {
// 	return nil, "", fmt.Errorf("failed to check user: %w", err)
// }
// if exists {
// 	return nil, "", fmt.Errorf("user already exists")
// }

// // 3. Создание пользователя
// user := &domain.User{
// 	Email:  req.Email,
// 	Phone:  req.Phone,
// 	Name:   req.Name,
// 	Active: "Y",
// }
// if err := s.userRepo.Create(ctx, user); err != nil {
// 	return nil, "", fmt.Errorf("failed to create user: %w", err)
// }

// // 4. Генерация токенов
// accessToken, err := s.tokenService.GenerateAccessToken(user.ID)
// if err != nil {
// 	return nil, "", fmt.Errorf("failed to create access token: %w", err)
// }
// refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, "", user.TokenVersion)
// if err != nil {
// 	return nil, "", fmt.Errorf("failed to create refresh token: %w", err)
// }

// return &AuthResponse{
// 	UserID:      user.ID,
// 	AccessToken: accessToken,
// 	ExpiresIn:   3600,
// }, refreshToken, nil

// }
