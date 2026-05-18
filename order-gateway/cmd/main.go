package main

import (
	"fmt"
	"log"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/order-gateway/internal/client"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/order-gateway/internal/config"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/order-gateway/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	orderClient, err := client.NewOrderClient(cfg.OrderServiceAddr)
	if err != nil {
		log.Fatalf("failed to connect to order service: %v", err)
	}
	defer orderClient.Close()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "order-gateway is running"})
	})

	orderHandler := handler.NewOrderHandler(orderClient)
	orderHandler.RegisterRoutes(router)

	addr := fmt.Sprintf(":%s", cfg.HTTPPort)
	log.Printf("Order Gateway HTTP server started on port %s", cfg.HTTPPort)

	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to run http server: %v", err)
	}
}
