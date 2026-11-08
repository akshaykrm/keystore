-- +goose Up
CREATE TABLE products (
    id INTEGER PRIMARY KEY,
    sku_id TEXT NOT NULL UNIQUE,
    product_id TEXT NOT NULL,
    name TEXT NOT NULL,
    descriptions name TEXT NOT NULL,
    workspace_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME DEFAULT NULL
);

-- +goose Down
DROP TABLE products;
