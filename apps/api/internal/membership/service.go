package membership

import "time"

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Create(newMembership CreateMembershipInput) error {
	now := time.Now().UTC()

	membership := Membership{
		WorkspaceID: newMembership.WorkspaceID,
		userID:      newMembership.userID,
		Role:        newMembership.Role,
	}

	err := s.repo.Create(membership)

}
