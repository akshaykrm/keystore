package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/akshaykrm/keystore/apps/api/db"
	"github.com/akshaykrm/keystore/apps/api/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	query      db.Querier
	workspace  workspaceCreator
	membership membershiptCreator
}

func NewService(q db.Querier, w workspaceCreator, m membershiptCreator) *Service {
	return &Service{
		query:      q,
		workspace:  w,
		membership: m,
	}
}

func (s *Service) Create(ctx context.Context, newUser CreateUserInput) error {
	now := time.Now().UTC()
	hashedPassword, err := auth.HashPassword(newUser.Password)
	if err != nil {
		return err
	}

	user := db.CreateUserParams{
		Email:        newUser.Email,
		PasswordHash: hashedPassword,
		Name:         newUser.Name,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	createdUser, err := s.query.CreateUser(ctx, user)

	if err != nil {
		fmt.Printf("User create failed %w", err)
		return err
	}

	workspaceID, err := s.workspace.CreateDefaultWorkspaceForUser(ctx, createdUser.Name)
	if err != nil {
		fmt.Printf("Workspace create failed %w", err)
		return err
	}

	fmt.Printf("Workspace Created with ID: %s \n", workspaceID)

	if err := s.membership.CreateDefaultMembership(ctx, workspaceID, createdUser.ID); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetAll(ctx context.Context) ([]db.GetAllUsersRow, error) {
	users, err := s.query.GetAllUsers(ctx)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Service) Login(ctx context.Context, loginReq LoginRequest) (string, error) {
	user, err := s.query.GetUserPasswordByEmail(ctx, loginReq.Email)
	if err != nil {
		return "", err
	}

	if err := auth.CompareHashedPassword(user.PasswordHash, loginReq.Password); err != nil {
		return "", fmt.Errorf("Password check failed: %w", err)
	}

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(24 * time.Hour)
	claims := auth.Claims{
		User:      user.Email,
		ID:        user.ID,
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}
	token, err := auth.CreateNewToken(claims)
	if err != nil {

		return "", fmt.Errorf("Token generation failed: %w", err)
	}
	return token, nil
}

func (s *Service) GetByEmail(ctx context.Context, Email string) (db.GetUserByEmailRow, error) {
	user, err := s.query.GetUserByEmail(ctx, Email)
	if err != nil {
		return db.GetUserByEmailRow{}, err
	}
	return user, nil
}

func (s *Service) GetByID(ctx context.Context, ID int64) (db.GetUserByIDRow, error) {
	user, err := s.query.GetUserByID(ctx, ID)
	if err != nil {
		return db.GetUserByIDRow{}, err
	}

	return user, nil
}

func (s *Service) UpdateById(ctx context.Context, ID int64, req UpdateUserInput) (UserResponse, error) {
	user, err := s.query.GetUserByID(ctx, ID)
	if err != nil {
		return UserResponse{}, err
	}

	now := time.Now().UTC()
	payload := db.UpdateUserByIDParams{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		UpdatedAt: now,
	}

	if err := s.query.UpdateUserByID(ctx, payload); err != nil {
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

func (s *Service) DeleteById(ctx context.Context, ID int64) error {
	now := sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}
	payload := db.DeleteUserByIDParams{
		DeletedAt: now,
		ID:        ID,
	}

	return s.query.DeleteUserByID(ctx, payload)
}
