package auth

import "time"

type AuthCode struct {
	ID        int       `db:"id"`
	Email     string    `db:"email"`
	Phone     string    `db:"phone"`
	Code      string    `db:"code"`
	ExpiresAt time.Time `db:"expires_at"`
	Used      bool      `db:"used"`
	CreatedAt time.Time `db:"created_at"`
}