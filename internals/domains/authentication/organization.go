package authentication

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID                  uuid.UUID
	Email               string
	Password            string
	FirstName           string
	LastName            string
	FullName            string
	EmailVerified       bool
	RefreshTokenVersion int
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type AddUserParams struct {
	Email     string `json:"email"    binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

type VerifyCodeParams struct {
	Email string `json:"email"    binding:"required,email"`
	Code  string `json:"code"     binding:"required"`
}

type ResendCodeParams struct {
	Email string `json:"email"    binding:"required,email"`
}

type GetUserParams struct {
	Email string    `json:"email"    binding:"omitempty,email"`
	ID    uuid.UUID `json:"id"       binding:"omitempty,uuid"`
}

type LoginParams struct {
	Email    string `json:"email"        binding:"required,email"`
	Password string `json:"password"     binding:"required"`
}

type UserProfileParams struct {
	ID            uuid.UUID
	FirstName     string `json:"firstName" binding:"required"`
	LastName      string `json:"lastName"  binding:"required"`
	EmailVerified bool
}
