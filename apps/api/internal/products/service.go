package products

import (
	"context"
	"fmt"
	"time"

	"github.com/akshaykrm/keystore/apps/api/db"
)

type Service struct {
	queries *db.Queries
}

func NewService(q *db.Queries) *Service {
	fmt.Println("Creating new Service:", q)
	return &Service{
		queries: q,
	}
}

func (s *Service) Create(ctx context.Context, payload CreateProduct) (db.Product, error) {
	p, err := s.queries.CreateProducts(ctx, db.CreateProductsParams{
		Name:         payload.Name,
		Descriptions: payload.Description,
		SkuID:        payload.SkuID,
		ProductID:    payload.ProductID,
		WorkspaceID:  payload.WorkspaceID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})

	if err != nil {
		if db.IsUniqueViolation(err) {
			return p, ErrProductUniqueViolation
		}
		return p, fmt.Errorf("product create failed: %w", err)
	}

	return p, nil

}

func (s *Service) GetAll(ctx context.Context) ([]db.GetAllProductsRow, error) {

	products, err := s.queries.GetAllProducts(ctx)

	if err != nil {
		return products, fmt.Errorf("Product get all failed: %w", err)
	}
	return products, nil
}

func (s *Service) GetOneByID(ctx context.Context, id int64) (db.GetOneProductByIDRow, error) {
	product, err := s.queries.GetOneProductByID(ctx, id)

	if err != nil {
		return product, fmt.Errorf("Product get by id failed: %w", err)
	}
	return product, nil
}

func (s *Service) UpdateByID(ctx context.Context, payload updateProduct) error {
	err := s.queries.UpdateProductById(ctx, db.UpdateProductByIdParams{
		ID:           payload.ID,
		SkuID:        payload.SkuID,
		ProductID:    payload.ProductID,
		Name:         payload.Name,
		Descriptions: payload.Description,
		WorkspaceID:  payload.WorkspaceID,
		UpdatedAt:    time.Now(),
	})

	if err != nil {
		return fmt.Errorf("Product update by id failed: %w", err)
	}

	return nil

}
