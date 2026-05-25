package user

// аналог @Injectable() Service в NestJS
type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) FindAll() []User {
	return s.repo.FindAll()
}

func (s *Service) FindOne(id string) (User, error) {
	return s.repo.FindByID(id)
}

func (s *Service) FindByEmail(email string) (User, error) {
	return s.repo.FindByEmail(email)
}

func (s *Service) Create(dto CreateUserDto) (User, error) {
	return s.repo.Create(dto)
}

func (s *Service) Update(id string, dto UpdateUserDto) (User, error) {
	return s.repo.Update(id, dto)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}