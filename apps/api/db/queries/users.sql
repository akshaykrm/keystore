-- name: CreateUser :one
INSERT INTO
    users (
        email,
        password_hash,
        name,
        created_at,
        updated_at
    )
VALUES
    (?, ?, ?, ?, ?)
RETURNING
    id,
    email,
    name,
    created_at,
    updated_at;

-- name: GetAllUsers :many
SELECT
    id,
    email,
    name,
    created_at,
    updated_at
FROM
    users
WHERE
    deleted_at IS NULL;

-- name: GetUserPasswordByEmail :one
SELECT
    id,
    password_hash,
    email
FROM
    users
WHERE
    email = ?
    AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT
    id,
    name,
    email,
    created_at,
    updated_at
FROM
    users
WHERE
    email = ?
    AND deleted_at IS NULL;

-- name: GetUserByID :one
SELECT
    id,
    name,
    email,
    created_at,
    updated_at
FROM
    users
WHERE
    id = ?
    AND deleted_at IS NULL;

-- name: UpdateUserByID :exec
UPDATE
    users
SET
    email = ?,
    name = ?,
    updated_at = ?
WHERE
    id = ?;

-- name: DeleteUserByID :exec
UPDATE
    users
SET
    deleted_at = ?
WHERE
    id = ?;
