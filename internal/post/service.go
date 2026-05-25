package post

import (
	"context"
	"errors"

	"go-first-api/internal/shared"

	"gorm.io/gorm"
)

var (
	ErrNotFound      = errors.New("post not found")
	ErrInvalidStatus = errors.New("invalid status")
	ErrForbidden     = errors.New("forbidden")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) FindAll(ctx context.Context, q shared.PaginationQuery) (shared.PaginatedResult[Post], error) {
	q.Normalize()
	posts, total, err := s.repo.FindAll(ctx, q)
	if err != nil {
		return shared.PaginatedResult[Post]{}, err
	}
	return shared.NewPaginatedResult(posts, total, q), nil
}

func (s *Service) FindByID(ctx context.Context, id string) (Post, error) {
	p, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Post{}, ErrNotFound
	}
	return p, err
}

func (s *Service) FindByUserID(ctx context.Context, userID string, q shared.PaginationQuery) (shared.PaginatedResult[Post], error) {
	q.Normalize()
	posts, total, err := s.repo.FindByUserID(ctx, userID, q)
	if err != nil {
		return shared.PaginatedResult[Post]{}, err
	}
	return shared.NewPaginatedResult(posts, total, q), nil
}

func (s *Service) Create(ctx context.Context, userID string, dto CreateDTO) (Post, error) {
	p := Post{
		Title:  dto.Title,
		Body:   dto.Body,
		UserID: userID,
		Status: StatusUnpublished,
	}
	if err := s.repo.Create(ctx, &p); err != nil {
		return Post{}, err
	}
	return p, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id, callerID string, dto UpdateStatusDTO) (Post, error) {
	if !dto.Status.IsValid() {
		return Post{}, ErrInvalidStatus
	}
	p, err := s.FindByID(ctx, id)
	if err != nil {
		return Post{}, err
	}
	if p.UserID != callerID {
		return Post{}, ErrForbidden
	}
	p.Status = dto.Status
	if err := s.repo.Save(ctx, &p); err != nil {
		return Post{}, err
	}
	return p, nil
}

func (s *Service) Delete(ctx context.Context, id, callerID string) error {
	p, err := s.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if p.UserID != callerID {
		return ErrForbidden
	}
	return s.repo.Delete(ctx, id)
}
