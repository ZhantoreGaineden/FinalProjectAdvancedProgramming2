package usecase

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/order-service/internal/entity"
	"github.com/nats-io/nats.go"
)

type OrderRepository interface {
	Create(ctx context.Context, order entity.Order) (entity.Order, error)
	GetByID(ctx context.Context, id string) (entity.Order, error)
	ListByUserID(ctx context.Context, userID string) ([]entity.Order, error)
	UpdateStatus(ctx context.Context, id, status string) (entity.Order, error)
	Cancel(ctx context.Context, id string) (entity.Order, error)
}

type OrderUsecase struct {
	repo OrderRepository
	nats *nats.Conn
}

func NewOrderUsecase(repo OrderRepository, natsConn *nats.Conn) *OrderUsecase {
	return &OrderUsecase{
		repo: repo,
		nats: natsConn,
	}
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, order entity.Order) (entity.Order, error) {
	if order.UserID == "" {
		return entity.Order{}, errors.New("user id is required")
	}
	if order.UserEmail == "" {
		return entity.Order{}, errors.New("user email is required")
	}
	if len(order.Items) == 0 {
		return entity.Order{}, errors.New("order items are required")
	}

	created, err := u.repo.Create(ctx, order)
	if err != nil {
		return entity.Order{}, err
	}

	u.publishOrderCreated(created)

	return created, nil
}

func (u *OrderUsecase) GetOrder(ctx context.Context, id string) (entity.Order, error) {
	if id == "" {
		return entity.Order{}, errors.New("order id is required")
	}

	return u.repo.GetByID(ctx, id)
}

func (u *OrderUsecase) ListUserOrders(ctx context.Context, userID string) ([]entity.Order, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	return u.repo.ListByUserID(ctx, userID)
}

func (u *OrderUsecase) UpdateOrderStatus(ctx context.Context, id, status string) (entity.Order, error) {
	if id == "" {
		return entity.Order{}, errors.New("order id is required")
	}
	if status == "" {
		return entity.Order{}, errors.New("status is required")
	}

	updated, err := u.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		return entity.Order{}, err
	}

	u.publishOrderStatusUpdated(updated)

	return updated, nil
}

func (u *OrderUsecase) CancelOrder(ctx context.Context, id string) (entity.Order, error) {
	if id == "" {
		return entity.Order{}, errors.New("order id is required")
	}

	cancelled, err := u.repo.Cancel(ctx, id)
	if err != nil {
		return entity.Order{}, err
	}

	u.publishOrderStatusUpdated(cancelled)

	return cancelled, nil
}

func (u *OrderUsecase) publishOrderCreated(order entity.Order) {
	if u.nats == nil {
		return
	}

	event := map[string]interface{}{
		"order_id":    order.ID,
		"user_id":     order.UserID,
		"user_email":  order.UserEmail,
		"total_price": order.TotalPrice,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	_ = u.nats.Publish("order.created", data)
}

func (u *OrderUsecase) publishOrderStatusUpdated(order entity.Order) {
	if u.nats == nil {
		return
	}

	event := map[string]interface{}{
		"order_id":   order.ID,
		"user_email": order.UserEmail,
		"status":     order.Status,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	_ = u.nats.Publish("order.status_updated", data)
}
