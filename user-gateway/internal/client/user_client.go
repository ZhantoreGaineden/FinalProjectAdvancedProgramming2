package client

import (
	"context"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/proto/gen/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	conn   *grpc.ClientConn
	client userpb.UserServiceClient
}

func NewUserClient(addr string) (*UserClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &UserClient{
		conn:   conn,
		client: userpb.NewUserServiceClient(conn),
	}, nil
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}

func (c *UserClient) RegisterUser(ctx context.Context, req *userpb.RegisterUserRequest) (*userpb.UserResponse, error) {
	return c.client.RegisterUser(ctx, req)
}

func (c *UserClient) LoginUser(ctx context.Context, req *userpb.LoginUserRequest) (*userpb.LoginResponse, error) {
	return c.client.LoginUser(ctx, req)
}

func (c *UserClient) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	return c.client.GetUser(ctx, req)
}

func (c *UserClient) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {
	return c.client.UpdateUser(ctx, req)
}

func (c *UserClient) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	return c.client.DeleteUser(ctx, req)
}
