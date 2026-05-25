package pagination

type Query struct {
	Page  int `form:"page"  binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

func (q *Query) Normalize() {
	if q.Page == 0 {
		q.Page = 1
	}
	if q.Limit == 0 {
		q.Limit = 10
	}
}

func (q Query) Offset() int {
	return (q.Page - 1) * q.Limit
}

type Meta struct {
	Total           int64 `json:"total"`
	Page            int   `json:"page"`
	Limit           int   `json:"limit"`
	PageCount       int   `json:"pageCount"`
	HasNextPage     bool  `json:"hasNextPage"`
	HasPreviousPage bool  `json:"hasPreviousPage"`
}

func NewMeta(total int64, q Query) Meta {
	pageCount := int((total + int64(q.Limit) - 1) / int64(q.Limit))
	return Meta{
		Total:           total,
		Page:            q.Page,
		Limit:           q.Limit,
		PageCount:       pageCount,
		HasNextPage:     q.Page < pageCount,
		HasPreviousPage: q.Page > 1,
	}
}

type Result[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"_meta"`
}

func NewResult[T any](items []T, total int64, q Query) Result[T] {
	return Result[T]{
		Data: items,
		Meta: NewMeta(total, q),
	}
}