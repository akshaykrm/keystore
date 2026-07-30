package user

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
		db: db,
	}
}

func (r *Repository) Create(newUser User) error {
	query := `
		INSERT INTO users (
			id,
			email,
			password_hash,
			name,
			created_at,
			updated_at
		) 
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		newUser.ID,
		newUser.Email,
		newUser.Password,
		newUser.Name,
		newUser.CreatedAt,
		newUser.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	return nil

}

func (r *Repository) GetAll() ([]User, error) {
	query := `SELECT id, email, name, created_at, updated_at  FROM users;`

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
			return nil, fmt.Errorf("Scan users: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)

	}

	return users, nil

}

func (r *Repository) GetByID(ID string) (User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?;`

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

	_, err := r.db.Exec(query, user.Email, user.Name, user.UpdatedAt, user.ID)

	fmt.Println(err)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}
