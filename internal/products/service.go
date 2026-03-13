package products

import (
	"context"
	repo "kwadw0/gocommerce/internal/adapters/postgres/sqlc"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	GetAllProducts(ctx context.Context) ([] repo.Product, error)
	AddProduct(ctx context.Context, product CreateProduct) (repo.Product, error)
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &service{repo: repo}
}

func (s *service) GetAllProducts(ctx context.Context) ([] repo.Product, error) {
	products, err := s.repo.ListProducts(ctx)
	return products, err
}

func (s *service) AddProduct(ctx context.Context, payload CreateProduct) (repo.Product, error) {
	product, err := s.repo.CreateProduct(ctx, repo.CreateProductParams{
		Name:        payload.Name,
		Description: pgtype.Text{String: payload.Description, Valid: true},
		Price:       pgtype.Numeric{Int: big.NewInt(int64(payload.Price * 100)), Exp: -2, Valid: true},
		Quantity:    int32(payload.Quantity),
	})
	return product, err
}


// Add this helper near the bottom
func mapToProductResponse(p repo.Product) ProductResponse {
	// We need to carefully take the price out of its "Numeric" box
	priceF, _ := p.Price.Float64Value()

	return ProductResponse{
		UUID:        p.Uuid,
		Name:        p.Name,
		Description: p.Description.String,
		Price:       priceF.Float64,
		Quantity:    p.Quantity,
	}
}
