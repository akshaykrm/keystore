package user

import (
	"fmt"
	"time"

	"github.com/akshaykrm/keystore/apps/api/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	repo      *Repository
	workspace workspaceCreator
}

func NewService(r *Repository, w workspaceCreator) *Service {
	return &Service{repo: r, workspace: w}
}

func (s *Service) Create(newUser CreateUserInput) error {
	now := time.Now().UTC()
	hashedPassword, err := auth.HashPassword(newUser.Password)
	if err != nil {
		return err
	}

	user := User{
		Email:     newUser.Email,
		Password:  hashedPassword,
		Name:      newUser.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(user); err != nil {
		return err
	}

	if err := s.workspace.CreateDefaultWorkspaceForUser(newUser.Name); err != nil {
		return err
	}

	//TODO: after creating the workspace add a record in the membership table with owner role

	return nil
}

func (s *Service) GetAll() ([]UserResponse, error) {
	users, err := s.repo.GetAll()

	if err != nil {
		return nil, err
	}

	userResponse := make([]UserResponse, 0, len(users))
	for _, user := range users {
		userResponse = append(userResponse, toUserResponse(user))
	}
	return userResponse, nil
}

func (s *Service) Login(loginReq LoginRequest) (string, error) {
	user, err := s.repo.GetPasswordByEmail(loginReq.Email)
	if err != nil {
		return "", err
	}

	if err := auth.CompareHashedPassword(user.Password, loginReq.Password); err != nil {
		return "", fmt.Errorf("Password check failed: %w", err)
	}

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(24 * time.Hour)
	claims := auth.Claims{
		User:      user.Email,
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}
	token, err := auth.CreateNewToken(claims)
	if err != nil {

		return "", fmt.Errorf("Token generation failed: %w", err)
	}
	return token, nil
}

func (s *Service) GetByEmail(Email string) (UserResponse, error) {
	user, err := s.repo.GetByEmail(Email)
	if err != nil {
		return UserResponse{}, err
	}
	return toUserResponse(user), nil
}

func (s *Service) GetByID(ID string) (UserResponse, error) {
	user, err := s.repo.GetByID(ID)
	if err != nil {
		return UserResponse{}, err
	}

	return toUserResponse(user), nil
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

func toUserResponse(user User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
