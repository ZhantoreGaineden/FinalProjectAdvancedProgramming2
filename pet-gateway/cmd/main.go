package main

import (
	"fmt"
	"log"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/pet-gateway/internal/client"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/pet-gateway/internal/config"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/pet-gateway/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	petClient, err := client.NewPetClient(cfg.PetServiceAddr)
	if err != nil {
		log.Fatalf("failed to connect to pet service: %v", err)
	}
	defer petClient.Close()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "pet-gateway is running"})
	})

	petHandler := handler.NewPetHandler(petClient)
	petHandler.RegisterRoutes(router)

	addr := fmt.Sprintf(":%s", cfg.HTTPPort)
	log.Printf("Pet Gateway HTTP server started on port %s", cfg.HTTPPort)

	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to run http server: %v", err)
	}
}
