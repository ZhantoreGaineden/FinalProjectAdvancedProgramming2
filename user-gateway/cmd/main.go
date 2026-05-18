package main

import (
	"fmt"
	"log"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/user-gateway/internal/client"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/user-gateway/internal/config"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/user-gateway/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	userClient, err := client.NewUserClient(cfg.UserServiceAddr)
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}
	defer userClient.Close()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "user-gateway is running"})
	})

	userHandler := handler.NewUserHandler(userClient)
	userHandler.RegisterRoutes(router)

	addr := fmt.Sprintf(":%s", cfg.HTTPPort)
	log.Printf("User Gateway HTTP server started on port %s", cfg.HTTPPort)

	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to run http server: %v", err)
	}
}
