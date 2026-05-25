package post

import "go-first-api/pkg/pagination"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) FindAll(q pagination.Query) pagination.Result[Post] {
	posts, total := s.repo.FindAll(q)
	return pagination.NewResult(posts, total, q)
}

func (s *Service) FindOne(id string) (Post, error) {
	return s.repo.FindByID(id)
}

func (s *Service) FindByUserID(userID string, q pagination.Query) pagination.Result[Post] {
	posts, total := s.repo.FindByUserID(userID, q)
	return pagination.NewResult(posts, total, q)
}

func (s *Service) Create(dto CreatePostDto) (Post, error) {
	return s.repo.Create(dto)
}

func (s *Service) UpdateStatus(id string, dto UpdatePostStatusDto) (Post, error) {
	if !dto.Status.IsValid() {
		return Post{}, ErrInvalidStatus
	}
	return s.repo.UpdateStatus(id, dto.Status)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}