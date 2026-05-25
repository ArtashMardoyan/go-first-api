package post

import (
	"go-first-api/internal/shared"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	shared.Base
	ID     string `json:"id"     gorm:"primaryKey"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Status Status `json:"status" gorm:"default:unpublished"`
	UserID string `json:"userId" gorm:"column:userId;index"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	return nil
}