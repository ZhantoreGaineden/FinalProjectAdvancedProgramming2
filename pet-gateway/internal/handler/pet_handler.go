package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/pet-gateway/internal/client"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/proto/gen/petpb"
	"github.com/gin-gonic/gin"
)

type PetHandler struct {
	petClient *client.PetClient
}

func NewPetHandler(petClient *client.PetClient) *PetHandler {
	return &PetHandler{petClient: petClient}
}

type createPetRequest struct {
	Name     string  `json:"name" binding:"required"`
	Category string  `json:"category" binding:"required"`
	Breed    string  `json:"breed"`
	Age      int32   `json:"age"`
	Price    float64 `json:"price" binding:"required"`
	Status   string  `json:"status"`
}

type updatePetRequest struct {
	Name     string  `json:"name" binding:"required"`
	Category string  `json:"category" binding:"required"`
	Breed    string  `json:"breed"`
	Age      int32   `json:"age"`
	Price    float64 `json:"price" binding:"required"`
	Status   string  `json:"status"`
}

func (h *PetHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/pets")
	{
		api.POST("", h.CreatePet)
		api.GET("", h.ListPets)
		api.GET("/:id", h.GetPet)
		api.PUT("/:id", h.UpdatePet)
		api.DELETE("/:id", h.DeletePet)
	}
}

func (h *PetHandler) CreatePet(c *gin.Context) {
	var req createPetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	res, err := h.petClient.CreatePet(ctx, &petpb.CreatePetRequest{
		Name:     req.Name,
		Category: req.Category,
		Breed:    req.Breed,
		Age:      req.Age,
		Price:    req.Price,
		Status:   req.Status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res.GetPet())
}

func (h *PetHandler) GetPet(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := withTimeout(c)
	defer cancel()

	res, err := h.petClient.GetPet(ctx, &petpb.GetPetRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res.GetPet())
}

func (h *PetHandler) ListPets(c *gin.Context) {
	category := c.Query("category")
	status := c.Query("status")

	ctx, cancel := withTimeout(c)
	defer cancel()

	res, err := h.petClient.ListPets(ctx, &petpb.ListPetsRequest{
		Category: category,
		Status:   status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pets": res.GetPets()})
}

func (h *PetHandler) UpdatePet(c *gin.Context) {
	id := c.Param("id")

	var req updatePetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	res, err := h.petClient.UpdatePet(ctx, &petpb.UpdatePetRequest{
		Id:       id,
		Name:     req.Name,
		Category: req.Category,
		Breed:    req.Breed,
		Age:      req.Age,
		Price:    req.Price,
		Status:   req.Status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res.GetPet())
}

func (h *PetHandler) DeletePet(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := withTimeout(c)
	defer cancel()

	res, err := h.petClient.DeletePet(ctx, &petpb.DeletePetRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func withTimeout(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), 5*time.Second)
}
