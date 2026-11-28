package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/akshaykrm/keystore/apps/api/internal/apperr"
	"github.com/akshaykrm/keystore/apps/api/internal/auth"
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

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		httpx.Json(w,
			httpx.Response[any]{
				Message: "invalid payload",
				Status:  http.StatusBadRequest,
			},
		)
		return
	}

	err := c.validate.Struct(loginReq)
	if err != nil {
		httpx.Json(w, httpx.Response[any]{
			Message: "validation error",
			Status:  http.StatusBadRequest,
		})
		return
	}

	token, err := c.service.Login(r.Context(), loginReq)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			httpx.Json(w, httpx.Response[any]{
				Message: "the requested user does not exist",
				Error:   err,
				Status:  http.StatusNotFound,
			})
			return
		}
		if errors.Is(err, apperr.ErrAccessDenied) {
			httpx.Json(w, httpx.Response[any]{
				Message: "resource restricted",
				Error:   err,
				Status:  http.StatusUnauthorized,
			})
			return
		}

		httpx.InternalServerError(w, r, err)
		return

	}

	httpx.Json(w, httpx.Response[any]{
		Message: "authentication successful",
		Data:    LoginResponse{token},
		Status:  http.StatusOK,
	})
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Json(w,
			httpx.Response[any]{
				Message: "invalid payload",
				Status:  http.StatusBadRequest,
			},
		)
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
			httpx.Json(w, httpx.Response[any]{
				Message: "invalid request",
				Data:    errs,
				Error:   apperr.ErrInvalidRequest,
				Status:  http.StatusBadRequest,
			})
			return
		}
	}

	user := CreateUserInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	err = c.service.Create(r.Context(), user)

	if err != nil {
		if errors.Is(err, apperr.ErrUniqueViolation) {
			httpx.Json(w, httpx.Response[any]{
				Message: "please chose another email",
				Error:   err,
				Status:  http.StatusConflict,
			})
			return
		}

	}

	httpx.Json(w, httpx.Response[any]{
		Message: "user created",
		Status:  http.StatusCreated,
	})
}

func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	_, ok := auth.GetClaims(r)
	if !ok {
		httpx.Json(w,
			httpx.Response[any]{
				Message: "user denied access",
				Error:   apperr.ErrAccessDenied,
				Status:  http.StatusUnauthorized,
			},
		)
		return
	}

	users, err := c.service.GetAll(r.Context())

	if err != nil {
		httpx.InternalServerError(w, r, err)
		return
	}

	httpx.Json(w, httpx.Response[any]{
		Message: "active users list",
		Data:    users,
		Status:  http.StatusOK,
	})
}

func (c *Controller) GetByID(w http.ResponseWriter, r *http.Request) {
	_, ok := auth.GetClaims(r)
	if !ok {
		httpx.Json(w,
			httpx.Response[any]{
				Message: "user denied access",
				Error:   apperr.ErrAccessDenied,
				Status:  http.StatusUnauthorized,
			},
		)
		return
	}

	id := r.PathValue("id")

	num, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		httpx.Json(w, httpx.Response[any]{
			Message: "user id invalid",
			Status:  http.StatusBadRequest,
		})
		return
	}

	userResp, err := c.service.GetByID(r.Context(), num)

	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			httpx.Json(w, httpx.Response[any]{
				Message: "requested user not found",
				Error:   err,
				Status:  http.StatusNotFound,
			})
			return
		}
		httpx.InternalServerError(w, r, err)
		return
	}

	httpx.Json(w, httpx.Response[any]{
		Message: "user retrieved by id",
		Data:    userResp,
		Status:  http.StatusOK,
	})
}

func (c *Controller) UpdateById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Json(w,
			httpx.Response[any]{
				Message: "invalid payload",
				Status:  http.StatusBadRequest,
			},
		)
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
			httpx.Json(w, httpx.Response[any]{
				Message: "validation error",
				Data:    errs,
				Error:   apperr.ErrInvalidRequest,
				Status:  http.StatusBadRequest,
			})
			return
		}

	}

	user := UpdateUserInput{
		Email: req.Email,
		Name:  req.Name,
	}

	num, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		httpx.Json(w, httpx.Response[any]{
			Message: "Product id invalid",
			Status:  http.StatusBadRequest,
		})
		return
	}

	if _, err := c.service.UpdateById(r.Context(), num, user); err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			httpx.Json(w,
				httpx.Response[any]{
					Message: "user not found",
					Error:   apperr.ErrNotFound,
					Status:  http.StatusNotFound,
				},
			)
			return
		}
		if errors.Is(err, apperr.ErrUniqueViolation) {
			httpx.Json(w,
				httpx.Response[any]{
					Message: "email already taken",
					Error:   apperr.ErrUniqueViolation,
					Status:  http.StatusConflict,
				},
			)
			return
		}
		httpx.InternalServerError(w, r, err)
		return
	}

	httpx.Json(w, httpx.Response[any]{
		Message: "user updated",
		Status:  http.StatusOK,
	})
}

func (c *Controller) DeleteById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	num, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		httpx.Json(w, httpx.Response[any]{
			Message: "Product id invalid",
			Status:  http.StatusBadRequest,
		})
		return
	}

	if err := c.service.DeleteById(r.Context(), num); err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			httpx.Json(w,
				httpx.Response[any]{
					Message: "user not found",
					Error:   apperr.ErrNotFound,
					Status:  http.StatusNotFound,
				},
			)
			return
		}
		httpx.InternalServerError(w, r, err)
		return
	}

	httpx.Json(w, httpx.Response[any]{
		Message: "user deleted",
		Status:  http.StatusOK,
	})
}
