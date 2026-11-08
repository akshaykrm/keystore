package products

import "time"

type Product struct {
	ID          string
	SkuID       string
	ProductID   string
	Name        string
	Description string
	WorkspaceID string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
