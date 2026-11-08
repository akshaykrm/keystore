package membership

import (
	"context"
	"fmt"
	"time"

	"github.com/akshaykrm/keystore/apps/api/db"
)

type Service struct {
	query db.Querier
}

func NewService(q db.Querier) *Service {
	return &Service{query: q}
}

func (s *Service) CreateDefaultMembership(ctx context.Context, workspaceID, userID int64) error {
	payload := CreateMembershipInput{
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        RoleOwner,
	}

	return s.Create(ctx, payload)
}

func (s *Service) Create(ctx context.Context, newMembership CreateMembershipInput) error {
	now := time.Now().UTC()

	membership := db.CreateMembershipParams{
		WorkspaceID: newMembership.WorkspaceID,
		UserID:      newMembership.UserID,
		Role:        newMembership.Role,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.query.CreateMembership(ctx, membership); err != nil {
		fmt.Printf("Membership Create Failed %v", err)
		return err
	}

	return nil
}
