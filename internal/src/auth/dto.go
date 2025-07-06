package auth

type SendCodeRequest struct {
	Email string `json:"email,omitempty"` // обязательно email или phone
	Phone string `json:"phone,omitempty"`
	Type  string `json:"type,omitempty"`
}

// ------------------------------------------
// RegistrationRequest — DTO для POST /auth/register
type RegistrationRequest struct {
	Email      *string `json:"email,omitempty"` // обязательно email или phone
	Phone      *string `json:"phone,omitempty"`
	Password   string  `json:"password" validate:"required,min=6"`
	Name       *string `json:"name,omitempty"`
	LastName   *string `json:"last_name,omitempty"`
	SecondName *string `json:"second_name,omitempty"`
	CityID     *int64  `json:"city_id,omitempty"`
	UserType   string  `json:"user_type" validate:"required,oneof=individual legal"`
	INN        *string `json:"inn,omitempty"`
}

// LoginRequest — DTO для POST /auth/login
type LoginRequest struct {
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Password string  `json:"password" validate:"required"`
}

// AuthResponse — Ответ на успешный логин/регистрацию
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // в секундах
}
