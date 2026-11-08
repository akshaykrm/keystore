package membership

import (
	"time"
)

type Membership struct {
	ID          string
	WorkspaceID int64
	UserID      int64
	Role        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
