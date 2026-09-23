package client

import (
	"context"

	analyticspb "order/internal/grpc/analytics"

	"google.golang.org/grpc"
)

type AnalyticsClient struct {
	client analyticspb.AnalyticsServiceClient
	conn   *grpc.ClientConn
}

func NewAnalyticsClient(address string) (*AnalyticsClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &AnalyticsClient{
		client: analyticspb.NewAnalyticsServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *AnalyticsClient) RecordOrder(
	ctx context.Context,
	id string,
	productID string,
	quantity int,
	price float64,
	status string,
	createdAt string,
) error {
	_, err := c.client.RecordOrder(
		ctx,
		&analyticspb.RecordOrderRequest{
			Id:        id,
			ProductId: productID,
			Quantity:  int32(quantity),
			Price:     price,
			Status:    status,
			CreatedAt: createdAt,
		},
	)

	return err
}

func (c *AnalyticsClient) Close() error {
	return c.conn.Close()
}
