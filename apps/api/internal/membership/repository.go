package membership

import (
	"database/sql"
	"errors"
	"fmt"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(newMembership Membership) error {
	query := `
		INSERT INTO memberships (
			role,
			workspace_id,
			user_id,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		newMembership.Role,
		newMembership.WorkspaceID,
		newMembership.UserID,
		newMembership.CreatedAt,
		newMembership.UpdatedAt,
	)

	if err != nil {
		sqliteErr, ok := errors.AsType[*sqlite.Error](err)
		if ok {
			switch sqliteErr.Code() {
			case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
				return ErrMemberShipConflict
			}
		}
		return fmt.Errorf("create membership: %w", err)
	}

	return nil

}
