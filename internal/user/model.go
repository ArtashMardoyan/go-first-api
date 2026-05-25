package user

type User struct {
	ID    string `json:"id"    gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email" gorm:"uniqueIndex"`
	Age   int    `json:"age"`
}

type CreateUserDto struct {
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age"   binding:"required,min=1"`
}

type UpdateUserDto struct {
	Name  string `json:"name"`
	Email string `json:"email"  binding:"omitempty,email"`
	Age   int    `json:"age"    binding:"omitempty,min=1"`
}