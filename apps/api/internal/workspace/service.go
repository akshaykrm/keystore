package workspace

import (
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) Create(w CreateWorkspacePayload) error {
	workspace := Workspace{
		Name: w.Name,
		Slug: w.Slug,
	}

	return s.repo.Create(workspace)
}

func (s *Service) GetAll() ([]WorkspaceList, error) {
	workspaces, err := s.repo.GetAll()

	if err != nil {
		return nil, err
	}

	workspaceList := make([]WorkspaceList, 0, len(workspaces))
	for _, w := range workspaces {
		workspaceList = append(workspaceList, toWorkspaceResponse(w))
	}
	return workspaceList, nil

}

func (s *Service) GetById(ID string) (WorkspaceList, error) {
	workspace, err := s.repo.GetById(ID)
	if err != nil {
		return WorkspaceList{}, err
	}

	return toWorkspaceResponse(workspace), nil

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
