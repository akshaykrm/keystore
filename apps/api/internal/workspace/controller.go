package workspace

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/akshaykrm/keystore/apps/api/internal/httpx"
	"github.com/go-playground/validator/v10"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{
		service: s,
	}
}

func newValidator() *validator.Validate {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return v
}

type ValidationErrors map[string]string

func (v ValidationErrors) Error() string {
	return "validation failed"
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var payload createWorkspaceRequestBody

	if err := httpx.Decode(r.Body, &payload); err != nil {
		httpx.Json(w, httpx.Response{
			Message: "Invalid Body",
			Status:  http.StatusBadRequest,
		})
		return
	}

	if err := validate(payload); err != nil {
		httpx.Json(w, httpx.Response{
			Message: "Invalid request",
			Status:  http.StatusBadRequest,
			Error:   err.Error(),
		})
		return

	}

	workspace := CreateWorkspacePayload{
		Name: payload.Name,
		Slug: payload.Slug,
	}

	if err := c.service.Create(workspace); err != nil {
		httpx.Json(w, httpx.Response{
			Message: "something went wrong while creating workspace",
			Status:  http.StatusInternalServerError,
		})
		fmt.Printf("Failed to create workspace: %v", err)
		return
	}

	httpx.Json(w, httpx.Response{
		Message: "Created workspace",
		Status:  http.StatusCreated,
	})
}

func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	workspaces, err := c.service.GetAll()

	if err != nil {
		httpx.Json(w, httpx.Response{
			Message: "something went wrong while retrieving workspaces",
			Status:  http.StatusInternalServerError,
		})
		fmt.Printf("Failed to get workspace list: %v", err)
		return
	}

	httpx.Json(w, httpx.Response{
		Message: "workspace list",
		Data:    workspaces,
		Status:  http.StatusOK,
	})
	return
}

func (c *Controller) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue(("id"))

	workspace, err := c.service.GetById(id)

	if err != nil {
		if errors.Is(err, ErrWorkspaceNotFound) {
			httpx.Json(w, httpx.Response{
				Message: "Workspace not found",
				Status:  http.StatusNotFound,
			})
		} else {
			httpx.Json(w, httpx.Response{
				Message: "something went wrong while retrieving workspace",
				Status:  http.StatusInternalServerError,
			})
			fmt.Printf("Failed to get workspace by id '%s': %v", id, err)
		}
		return
	}

	httpx.Json(w, httpx.Response{
		Message: "Workspace found",
		Data:    workspace,
		Status:  http.StatusOK,
	})
}

func (c *Controller) UpdateByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue(("id"))
	var body updateWorkspaceRequestBody

	if err := httpx.Decode(r.Body, &body); err != nil {
		httpx.Json(w, httpx.Response{
			Message: "Invalid Body",
			Status:  http.StatusBadRequest,
		})
		return
	}

	if err := validate(body); err != nil {
		httpx.Json(w, httpx.Response{
			Message: "Invalid request",
			Status:  http.StatusBadRequest,
			Error:   err.Error(),
		})
		return
	}

	payload := UpdateWorkspaceInput{
		Name: body.Name,
		Slug: body.Slug,
	}
	if _, err := c.service.UpdateById(id, payload); err != nil {
		if err == ErrWorkspaceNotFound {
			httpx.Json(w, httpx.Response{
				Message: "workspace not found",
				Error:   err.Error(),
				Status:  http.StatusNotFound,
			})
			return
		}
		httpx.Json(w, httpx.Response{
			Message: "something went wrong while updating workspace",
			Status:  http.StatusInternalServerError,
		})
		fmt.Printf("Failed to update workspace by id '%s': %v", id, err)
		return
	}

	httpx.Json(w, httpx.Response{
		Message: "successfully updated workspace",
		Status:  http.StatusOK,
	})

}

func (c *Controller) DeleteById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := c.service.DeleteById(id); err != nil {
		if errors.Is(err, ErrWorkspaceNotFound) {
			httpx.Json(w, httpx.Response{
				Message: "workspace not found",
				Error:   err.Error(),
				Status:  http.StatusNotFound,
			})
			return
		}
		httpx.Json(w, httpx.Response{
			Message: "something went wrong while deleting workspace",
			Status:  http.StatusInternalServerError,
		})
		fmt.Printf("Failed to delete workspace by id '%s': %v", id, err)
		return
	}

	httpx.Json(w, httpx.Response{
		Message: "workspace deleted",
		Status:  http.StatusOK,
	})
}
