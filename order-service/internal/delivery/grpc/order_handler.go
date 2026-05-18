package grpcdelivery

import (
	"context"
	"time"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/order-service/internal/entity"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/order-service/internal/usecase"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/proto/gen/orderpb"
)

type OrderHandler struct {
	orderpb.UnimplementedOrderServiceServer
	usecase *usecase.OrderUsecase
}

func NewOrderHandler(usecase *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{usecase: usecase}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.OrderResponse, error) {
	items := make([]entity.OrderItem, 0, len(req.GetItems()))

	for _, item := range req.GetItems() {
		items = append(items, entity.OrderItem{
			PetID: item.GetPetId(),
			Price: item.GetPrice(),
		})
	}

	order := entity.Order{
		UserID:    req.GetUserId(),
		UserEmail: req.GetUserEmail(),
		Status:    "created",
		Items:     items,
	}

	created, err := h.usecase.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return &orderpb.OrderResponse{Order: toProtoOrder(created)}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.OrderResponse, error) {
	order, err := h.usecase.GetOrder(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &orderpb.OrderResponse{Order: toProtoOrder(order)}, nil
}

func (h *OrderHandler) ListUserOrders(ctx context.Context, req *orderpb.ListUserOrdersRequest) (*orderpb.ListOrdersResponse, error) {
	orders, err := h.usecase.ListUserOrders(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	response := &orderpb.ListOrdersResponse{
		Orders: make([]*orderpb.Order, 0, len(orders)),
	}

	for _, order := range orders {
		response.Orders = append(response.Orders, toProtoOrder(order))
	}

	return response, nil
}

func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, req *orderpb.UpdateOrderStatusRequest) (*orderpb.OrderResponse, error) {
	order, err := h.usecase.UpdateOrderStatus(ctx, req.GetId(), req.GetStatus())
	if err != nil {
		return nil, err
	}

	return &orderpb.OrderResponse{Order: toProtoOrder(order)}, nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context, req *orderpb.CancelOrderRequest) (*orderpb.CancelOrderResponse, error) {
	_, err := h.usecase.CancelOrder(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &orderpb.CancelOrderResponse{
		Success: true,
		Message: "order cancelled successfully",
	}, nil
}

func toProtoOrder(order entity.Order) *orderpb.Order {
	createdAt := ""
	if !order.CreatedAt.IsZero() {
		createdAt = order.CreatedAt.Format(time.RFC3339)
	}

	items := make([]*orderpb.OrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, &orderpb.OrderItem{
			Id:      item.ID,
			OrderId: item.OrderID,
			PetId:   item.PetID,
			Price:   item.Price,
		})
	}

	return &orderpb.Order{
		Id:         order.ID,
		UserId:     order.UserID,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreatedAt:  createdAt,
		Items:      items,
	}
}
