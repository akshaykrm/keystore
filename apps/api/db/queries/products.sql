-- name: CreateProducts :one
INSERT INTO
    products (
        sku_id,
        product_id,
        name,
        descriptions,
        workspace_id,
        created_at,
        updated_at
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?)
RETURNING
    *;

-- name: GetAllProducts :many
SELECT
    id,
    sku_id,
    product_id,
    name,
    descriptions,
    workspace_id,
    created_at,
    updated_at
FROM
    products;

-- name: GetOneProductByID :one
SELECT
    id,
    sku_id,
    product_id,
    name,
    descriptions,
    workspace_id,
    created_at,
    updated_at
FROM
    products
WHERE
    id = ?;

-- name: UpdateProductById :exec
UPDATE
    products
SET
    sku_id = ?,
    product_id = ?,
    name = ?,
    descriptions = ?,
    workspace_id = ?,
    updated_at = ?
WHERE
    id = ?;
