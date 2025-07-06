package user

import (
	"time"
)

type UserType string

const (
	UserTypeIndividual UserType = "individual"
	UserTypeLegal      UserType = "legal"
)

type UserEnt struct {
	ID           int64     `db:"id"`
	Email        *string   `db:"email"`
	Phone        *string   `db:"phone"`
	Name         *string   `db:"name"`
	LastName     *string   `db:"last_name"`
	SecondName   *string   `db:"second_name"`
	CityID       *int64    `db:"city_id"`
	UserType     UserType  `db:"user_type"`
	INN          *string   `db:"inn"`
	Active       string    `db:"active"`
	PasswordHash string    `db:"password_hash"`
	Checkword    *string   `db:"checkword"`
	TokenVersion *string   `db:"token_version"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
