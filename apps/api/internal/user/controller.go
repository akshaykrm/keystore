package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/akshaykrm/keystore/apps/api/internal/httpx"
	"github.com/go-playground/validator/v10"
)

type Controller struct {
	service  *Service
	validate *validator.Validate
}

func NewController(s *Service) *Controller {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Controller{
		service:  s,
		validate: v,
	}
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err := c.validate.Struct(req)
	if err != nil {
		var validationError validator.ValidationErrors
		if errors.As(err, &validationError) {
			errs := make(map[string]string)
			for _, v := range validationError {
				switch v.Tag() {
				case "required":
					errs[v.Field()] = "is required"
				case "email":
					errs[v.Field()] = "must be valid email"
				case "min":
					errs[v.Field()] = fmt.Sprintf("must be at least %s characters", v.Param())
				case "max":
					errs[v.Field()] = fmt.Sprintf("must be at most %s characters", v.Param())
				}
			}
			resp := httpx.ErrorResponse{
				Message: "validation error",
				Error:   errs,
			}
			httpx.Error2(w, resp, http.StatusBadRequest)
			return
		}
	}

	user := CreateUserInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	if err := c.service.Create(user); err != nil {
		if errors.Is(err, ErrEmailConflict) {
			httpx.Error(w, "email already exists", http.StatusConflict)
			return
		}
		httpx.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	resp := httpx.SuccessResponse{
		Message: "user created",
	}

	httpx.Write(w, resp, http.StatusCreated)
}

func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := c.service.GetAll()

	if err != nil {
		httpx.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := httpx.SuccessResponse{
		Data:    users,
		Message: "user list",
	}
	httpx.Write(w, resp, http.StatusOK)
}

func (c *Controller) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	userResp, err := c.service.GetByID(id)

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			httpx.Error(w, "User not found", http.StatusNotFound)
			return
		}
		httpx.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp := httpx.SuccessResponse{
		Data:    userResp,
		Message: "user",
	}

	httpx.Write(w, resp, http.StatusOK)
}

func (c *Controller) UpdateById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err := c.validate.Struct(req)
	if err != nil {
		var validationError validator.ValidationErrors
		if errors.As(err, &validationError) {
			errs := make(map[string]string)
			for _, v := range validationError {
				switch v.Tag() {
				case "email":
					errs[v.Field()] = "must be valid email"
				case "min":
					errs[v.Field()] = fmt.Sprintf("must be at least %s characters", v.Param())
				case "max":
					errs[v.Field()] = fmt.Sprintf("must be at most %s characters", v.Param())

				}
			}
			resp := httpx.ErrorResponse{
				Message: "validation error",
				Error:   errs,
			}
			httpx.Error2(w, resp, http.StatusBadRequest)
			return
		}

	}

	user := UpdateUserInput{
		Email: req.Email,
		Name:  req.Name,
	}

	if _, err := c.service.UpdateById(id, user); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			httpx.Error(w, "user not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrEmailConflict) {
			httpx.Error(w, "email already taken", http.StatusConflict)
			return
		}
		httpx.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}

	resp := httpx.SuccessResponse{
		Message: "user updated",
	}

	httpx.Write(w, resp, http.StatusOK)
}

func (c *Controller) DeleteById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := c.service.DeleteById(id); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			httpx.Error(w, "user not found", http.StatusNotFound)
			return
		}
		httpx.Error(w, "failed to delete user", http.StatusInternalServerError)
		return
	}

	resp := httpx.SuccessResponse{
		Message: "user deleted",
	}

	httpx.Write(w, resp, http.StatusOK)
}
