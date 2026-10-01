package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	productpb "marketplace-api/gen/product"
	"product/internal/model"
	"product/internal/service"
)

type ProductService interface {
	GetByID(ctx context.Context, id string) (model.Product, error)
	GetAll(ctx context.Context) ([]model.Product, error)
	DecreaseStock(ctx context.Context, id string, quantity int) error
}

type Server struct {
	productpb.UnimplementedProductServiceServer
	service ProductService
}

func NewServer(service ProductService) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) GetProduct(
	ctx context.Context,
	req *productpb.GetProductRequest,
) (*productpb.GetProductResponse, error) {
	product, err := s.service.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &productpb.GetProductResponse{
		Id:    product.ID,
		Name:  product.Name,
		Price: product.Price,
		Stock: int32(product.Stock),
	}, nil
}

func (s *Server) GetProducts(
	ctx context.Context,
	req *productpb.GetProductsRequest,
) (*productpb.GetProductsResponse, error) {
	products, err := s.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	response := &productpb.GetProductsResponse{}

	for _, product := range products {
		response.Products = append(
			response.Products,
			&productpb.GetProductResponse{
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
	req *productpb.DecreaseStockRequest,
) (*productpb.DecreaseStockResponse, error) {
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

	return &productpb.DecreaseStockResponse{
		Success: true,
	}, nil
}

// compile-time check
var _ ProductService = (*service.ProductService)(nil)
