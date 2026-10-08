package membership

import (
	"fmt"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) CreateDefaultMembership(workspaceID, userID string) error {
	payload := CreateMembershipInput{
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        RoleOwner,
	}
	fmt.Println(payload, userID)
	return s.Create(payload)
}

func (s *Service) Create(newMembership CreateMembershipInput) error {
	now := time.Now().UTC()

	membership := Membership{
		WorkspaceID: newMembership.WorkspaceID,
		UserID:      newMembership.UserID,
		Role:        newMembership.Role,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(membership); err != nil {
		fmt.Printf("Membership Create Failed %v", err)
		return err
	}

	return nil
}
