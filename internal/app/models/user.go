package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type UserInput struct {
	ID       int64  `json:"-"`
	Username string `json:"username" validate:"required,min=5,max=20"`
	Password string `json:"password" validate:"required,min=8,max=16"`
	Email    string `json:"email" validate:"required,email"`
}

type GetUserResponse struct {
	ID       int64     `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	CreateAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
}

func (u *UserInput) Validate() error {
	return validate.Struct(u)
}

type UserIdResponse struct {
	UserId int `json:"id"`
}

type SignInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
