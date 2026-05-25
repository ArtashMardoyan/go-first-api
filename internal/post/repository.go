package post

import (
	"context"

	"go-first-api/internal/shared"
)

type Repository interface {
	FindAll(ctx context.Context, q shared.PaginationQuery) ([]Post, int64, error)
	FindByID(ctx context.Context, id string) (Post, error)
	FindByUserID(ctx context.Context, userID string, q shared.PaginationQuery) ([]Post, int64, error)
	Create(ctx context.Context, p *Post) error
	Save(ctx context.Context, p *Post) error
	Delete(ctx context.Context, id string) error
}