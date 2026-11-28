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

// @Summary Create a product
// @Tags products
// @Accept json
// @Produce json
// @Param request body CreateProductRequest true "product to create"
// @Success 201 {object} CreateProductResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /products [post]
func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Json(w, httpx.Response[any]{
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
			httpx.Json(w, httpx.Response[any]{
				Message: "Invalid request",
				Status:  http.StatusBadRequest,
				Error:   apperr.ErrInvalidRequest,
				Data:    err,
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
		if errors.Is(err, apperr.ErrUniqueViolation) {
			httpx.Json(w, httpx.Response[any]{
				Message: "A product with SKU already exists.",
				Error:   err,
				Status:  http.StatusConflict,
			})
			return
		}
		httpx.InternalServerError(w, r, err)
		return
	}

	responseObject := ProductResponse{
		ID:           p.ID,
		SkuID:        p.SkuID,
		ProductID:    p.ProductID,
		Name:         p.Name,
		Descriptions: p.Descriptions,
		WorkspaceID:  p.WorkspaceID,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}

	httpx.Json(w, httpx.Response[ProductResponse]{
		Message: "product created",
		Status:  http.StatusCreated,
		Data:    responseObject,
	})
}

func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := c.service.GetAll(r.Context())

	if err != nil {
		httpx.InternalServerError(w, r, err)
		return
	}

	httpx.Json(w, httpx.Response[any]{
		Message: "product list",
		Data:    products,
		Status:  http.StatusOK,
	})

}

func (c *Controller) GetOneByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_, ok := auth.GetClaims(r)

	if !ok {
		httpx.Json(w, httpx.Response[any]{
			Message: "user denied access",
			Error:   apperr.ErrAccessDenied,
			Status:  http.StatusUnauthorized,
		})
		return
	}
	num, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		httpx.Json(w, httpx.Response[any]{
			Message: "Product id invalid",
			Status:  http.StatusBadRequest,
		})
		return
	}
	product, err := c.service.GetOneByID(r.Context(), num)

	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			httpx.Json(w, httpx.Response[any]{
				Message: "the requested product does not exist",
				Error:   err,
				Status:  http.StatusNotFound,
			})
			return
		}
		httpx.InternalServerError(w, r, err)
		return
	}
	httpx.Json(w, httpx.Response[any]{
		Message: "product found",
		Data:    product,
		Status:  http.StatusOK,
	})
}

func (c *Controller) UpdateByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue(("id"))
	var body updateProductRequest

	if err := httpx.Decode(r.Body, &body); err != nil {
		httpx.Json(w, httpx.Response[any]{
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
			httpx.Json(w, httpx.Response[any]{
				Message: "Invalid request",
				Status:  http.StatusBadRequest,
				Error:   err,
			})
			return
		}
	}

	num, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		httpx.Json(w, httpx.Response[any]{
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
		httpx.Json(w, httpx.Response[any]{
			Message: "something went wrong while updating workspace",
			Error:   err,
			Status:  http.StatusInternalServerError,
		})
		return
	}
	httpx.Json(w, httpx.Response[any]{
		Message: "Updated product",
		Status:  http.StatusOK,
	})
}
