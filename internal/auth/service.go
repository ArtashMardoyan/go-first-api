package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"go-first-api/internal/user"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Service struct {
	userRepo user.Repository
}

func NewService(userRepo user.Repository) *Service {
	return &Service{userRepo: userRepo}
}

func (s *Service) Login(ctx context.Context, dto LoginDTO) (LoginResponseDTO, error) {
	u, err := s.userRepo.FindByEmail(ctx, dto.Email)
	if err != nil {
		return LoginResponseDTO{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(dto.Password)); err != nil {
		return LoginResponseDTO{}, ErrInvalidCredentials
	}

	if u.Status == user.StatusDeactivated {
		return LoginResponseDTO{}, errors.New("account is deactivated")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return LoginResponseDTO{}, err
	}

	return LoginResponseDTO{AccessToken: signed, User: u}, nil
}