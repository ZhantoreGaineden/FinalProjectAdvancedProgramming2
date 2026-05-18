package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/proto/gen/userpb"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/user-service/internal/config"
	grpcdelivery "github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/user-service/internal/delivery/grpc"
	postgresrepo "github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/user-service/internal/repository/postgres"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/user-service/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("failed to ping postgres: %v", err)
	}

	natsConn, err := nats.Connect(cfg.NATSURL)
	if err != nil {
		log.Printf("warning: failed to connect to nats: %v", err)
		natsConn = nil
	}
	if natsConn != nil {
		defer natsConn.Close()
	}

	userRepo := postgresrepo.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, natsConn)
	userHandler := grpcdelivery.NewUserHandler(userUsecase)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	userpb.RegisterUserServiceServer(server, userHandler)

	log.Printf("User Service gRPC server started on port %s", cfg.GRPCPort)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve grpc server: %v", err)
	}
}
