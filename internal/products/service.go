package products

import "context"

type Service interface {
	GetAllProducts(ctx context.Context) (error)
}

type service struct {
	
}

func NewService() Service {
	return &service{}
}

func (s *service) GetAllProducts(ctx context.Context) (error) {
	return nil
}
