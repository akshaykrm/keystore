package workspace

import "time"

type CreateWorkspacePayload struct {
	Name   string
	Slug   string
	UserId int64
}

type UpdateWorkspaceInput struct {
	Name string
	Slug string
}

type createWorkspaceRequestBody struct {
	Name string `json:"name" validate:"required,min=3"`
	Slug string `json:"slug" validate:"required,min=3"`
}

type updateWorkspaceRequestBody struct {
	Name string `json:"name" validate:"omitempty,min=3"`
	Slug string `json:"slug" validate:"omitempty,min=3"`
}

type WorkspaceList struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListWorkspaceFilter struct {
	UserID int64
}
