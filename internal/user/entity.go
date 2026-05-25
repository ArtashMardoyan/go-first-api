package user

import (
	"go-first-api/internal/shared"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	shared.Base
	ID       string `json:"id"       gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email"    gorm:"uniqueIndex"`
	Age      int    `json:"age"`
	Status   Status `json:"status"   gorm:"default:activated"`
	Password string `json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	return nil
}