package products

import "time"

// CreateProductRequest is the request body to create a product
type CreateProductRequest struct {
	// Name is the products display name.
	Name string `json:"name" validate:"required,min=3,max=30"`

	// Description is the products description.
	Description string `json:"description" validate:"required,min=6,max=100"`

	// SkuID is unique identifier for the product.
	SkuID string `json:"sku_id" validate:"required,min=6,max=10"`

	// ProductID is the identification for the type of product.
	ProductID string `json:"product_id" validate:"required,min=6,max=12"`

	// WorkspaceID is used to identify in which workspace the product belongs to.
	WorkspaceID int64 `json:"workspace_id" validate:"required"`
}

// ProductResponse is the public representation of a product.
// It is the only product shape sent to clients, so it does not expose
// database-only details such as sql.NullTime.
type ProductResponse struct {
	// ID is the internal identifier of the product.
	ID int64 `json:"id"`

	// SkuID is the stock-keeping unit that identifies the product in inventory.
	SkuID string `json:"sku_id"`

	// ProductID is provisional. Describe it once its source is decided.
	ProductID string `json:"product_id"`

	// Name is the product's display name.
	Name string `json:"name"`

	// Description is a short summary of the product.
	Descriptions interface{} `json:"descriptions"`

	// WorkspaceID is the workspace the product belongs to.
	WorkspaceID int64 `json:"workspace_id"`

	// CreatedAt is when the product was created, in UTC.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is when the product was last changed, in UTC.
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateProductResponse is the body returned when a product is created.
type CreateProductResponse struct {
	// Message is a short summary of the result, such as "product created".
	Message string `json:"message"`

	// Data is the product that was created.
	Data ProductResponse `json:"data"`
}

type updateProductRequest struct {
	Name        string `json:"name" validate:"omitempty,min=3,max=30"`
	Description string `json:"description" validate:"omitempty,min=6,max=100"`
	SkuID       string `json:"sku_id" validate:"omitempty,min=6,max=10"`
	ProductID   string `json:"product_id" validate:"omitempty,min=6,max=12"`
	WorkspaceID int64  `json:"workspace_id" validate:"omitempty"`
}

type CreateProduct struct {
	Name        string
	Description string
	SkuID       string
	ProductID   string
	WorkspaceID int64
}

type updateProduct struct {
	ID          int64
	Name        string
	Description string
	SkuID       string
	ProductID   string
	WorkspaceID int64
}
