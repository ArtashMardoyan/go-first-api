package post

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("post not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll() []Post {
	var posts []Post
	r.db.Find(&posts)
	return posts
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

func (r *Repository) FindByUserID(userID string) []Post {
	var posts []Post
	r.db.Where(`"userId" = ?`, userID).Find(&posts)
	return posts
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