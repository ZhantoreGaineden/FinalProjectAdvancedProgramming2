package client

import (
	"context"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/proto/gen/orderpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderClient struct {
	conn   *grpc.ClientConn
	client orderpb.OrderServiceClient
}

func NewOrderClient(addr string) (*OrderClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &OrderClient{
		conn:   conn,
		client: orderpb.NewOrderServiceClient(conn),
	}, nil
}

func (c *OrderClient) Close() error {
	return c.conn.Close()
}

func (c *OrderClient) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.OrderResponse, error) {
	return c.client.CreateOrder(ctx, req)
}

func (c *OrderClient) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.OrderResponse, error) {
	return c.client.GetOrder(ctx, req)
}

func (c *OrderClient) ListUserOrders(ctx context.Context, req *orderpb.ListUserOrdersRequest) (*orderpb.ListOrdersResponse, error) {
	return c.client.ListUserOrders(ctx, req)
}

func (c *OrderClient) UpdateOrderStatus(ctx context.Context, req *orderpb.UpdateOrderStatusRequest) (*orderpb.OrderResponse, error) {
	return c.client.UpdateOrderStatus(ctx, req)
}

func (c *OrderClient) CancelOrder(ctx context.Context, req *orderpb.CancelOrderRequest) (*orderpb.CancelOrderResponse, error) {
	return c.client.CancelOrder(ctx, req)
}
