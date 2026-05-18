package main

import (
	"context"
	"fmt"
	"log"
	"net"

	grpcdelivery "github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/pet-service/internal/delivery/grpc"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/pet-service/internal/config"
	postgresrepo "github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/pet-service/internal/repository/postgres"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/pet-service/internal/usecase"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/proto/gen/petpb"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	defer redisClient.Close()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("warning: failed to ping redis: %v", err)
	}

	petRepo := postgresrepo.NewPetRepository(db)
	petUsecase := usecase.NewPetUsecase(petRepo, redisClient)
	petHandler := grpcdelivery.NewPetHandler(petUsecase)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	petpb.RegisterPetServiceServer(server, petHandler)

	log.Printf("Pet Service gRPC server started on port %s", cfg.GRPCPort)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve grpc server: %v", err)
	}
}
