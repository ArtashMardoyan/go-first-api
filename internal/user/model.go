package user

type User struct {
	ID       string `json:"id"       gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email"    gorm:"uniqueIndex"`
	Age      int    `json:"age"`
	Password string `json:"-"`
}

type CreateUserDto struct {
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Age      int    `json:"age"      binding:"required,min=1"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserDto struct {
	Name string `json:"name"`
	Age  int    `json:"age" binding:"omitempty,min=1"`
}