package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/akshaykrm/keystore/apps/api/internal/httpx"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{
		service: s,
	}
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJsonError(w, "invalid body", http.StatusBadRequest)
		return
	}

	user := CreateUserInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	if err := c.service.Create(user); err != nil {
		if errors.Is(err, ErrEmailConflict) {
			httpx.WriteJsonError(w, "email already exists", http.StatusConflict)
			return
		}
		httpx.WriteJsonError(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	httpx.WriteJsonSuccess(w, "User Created", http.StatusCreated)
}

func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := c.service.GetAll()

	if err != nil {
		httpx.WriteJsonError(w, err.Error(), http.StatusInternalServerError)
	}

	httpx.WriteJsonSuccess(w, users, http.StatusOK)
}

func (c *Controller) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	userResp, err := c.service.GetByID(id)

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			httpx.WriteJsonError(w, "User not found", http.StatusNotFound)
			return
		}
		httpx.WriteJsonError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	httpx.WriteJsonSuccess(w, userResp, http.StatusOK)
}

func (c *Controller) UpdateById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJsonError(w, "invalid body", http.StatusBadRequest)
	}

	user := UpdateUserInput{
		Email: req.Email,
		Name:  req.Name,
	}

	if _, err := c.service.UpdateById(id, user); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			httpx.WriteJsonError(w, "user not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrEmailConflict) {
			httpx.WriteJsonError(w, "email already taken", http.StatusConflict)
			return
		}
		httpx.WriteJsonError(w, "failed to update user", http.StatusInternalServerError)
		return
	}

	httpx.WriteJsonSuccess(w, "User Updated", http.StatusOK)
}

func (c *Controller) DeleteById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := c.service.DeleteById(id); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			httpx.WriteJsonError(w, "user not found", http.StatusNotFound)
			return
		}
		httpx.WriteJsonError(w, "failed to delete user", http.StatusInternalServerError)
	}
	httpx.WriteJsonSuccess(w, "User Deleted", http.StatusOK)
}
