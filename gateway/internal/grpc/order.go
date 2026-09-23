package grpc

import (
	"context"

	orderpb "gateway/internal/grpc/order"
	"google.golang.org/grpc"
)

type OrderClient struct {
	client orderpb.OrderServiceClient
	conn   *grpc.ClientConn
}

func NewOrderClient(address string) (*OrderClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &OrderClient{
		client: orderpb.NewOrderServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *OrderClient) GetOrder(
	ctx context.Context,
	id string,
) (*orderpb.GetOrderResponse, error) {
	return c.client.GetOrder(
		ctx,
		&orderpb.GetOrderRequest{
			Id: id,
		},
	)
}

func (c *OrderClient) CreateOrder(
	ctx context.Context,
	productID string,
	quantity int32,
	price float64,
) (*orderpb.CreateOrderResponse, error) {
	return c.client.CreateOrder(
		ctx,
		&orderpb.CreateOrderRequest{
			ProductId: productID,
			Quantity:  quantity,
			Price:     price,
		},
	)
}

func (c *OrderClient) CancelOrder(
	ctx context.Context,
	id string,
) (*orderpb.CancelOrderResponse, error) {
	return c.client.CancelOrder(
		ctx,
		&orderpb.CancelOrderRequest{
			Id: id,
		},
	)
}

func (c *OrderClient) Close() error {
	return c.conn.Close()
}
