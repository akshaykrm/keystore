package workspace

import (
	"database/sql"
	"errors"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db,
	}
}

func (r *Repository) Create(w Workspace) (Workspace, error) {
	query := `
		INSERT INTO workspaces (
			name,
			slug,
			created_at,
			updated_at
		) 
		VALUES (?, ?, ?, ?)
		RETURNING id, name, slug, created_at, updated_at
	`

	row := r.db.QueryRow(query, w.Name, w.Slug, w.CreatedAt, w.UpdatedAt)

	created := Workspace{}
	if err := row.Scan(
		&created.ID,
		&created.Name,
		&created.Slug,
		&created.CreatedAt,
		&created.UpdatedAt,
	); err != nil {
		return created, fmt.Errorf("Insert workspace failed: %w", err)
	}

	return created, nil
}

func (r *Repository) GetAll(filter ListWorkspaceFilter) ([]Workspace, error) {
	query := `
		SELECT 
			w.id,
			w.name,
			w.slug,
			w.created_at,
			w.updated_at
		FROM
			workspaces w
		INNER JOIN
			memberships m ON w.id = m.workspace_id

		WHERE
			m.user_id = ?
		AND
			w.deleted_at IS NULL
	`

	args := []any{filter.UserID}
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query workspaces: %w", err)
	}

	defer rows.Close()

	var workspaces []Workspace

	for rows.Next() {
		var workspace Workspace
		if err := rows.Scan(
			&workspace.ID,
			&workspace.Name,
			&workspace.Slug,
			&workspace.CreatedAt,
			&workspace.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan workspaces: %w", err)
		}
		workspaces = append(workspaces, workspace)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate workspaces: %w", err)
	}

	return workspaces, nil

}

func (r *Repository) GetById(ID string) (Workspace, error) {
	query := `SELECT 
			id,
			name,
			slug,
			created_at,
			updated_at 
		FROM 
			workspaces 
		WHERE 
			id = ? 
		AND 
			deleted_at IS NULL;`

	var workspace Workspace
	err := r.db.QueryRow(query, ID).Scan(
		&workspace.ID,
		&workspace.Name,
		&workspace.Slug,
		&workspace.CreatedAt,
		&workspace.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Workspace{}, ErrWorkspaceNotFound
		}
		return Workspace{}, fmt.Errorf("get workspace: %w", err)
	}

	return workspace, nil
}

func (r *Repository) GetByIdAndUser(id, userID string) (*Workspace, error) {
	query := `SELECT 
			w.id,
			w.name,
			w.slug,
			w.created_at,
			w.updated_at 
		FROM 
			workspaces w
		INNER JOIN 
			memberships m
		ON 
			w.id = m.workspace_id
			
		WHERE 
			w.id = ? 
		AND
			m.user_id = ?
		AND 
			w.deleted_at IS NULL;`

	row := r.db.QueryRow(query, id, userID)
	var workspace Workspace

	if err := row.Scan(
		&workspace.ID,
		&workspace.Name,
		&workspace.Slug,
		&workspace.CreatedAt,
		&workspace.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWorkspaceNotFound
		}
		return nil, fmt.Errorf("Failed to query workspaces: %w", err)
	}

	return &workspace, nil
}

func (r *Repository) UpdateById(w Workspace) error {
	query := `
		UPDATE 
			workspaces
		SET
			name = ?,
			slug = ?,
			updated_at = ?
		WHERE
			id = ?
	`

	result, err := r.db.Exec(query, w.Name, w.Slug, w.UpdatedAt, w.ID)

	if err != nil {
		return fmt.Errorf("update workspace: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrWorkspaceNotFound
	}

	return nil
}

func (r *Repository) DeleteById(workspace Workspace) error {
	query := `
		UPDATE
			workspaces
		SET 
			deleted_at = ?
		WHERE
			id = ?;
	`

	result, err := r.db.Exec(query, workspace.DeletedAt, workspace.ID)

	if err != nil {
		return fmt.Errorf("delete workspace: %w", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrWorkspaceNotFound
	}

	return nil
}
