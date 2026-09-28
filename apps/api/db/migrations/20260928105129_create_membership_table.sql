-- +goose Up
CREATE TABLE memberships (
  id INTEGER PRIMARY KEY,
  role TEXT NOT NULL CHECK (role IN ('owner', 'member')),
  workspace_id TEXT NOT NULL,
  user_id TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  deleted_at DATETIME DEFAULT NULL,
  FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE (workspace_id, user_id)
);

-- +goose Down
DROP TABLE memberships;
