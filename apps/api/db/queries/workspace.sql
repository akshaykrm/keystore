-- name: CreateWorkspace :one
INSERT INTO
    workspaces (
        name,
        slug,
        created_at,
        updated_at
    )
VALUES
    (?, ?, ?, ?)
RETURNING
    id,
    name,
    slug,
    created_at,
    updated_at;

-- name: GetAllWorkspaces :many
SELECT
    w.id,
    w.name,
    w.slug,
    w.created_at,
    w.updated_at
FROM
    workspaces w
    INNER JOIN memberships m ON w.id = m.workspace_id
WHERE
    m.user_id = ?
    AND w.deleted_at IS NULL;

-- name: GetWorkspaceByID :one
SELECT
    id,
    name,
    slug,
    created_at,
    updated_at
FROM
    workspaces
WHERE
    id = ?
    AND deleted_at IS NULL;

-- name: GetWorkspaceByUserAndID :one
SELECT
    w.id,
    w.name,
    w.slug,
    w.created_at,
    w.updated_at
FROM
    workspaces w
    INNER JOIN memberships m ON w.id = m.workspace_id
WHERE
    w.id = ?
    AND m.user_id = ?
    AND w.deleted_at IS NULL;

-- name: UpdateWorkspaceByID :execresult
UPDATE
    workspaces
SET
    name = ?,
    slug = ?,
    updated_at = ?
WHERE
    id = ?
    AND deleted_at IS NULL;

-- name: DeleteWorkspaceByID :execresult
UPDATE
    workspaces
SET
    deleted_at = ?
WHERE
    id = ?
    AND deleted_at IS NULL;
