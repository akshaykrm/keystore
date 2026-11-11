package products

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
	validate *validator.Validate
	service  *Service
}

func NewController(service *Service) *Controller {

	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Controller{
		validate: v,
		service:  service,
	}
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println(err)
		httpx.Json(w, httpx.Response{
			Message: "Invalid Body",
			Status:  http.StatusBadRequest,
		})
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
			httpx.Json(w, httpx.Response{
				Message: "Invalid request",
				Status:  http.StatusBadRequest,
				Error:   err.Error(),
			})
			return
		}
	}

	product := CreateProduct{
		Name:        req.Name,
		Description: req.Description,
		SkuID:       req.SkuID,
		ProductID:   req.ProductID,
		WorkspaceID: req.WorkspaceID,
	}

	p, err := c.service.Create(r.Context(), product)
	if err != nil {
		if errors.Is(err, ErrProductUniqueViolation) {
			httpx.Json(w, httpx.Response{
				Message: "A product with SKU already exists.",
				Error:   err.Error(),
				Status:  http.StatusConflict,
			})
			return
		}
		httpx.InternalServerError(w, r, err)
		return
	}

	httpx.Json(w, httpx.Response{
		Message: "product created",
		Status:  http.StatusCreated,
		Data:    p,
	})
}

func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := c.service.GetAll(r.Context())

	if err != nil {
		httpx.InternalServerError(w, r, err)
		return
	}

	httpx.Json(w, httpx.Response{
		Message: "product list",
		Data:    products,
		Status:  http.StatusOK,
	})

}

func (c *Controller) GetOneByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_, ok := auth.GetClaims(r)

	if !ok {
		httpx.Error(w, auth.AccessDenied.Error(), http.StatusUnauthorized)
		return
	}
	num, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		httpx.Json(w, httpx.Response{
			Message: "Product id invalid",
			Status:  http.StatusBadRequest,
		})
		return
	}
	product, err := c.service.GetOneByID(r.Context(), num)

	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			httpx.Json(w, httpx.Response{
				Message: "the requested product does not exists",
				Error:   err.Error(),
				Status:  http.StatusNotFound,
			})
			return
		}
		httpx.InternalServerError(w, r, err)
		return
	}
	httpx.Json(w, httpx.Response{
		Message: "product found",
		Data:    product,
		Status:  http.StatusOK,
	})
}

func (c *Controller) UpdateByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue(("id"))
	var body updateProductRequest

	if err := httpx.Decode(r.Body, &body); err != nil {
		httpx.Json(w, httpx.Response{
			Message: "Invalid Body",
			Status:  http.StatusBadRequest,
		})
		return
	}

	err := c.validate.Struct(body)
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
			httpx.Json(w, httpx.Response{
				Message: "Invalid request",
				Status:  http.StatusBadRequest,
				Error:   err.Error(),
			})
			return
		}
	}

	num, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		httpx.Json(w, httpx.Response{
			Message: "Product id invalid",
			Status:  http.StatusBadRequest,
		})
		return
	}

	payload := updateProduct{
		ID:          num,
		SkuID:       body.SkuID,
		ProductID:   body.ProductID,
		Name:        body.Name,
		Description: body.Description,
		WorkspaceID: body.WorkspaceID,
	}

	if err := c.service.UpdateByID(r.Context(), payload); err != nil {
		httpx.Json(w, httpx.Response{
			Message: "something went wrong while updating workspace",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}
	httpx.Json(w, httpx.Response{
		Message: "Updated product",
		Status:  http.StatusOK,
	})
}
