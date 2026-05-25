package middleware

import (
	"errors"
	"net/http"
	"strings"

	"go-first-api/internal/modules/user"
	"go-first-api/internal/shared"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func JWT(userRepo user.Repository, jwtSecret string) gin.HandlerFunc {
	secret := []byte(jwtSecret)
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			shared.Error(c, http.StatusUnauthorized, "missing or invalid authorization header")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims := &jwtClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			shared.Error(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		u, err := userRepo.FindByID(c.Request.Context(), claims.UserID)
		if err != nil {
			shared.Error(c, http.StatusUnauthorized, "user not found")
			c.Abort()
			return
		}

		c.Set(shared.ContextUserKey, u)
		c.Next()
	}
}
