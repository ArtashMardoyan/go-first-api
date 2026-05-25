package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"go-first-api/internal/user"
)

type LoginDto struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string    `json:"accessToken"`
	User        user.User `json:"user"`
}

type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}