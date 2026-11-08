package workspace

import (
	"context"
	"database/sql"
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
	return s.query.GetWorkspaceByID(ctx, ID)
}

func (s *Service) UpdateById(ctx context.Context, ID int64, payload UpdateWorkspaceInput) error {
	now := time.Now().UTC()

	return s.query.UpdateWorkspaceByID(ctx,
		db.UpdateWorkspaceByIDParams{
			ID:        ID,
			Name:      payload.Name,
			Slug:      payload.Slug,
			UpdatedAt: now,
		})
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

	return s.query.DeleteWorkspaceByID(ctx, payload)
}
