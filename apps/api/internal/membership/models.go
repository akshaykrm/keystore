package membership

import (
	"time"
)

type Membership struct {
	ID          string
	WorkspaceID string
	UserID      string
	Role        Role
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
