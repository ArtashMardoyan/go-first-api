package auth

import "github.com/golang-jwt/jwt/v5"

type LoginDto struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	AccessToken string `json:"accessToken"`
}

type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}