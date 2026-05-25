package user

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ErrNotFound — sentinel error, аналог NotFoundException в NestJS
var ErrNotFound = errors.New("user not found")

// ErrEmailTaken — возвращается когда email уже занят другим пользователем
var ErrEmailTaken = errors.New("email already taken")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll() []User {
	var users []User
	r.db.Find(&users)
	return users
}

func (r *Repository) FindByID(id string) (User, error) {
	var u User
	if err := r.db.First(&u, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}
	return u, nil
}

func (r *Repository) Create(dto CreateUserDto) (User, error) {
	var count int64
	r.db.Model(&User{}).Where("email = ?", dto.Email).Count(&count)
	if count > 0 {
		return User{}, ErrEmailTaken
	}

	u := User{
		ID:    uuid.NewString(),
		Name:  dto.Name,
		Email: dto.Email,
		Age:   dto.Age,
	}
	if err := r.db.Create(&u).Error; err != nil {
		return User{}, err
	}
	return u, nil
}

func (r *Repository) Update(id string, dto UpdateUserDto) (User, error) {
	var u User
	if err := r.db.First(&u, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}

	if dto.Email != "" && dto.Email != u.Email {
		var count int64
		r.db.Model(&User{}).Where("email = ?", dto.Email).Count(&count)
		if count > 0 {
			return User{}, ErrEmailTaken
		}
		u.Email = dto.Email
	}
	if dto.Name != "" {
		u.Name = dto.Name
	}
	if dto.Age != 0 {
		u.Age = dto.Age
	}

	if err := r.db.Save(&u).Error; err != nil {
		return User{}, err
	}
	return u, nil
}

func (r *Repository) Delete(id string) error {
	result := r.db.Where("id = ?", id).Delete(&User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}