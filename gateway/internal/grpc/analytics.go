package grpc

import (
	"context"

	"google.golang.org/grpc"
	analyticspb "marketplace-api/gen/analytics"
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

func (c *AnalyticsClient) GetRevenue(
	ctx context.Context,
) (*analyticspb.GetRevenueResponse, error) {
	return c.client.GetRevenue(
		ctx,
		&analyticspb.GetRevenueRequest{},
	)
}

func (c *AnalyticsClient) GetOrders(
	ctx context.Context,
) (*analyticspb.GetOrdersResponse, error) {
	return c.client.GetOrders(
		ctx,
		&analyticspb.GetOrdersRequest{},
	)
}

func (c *AnalyticsClient) Close() error {
	return c.conn.Close()
}
