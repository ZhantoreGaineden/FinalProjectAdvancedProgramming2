package client

import (
	"context"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/proto/gen/petpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PetClient struct {
	conn   *grpc.ClientConn
	client petpb.PetServiceClient
}

func NewPetClient(addr string) (*PetClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &PetClient{
		conn:   conn,
		client: petpb.NewPetServiceClient(conn),
	}, nil
}

func (c *PetClient) Close() error {
	return c.conn.Close()
}

func (c *PetClient) CreatePet(ctx context.Context, req *petpb.CreatePetRequest) (*petpb.PetResponse, error) {
	return c.client.CreatePet(ctx, req)
}

func (c *PetClient) GetPet(ctx context.Context, req *petpb.GetPetRequest) (*petpb.PetResponse, error) {
	return c.client.GetPet(ctx, req)
}

func (c *PetClient) ListPets(ctx context.Context, req *petpb.ListPetsRequest) (*petpb.ListPetsResponse, error) {
	return c.client.ListPets(ctx, req)
}

func (c *PetClient) UpdatePet(ctx context.Context, req *petpb.UpdatePetRequest) (*petpb.PetResponse, error) {
	return c.client.UpdatePet(ctx, req)
}

func (c *PetClient) DeletePet(ctx context.Context, req *petpb.DeletePetRequest) (*petpb.DeletePetResponse, error) {
	return c.client.DeletePet(ctx, req)
}
