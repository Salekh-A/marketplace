package grpc

import (
	"context"

	productpb "gateway/internal/grpc/product"
	"google.golang.org/grpc"
)

type ProductClient struct {
	client productpb.ProductServiceClient
	conn   *grpc.ClientConn
}

func NewProductClient(address string) (*ProductClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &ProductClient{
		client: productpb.NewProductServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *ProductClient) GetProduct(
	ctx context.Context,
	id string,
) (*productpb.GetProductResponse, error) {
	return c.client.GetProduct(
		ctx,
		&productpb.GetProductRequest{
			Id: id,
		},
	)
}

func (c *ProductClient) GetProducts(
	ctx context.Context,
) (*productpb.GetProductsResponse, error) {
	return c.client.GetProducts(
		ctx,
		&productpb.GetProductsRequest{},
	)
}

func (c *ProductClient) Close() error {
	return c.conn.Close()
}
