package workspace

import (
	"strings"
	"time"
)

type membershiptCreator interface {
	CreateDefaultMembership(workspacID, userID string) error
}

type Service struct {
	repo       *Repository
	membership membershiptCreator
}

func NewService(r *Repository, m membershiptCreator) *Service {
	return &Service{
		repo:       r,
		membership: m,
	}
}

func (s *Service) CreateDefaultWorkspaceForUser(name string) (string, error) {
	workspaceName := name + "s" + " " + "Workspace"
	workspaceSlug := strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(workspaceName), "'", " "), " ", "-")

	payload := CreateWorkspacePayload{
		Name: workspaceName,
		Slug: workspaceSlug,
	}
	created, err := s.Create(payload)
	return created.ID, err
}

func (s *Service) Create(w CreateWorkspacePayload) (Workspace, error) {
	workspace := Workspace{
		Name: w.Name,
		Slug: w.Slug,
	}

	workspace, err := s.repo.Create(workspace)
	if err != nil {
		return Workspace{}, err
	}
	return Workspace{}, s.membership.CreateDefaultMembership(workspace.ID, w.UserId)

}

func (s *Service) GetAll(filter ListWorkspaceFilter) ([]WorkspaceList, error) {
	workspaces, err := s.repo.GetAll(filter)

	if err != nil {
		return nil, err
	}

	workspaceList := make([]WorkspaceList, 0, len(workspaces))
	for _, w := range workspaces {
		workspaceList = append(workspaceList, toWorkspaceResponse(w))
	}
	return workspaceList, nil

}

func (s *Service) GetById(ID, userId string) (WorkspaceList, error) {
	workspace, err := s.repo.GetByIdAndUser(ID, userId)
	if err != nil {
		return WorkspaceList{}, err
	}
	return toWorkspaceResponse(*workspace), nil
}

func (s *Service) UpdateById(ID string, payload UpdateWorkspaceInput) (WorkspaceList, error) {
	workspace, err := s.repo.GetById(ID)
	if err != nil {
		return WorkspaceList{}, err
	}

	now := time.Now().UTC()
	workspace.Name = payload.Name
	workspace.Slug = payload.Slug
	workspace.UpdatedAt = now

	if err := s.repo.UpdateById(workspace); err != nil {
		return WorkspaceList{}, err
	}

	return toWorkspaceResponse(workspace), nil
}

func (s *Service) DeleteById(ID string) error {
	workspace, err := s.repo.GetById(ID)

	if err != nil {
		return err
	}

	now := time.Now().UTC()
	workspace.DeletedAt = &now

	return s.repo.DeleteById(workspace)
}

func toWorkspaceResponse(w Workspace) WorkspaceList {
	return WorkspaceList{
		ID:        w.ID,
		Name:      w.Name,
		Slug:      w.Slug,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}
}
