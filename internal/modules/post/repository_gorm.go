package post

import (
	"context"

	"go-first-api/internal/shared"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) FindAll(ctx context.Context, q shared.PaginationQuery) ([]Post, int64, error) {
	var posts []Post
	var total int64

	if err := r.db.WithContext(ctx).Model(&Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Preload("User").Offset(q.Offset()).Limit(q.Limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *gormRepository) FindByID(ctx context.Context, id string) (Post, error) {
	var p Post

	err := r.db.WithContext(ctx).Preload("User").First(&p, "id = ?", id).Error

	return p, err
}

func (r *gormRepository) FindByUserID(ctx context.Context, userID string, q shared.PaginationQuery) ([]Post, int64, error) {
	var posts []Post
	var total int64

	base := r.db.WithContext(ctx).Model(&Post{}).Where(`"userId" = ?`, userID)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := base.Preload("User").Offset(q.Offset()).Limit(q.Limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *gormRepository) Create(ctx context.Context, p *Post) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *gormRepository) Save(ctx context.Context, p *Post) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *gormRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&Post{}, "id = ?", id).Error
}
