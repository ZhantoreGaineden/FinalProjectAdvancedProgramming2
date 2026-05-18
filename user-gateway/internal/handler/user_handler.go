package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/proto/gen/userpb"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/user-gateway/internal/client"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userClient *client.UserClient
}

func NewUserHandler(userClient *client.UserClient) *UserHandler {
	return &UserHandler{userClient: userClient}
}

type registerUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type updateUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/users")
	{
		api.POST("/register", h.RegisterUser)
		api.POST("/login", h.LoginUser)
		api.GET("/:id", h.GetUser)
		api.PUT("/:id", h.UpdateUser)
		api.DELETE("/:id", h.DeleteUser)
	}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req registerUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	res, err := h.userClient.RegisterUser(ctx, &userpb.RegisterUserRequest{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res.GetUser())
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req loginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	res, err := h.userClient.LoginUser(ctx, &userpb.LoginUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := withTimeout(c)
	defer cancel()

	res, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res.GetUser())
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := withTimeout(c)
	defer cancel()

	res, err := h.userClient.UpdateUser(ctx, &userpb.UpdateUserRequest{
		Id:       id,
		FullName: req.FullName,
		Email:    req.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res.GetUser())
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := withTimeout(c)
	defer cancel()

	res, err := h.userClient.DeleteUser(ctx, &userpb.DeleteUserRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func withTimeout(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), 5*time.Second)
}
