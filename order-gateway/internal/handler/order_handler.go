package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/order-gateway/internal/client"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/proto/gen/orderpb"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderClient *client.OrderClient
}

func NewOrderHandler(orderClient *client.OrderClient) *OrderHandler {
	return &OrderHandler{orderClient: orderClient}
}

type createOrderRequest struct {
	UserID    string            `json:"user_id" binding:"required"`
	UserEmail string            `json:"user_email" binding:"required"`
	Items     []createOrderItem `json:"items" binding:"required"`
}

type createOrderItem struct {
	PetID string  `json:"pet_id" binding:"required"`
	Price float64 `json:"price" binding:"required"`
}

type updateOrderStatusRequest struct {
	Status    string `json:"status" binding:"required"`
	UserEmail string `json:"user_email"`
}

func (h *OrderHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/orders", h.CreateOrder)
		api.GET("/orders/:id", h.GetOrder)
		api.GET("/users/:id/orders", h.ListUserOrders)
		api.PATCH("/orders/:id/status", h.UpdateOrderStatus)
		api.POST("/orders/:id/cancel", h.CancelOrder)
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items := make([]*orderpb.CreateOrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, &orderpb.CreateOrderItem{
			PetId: item.PetID,
			Price: item.Price,
		})
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	res, err := h.orderClient.CreateOrder(ctx, &orderpb.CreateOrderRequest{
		UserId:    req.UserID,
		UserEmail: req.UserEmail,
		Items:     items,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res.GetOrder())
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := withTimeout(c)
	defer cancel()

	res, err := h.orderClient.GetOrder(ctx, &orderpb.GetOrderRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res.GetOrder())
}

func (h *OrderHandler) ListUserOrders(c *gin.Context) {
	userID := c.Param("id")

	ctx, cancel := withTimeout(c)
	defer cancel()

	res, err := h.orderClient.ListUserOrders(ctx, &orderpb.ListUserOrdersRequest{UserId: userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": res.GetOrders()})
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var req updateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	res, err := h.orderClient.UpdateOrderStatus(ctx, &orderpb.UpdateOrderStatusRequest{
		Id:        id,
		Status:    req.Status,
		UserEmail: req.UserEmail,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res.GetOrder())
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := withTimeout(c)
	defer cancel()

	res, err := h.orderClient.CancelOrder(ctx, &orderpb.CancelOrderRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func withTimeout(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), 5*time.Second)
}
