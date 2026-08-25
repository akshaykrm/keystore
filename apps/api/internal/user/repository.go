package user

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

func (r *Repository) Create(newUser User) error {
	query := `
		INSERT INTO users (
			email,
			password_hash,
			name,
			created_at,
			updated_at
		) 
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		newUser.Email,
		newUser.Password,
		newUser.Name,
		newUser.CreatedAt,
		newUser.UpdatedAt,
	)

	if err != nil {
		sqliteErr, ok := errors.AsType[*sqlite.Error](err)
		if ok {
			switch sqliteErr.Code() {
			case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
				return ErrEmailConflict
			}
		}
		fmt.Println(err)
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *Repository) GetAll() ([]User, error) {
	query := `SELECT id, email, name, created_at, updated_at  FROM users WHERE deleted_at IS NULL;`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}

	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan users: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}

func (r *Repository) GetPasswordByEmail(Email string) (User, error) {
	query := `SELECT id, password_hash, email FROM users WHERE email = ? AND deleted_at IS NULL;`
	var user User
	err := r.db.QueryRow(query, Email).Scan(
		&user.ID,
		&user.Password,
		&user.Email,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}
func (r *Repository) GetByEmail(Email string) (User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE email = ? AND deleted_at IS NULL;`

	var user User
	err := r.db.QueryRow(query, Email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil

}
func (r *Repository) GetByID(ID string) (User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = ? AND deleted_at IS NULL;`

	var user User
	err := r.db.QueryRow(query, ID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

func (r *Repository) UpdateByID(user User) error {
	query := `
	UPDATE 
		users 
	SET 
		email = ?,
		name = ?,
		updated_at = ?
	WHERE id = ?
	`
	result, err := r.db.Exec(query, user.Email, user.Name, user.UpdatedAt, user.ID)

	if err != nil {
		sqliteErr, ok := errors.AsType[*sqlite.Error](err)
		if ok {
			switch sqliteErr.Code() {
			case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
				return ErrEmailConflict
			}
		}

		return fmt.Errorf("update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *Repository) DeleteById(user User) error {
	query := `
	UPDATE
		users
	SET 
		deleted_at = ?
	WHERE
		id = ?;
	`

	result, err := r.db.Exec(query, user.DeletedAt, user.ID)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil

}
