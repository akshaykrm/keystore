package products

type CreateProductRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=30"`
	Description string `json:"description" validate:"required,min=6,max=100"`
	SkuID       string `json:"sku_id" validate:"required,min=6,max=10"`
	ProductID   string `json:"product_id" validate:"required,min=6,max=12"`
	WorkspaceID int64  `json:"workspace_id" validate:"required"`
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
