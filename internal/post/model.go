package post

type Post struct {
	ID     string `json:"id"     gorm:"primaryKey"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID string `json:"userId" gorm:"column:userId;index"`
}

type CreatePostDto struct {
	Title  string `json:"title"  binding:"required"`
	Body   string `json:"body"   binding:"required"`
	UserID string `json:"userId" binding:"required"`
}