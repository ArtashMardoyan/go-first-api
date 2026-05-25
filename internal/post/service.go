package post

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) FindAll() []Post {
	return s.repo.FindAll()
}

func (s *Service) FindOne(id string) (Post, error) {
	return s.repo.FindByID(id)
}

func (s *Service) FindByUserID(userID string) []Post {
	return s.repo.FindByUserID(userID)
}

func (s *Service) Create(dto CreatePostDto) (Post, error) {
	return s.repo.Create(dto)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}