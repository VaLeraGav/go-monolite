package user

// UserResponse — DTO для GET /users/me
type UserResponse struct {
	ID         uint    `json:"id"`
	Email      *string `json:"email,omitempty"`
	Phone      *string `json:"phone,omitempty"`
	Name       *string `json:"name,omitempty"`
	LastName   *string `json:"last_name,omitempty"`
	SecondName *string `json:"second_name,omitempty"`
	CityID     *int64  `json:"city_id,omitempty"`
	UserType   string  `json:"user_type"`
	INN        *string `json:"inn,omitempty"`
	Active     string  `json:"active"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

// UpdateUserRequest — DTO для PUT /users/me
type UpdateUserRequest struct {
	Email      string `json:"email,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Name       string `json:"name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	SecondName string `json:"second_name,omitempty"`
	CityID     int64  `json:"city_id,omitempty"`
	INN        string `json:"inn,omitempty"`
}
