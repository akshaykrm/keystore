package user

import (
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Create(newUser CreateUserInput) error {
	now := time.Now().UTC()

	user := User{
		ID:        "01",
		Email:     newUser.Email,
		Password:  newUser.Password,
		Name:      newUser.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return s.repo.Create(user)

}

func (s *Service) GetAll() ([]UserResponse, error) {
	users, err := s.repo.GetAll()

	if err != nil {
		return nil, err
	}

	userResponse := make([]UserResponse, 0, len(users))
	for _, user := range users {
		userResponse = append(userResponse, UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
	return userResponse, nil
}

func (s *Service) GetByID(ID string) (UserResponse, error) {
	user, err := s.repo.GetByID(ID)
	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *Service) UpdateById(ID string, req UpdateUserInput) (UserResponse, error) {
	user, err := s.repo.GetByID(ID)
	if err != nil {
		return UserResponse{}, err
	}

	now := time.Now().UTC()
	user.Name = req.Name
	user.Email = req.Email
	user.UpdatedAt = now

	if err := s.repo.UpdateByID(user); err != nil {
		return UserResponse{}, err
	}

	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *Service) DeleteById(ID string) error {
	user, err := s.repo.GetByID(ID)

	if err != nil {
		return err
	}

	now := time.Now().UTC()
	user.DeletedAt = &now

	return s.repo.DeleteById(user)
}
