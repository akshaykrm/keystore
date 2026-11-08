-- name: CreateMembership :exec
INSERT INTO
    memberships (
        role,
        workspace_id,
        user_id,
        created_at,
        updated_at
    )
VALUES
    (?, ?, ?, ?, ?)
