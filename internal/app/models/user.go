package models

import "github.com/go-playground/validator/v10"

type User struct {
	ID       int64  `json:"-"`
	Username string `json:"username" validate:"required,min=5,max=20"`
	Password string `json:"password" validate:"required,min=8,max=16"`
	Email    string `json:"email" validate:"required,email"`
}

func (u *User) Validate() error {
	return validate.Struct(u)
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type UserIdResponse struct {
	UserId int `json:"id"`
}

type SignInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
