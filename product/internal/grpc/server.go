package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"product/internal/model"
	"product/internal/service"
)

type ProductService interface {
	GetByID(ctx context.Context, id string) (model.Product, error)
	GetAll(ctx context.Context) ([]model.Product, error)
	DecreaseStock(ctx context.Context, id string, quantity int) error
}

type Server struct {
	UnimplementedProductServiceServer
	service ProductService
}

func NewServer(service ProductService) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) GetProduct(
	ctx context.Context,
	req *GetProductRequest,
) (*GetProductResponse, error) {
	product, err := s.service.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &GetProductResponse{
		Id:    product.ID,
		Name:  product.Name,
		Price: product.Price,
		Stock: int32(product.Stock),
	}, nil
}

func (s *Server) GetProducts(
	ctx context.Context,
	req *GetProductsRequest,
) (*GetProductsResponse, error) {
	products, err := s.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	response := &GetProductsResponse{}

	for _, product := range products {
		response.Products = append(
			response.Products,
			&GetProductResponse{
				Id:    product.ID,
				Name:  product.Name,
				Price: product.Price,
				Stock: int32(product.Stock),
			},
		)
	}

	return response, nil
}

func (s *Server) DecreaseStock(
	ctx context.Context,
	req *DecreaseStockRequest,
) (*DecreaseStockResponse, error) {
	err := s.service.DecreaseStock(
		ctx,
		req.Id,
		int(req.Quantity),
	)
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			err.Error(),
		)
	}

	return &DecreaseStockResponse{
		Success: true,
	}, nil
}

// compile-time check
var _ ProductService = (*service.ProductService)(nil)
