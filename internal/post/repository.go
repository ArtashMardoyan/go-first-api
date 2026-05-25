package post

import (
	"errors"

	"go-first-api/pkg/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("post not found")
var ErrInvalidStatus = errors.New("invalid status: use 'published' or 'unpublished'")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll(q pagination.Query) ([]Post, int64) {
	var posts []Post
	var total int64
	r.db.Model(&Post{}).Count(&total)
	r.db.Offset(q.Offset()).Limit(q.Limit).Find(&posts)
	return posts, total
}

func (r *Repository) FindByID(id string) (Post, error) {
	var p Post
	if err := r.db.First(&p, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Post{}, ErrNotFound
		}
		return Post{}, err
	}
	return p, nil
}

func (r *Repository) FindByUserID(userID string, q pagination.Query) ([]Post, int64) {
	var posts []Post
	var total int64
	r.db.Model(&Post{}).Where(`"userId" = ?`, userID).Count(&total)
	r.db.Where(`"userId" = ?`, userID).Offset(q.Offset()).Limit(q.Limit).Find(&posts)
	return posts, total
}

func (r *Repository) Create(dto CreatePostDto) (Post, error) {
	p := Post{
		ID:     uuid.NewString(),
		Title:  dto.Title,
		Body:   dto.Body,
		UserID: dto.UserID,
	}
	if err := r.db.Create(&p).Error; err != nil {
		return Post{}, err
	}
	return p, nil
}

func (r *Repository) UpdateStatus(id string, status Status) (Post, error) {
	var p Post
	if err := r.db.First(&p, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Post{}, ErrNotFound
		}
		return Post{}, err
	}
	p.Status = status
	if err := r.db.Save(&p).Error; err != nil {
		return Post{}, err
	}
	return p, nil
}

func (r *Repository) Delete(id string) error {
	result := r.db.Where("id = ?", id).Delete(&Post{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}