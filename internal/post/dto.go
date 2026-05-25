package post

type CreateDTO struct {
	Title string `json:"title" binding:"required,min=2"`
	Body  string `json:"body"  binding:"required"`
}

type UpdateStatusDTO struct {
	Status Status `json:"status" binding:"required"`
}