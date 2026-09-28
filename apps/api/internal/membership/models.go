package membership

import (
	"errors"
	"time"
)

type Role string

const (
	RoleOwner  Role = "owner"
	RoleMember Role = "member"
)

var ErrInvalidRole = errors.New("invalid membership role")

func (r Role) IsValid() bool {
	switch r {
	case RoleOwner, RoleMember:
		return true
	default:
		false
	}
}

type Membership struct {
	ID          string
	WorkspaceID string
	UserID      string
	Role        Role
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
