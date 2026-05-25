package user

import (
	"errors"

	"go-first-api/pkg/pagination"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("user not found")
var ErrEmailTaken = errors.New("email already taken")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll(q pagination.Query) ([]User, int64) {
	var users []User
	var total int64
	r.db.Model(&User{}).Count(&total)
	r.db.Offset(q.Offset()).Limit(q.Limit).Find(&users)
	return users, total
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

func (r *Repository) FindByEmail(email string) (User, error) {
	var u User
	if err := r.db.First(&u, "email = ?", email).Error; err != nil {
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

	hashed, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	u := User{
		ID:       uuid.NewString(),
		Name:     dto.Name,
		Email:    dto.Email,
		Age:      dto.Age,
		Password: string(hashed),
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