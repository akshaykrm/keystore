package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/akshaykrm/keystore/apps/api/db"
)

type membershiptCreator interface {
	CreateDefaultMembership(ctx context.Context, workspacID, userID int64) error
}

type Service struct {
	query      db.Querier
	membership membershiptCreator
}

func NewService(q db.Querier, m membershiptCreator) *Service {
	return &Service{
		query:      q,
		membership: m,
	}
}

func (s *Service) CreateDefaultWorkspaceForUser(ctx context.Context, name string) (int64, error) {
	workspaceName := name + "s" + " " + "Workspace"
	workspaceSlug := strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(workspaceName), "'", " "), " ", "-")

	payload := CreateWorkspacePayload{
		Name: workspaceName,
		Slug: workspaceSlug,
	}
	created, err := s.Create(ctx, payload)
	return created.ID, err
}

func (s *Service) Create(ctx context.Context, w CreateWorkspacePayload) (db.CreateWorkspaceRow, error) {
	now := time.Now().UTC()
	payload := db.CreateWorkspaceParams{
		Name:      w.Name,
		Slug:      w.Slug,
		CreatedAt: now,
		UpdatedAt: now,
	}

	created, err := s.query.CreateWorkspace(ctx, payload)
	if err != nil {
		return db.CreateWorkspaceRow{}, err
	}

	return db.CreateWorkspaceRow{}, s.membership.CreateDefaultMembership(ctx, created.ID, w.UserId)

}

func (s *Service) GetAll(ctx context.Context, filter ListWorkspaceFilter) ([]db.GetAllWorkspacesRow, error) {
	return s.query.GetAllWorkspaces(ctx, filter.UserID)
}

func (s *Service) GetById(ctx context.Context, ID int64) (db.GetWorkspaceByIDRow, error) {
	workspace, err := s.query.GetWorkspaceByID(ctx, ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return workspace, ErrWorkspaceNotFound
		}
	}
	return workspace, nil

}

func (s *Service) UpdateById(ctx context.Context, ID int64, payload UpdateWorkspaceInput) error {
	now := time.Now().UTC()

	result, err := s.query.UpdateWorkspaceByID(ctx,
		db.UpdateWorkspaceByIDParams{
			ID:        ID,
			Name:      payload.Name,
			Slug:      payload.Slug,
			UpdatedAt: now,
		})

	if err != nil {
		return fmt.Errorf("updating workspace %d: %w", ID, err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected for workspace %d: %w", ID, err)
	} else if count == 0 {
		return ErrWorkspaceNotFound
	}

	return nil

}

func (s *Service) DeleteById(ctx context.Context, ID int64) error {
	now := sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}

	payload := db.DeleteWorkspaceByIDParams{
		ID:        ID,
		DeletedAt: now,
	}

	result, err := s.query.DeleteWorkspaceByID(ctx, payload)

	if err != nil {
		return fmt.Errorf("deleting workspace %d: %w", ID, err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected for workspace %d: %w", ID, err)
	} else if count == 0 {
		return ErrWorkspaceNotFound
	}

	return nil

}
