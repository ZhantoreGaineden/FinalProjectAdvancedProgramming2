package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/order-service/internal/entity"
)

type fakeOrderRepository struct {
	orders map[string]entity.Order
}

func newFakeOrderRepository() *fakeOrderRepository {
	return &fakeOrderRepository{
		orders: make(map[string]entity.Order),
	}
}

func (r *fakeOrderRepository) Create(ctx context.Context, order entity.Order) (entity.Order, error) {
	order.ID = "order-1"
	order.Status = "created"
	order.CreatedAt = time.Now()

	total := 0.0
	for i := range order.Items {
		order.Items[i].ID = "item-1"
		order.Items[i].OrderID = order.ID
		total += order.Items[i].Price
	}
	order.TotalPrice = total

	r.orders[order.ID] = order
	return order, nil
}

func (r *fakeOrderRepository) GetByID(ctx context.Context, id string) (entity.Order, error) {
	order, ok := r.orders[id]
	if !ok {
		return entity.Order{}, errors.New("order not found")
	}
	return order, nil
}

func (r *fakeOrderRepository) ListByUserID(ctx context.Context, userID string) ([]entity.Order, error) {
	result := make([]entity.Order, 0)
	for _, order := range r.orders {
		if order.UserID == userID {
			result = append(result, order)
		}
	}
	return result, nil
}

func (r *fakeOrderRepository) UpdateStatus(ctx context.Context, id, status string) (entity.Order, error) {
	order, ok := r.orders[id]
	if !ok {
		return entity.Order{}, errors.New("order not found")
	}

	order.Status = status
	r.orders[id] = order

	return order, nil
}

func (r *fakeOrderRepository) Cancel(ctx context.Context, id string) (entity.Order, error) {
	return r.UpdateStatus(ctx, id, "cancelled")
}

func TestOrderUsecaseCreateOrderSuccess(t *testing.T) {
	ctx := context.Background()
	repo := newFakeOrderRepository()
	usecase := NewOrderUsecase(repo, nil)

	order, err := usecase.CreateOrder(ctx, entity.Order{
		UserID:    "user-1",
		UserEmail: "zhantore@example.com",
		Items: []entity.OrderItem{
			{
				PetID: "pet-1",
				Price: 500,
			},
		},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if order.ID == "" {
		t.Fatal("expected order id")
	}

	if order.TotalPrice != 500 {
		t.Fatalf("expected total price 500, got %v", order.TotalPrice)
	}

	if order.Status != "created" {
		t.Fatalf("expected status created, got %s", order.Status)
	}
}

func TestOrderUsecaseCreateOrderRequiresItems(t *testing.T) {
	ctx := context.Background()
	repo := newFakeOrderRepository()
	usecase := NewOrderUsecase(repo, nil)

	_, err := usecase.CreateOrder(ctx, entity.Order{
		UserID:    "user-1",
		UserEmail: "zhantore@example.com",
	})

	if err == nil {
		t.Fatal("expected error when order items are empty")
	}
}

func TestOrderUsecaseUpdateOrderStatus(t *testing.T) {
	ctx := context.Background()
	repo := newFakeOrderRepository()
	usecase := NewOrderUsecase(repo, nil)

	order, err := usecase.CreateOrder(ctx, entity.Order{
		UserID:    "user-1",
		UserEmail: "zhantore@example.com",
		Items: []entity.OrderItem{
			{
				PetID: "pet-1",
				Price: 500,
			},
		},
	})
	if err != nil {
		t.Fatalf("create order failed: %v", err)
	}

	updated, err := usecase.UpdateOrderStatus(ctx, order.ID, "paid")
	if err != nil {
		t.Fatalf("update status failed: %v", err)
	}

	if updated.Status != "paid" {
		t.Fatalf("expected status paid, got %s", updated.Status)
	}
}

func TestOrderUsecaseCancelOrder(t *testing.T) {
	ctx := context.Background()
	repo := newFakeOrderRepository()
	usecase := NewOrderUsecase(repo, nil)

	order, err := usecase.CreateOrder(ctx, entity.Order{
		UserID:    "user-1",
		UserEmail: "zhantore@example.com",
		Items: []entity.OrderItem{
			{
				PetID: "pet-1",
				Price: 500,
			},
		},
	})
	if err != nil {
		t.Fatalf("create order failed: %v", err)
	}

	cancelled, err := usecase.CancelOrder(ctx, order.ID)
	if err != nil {
		t.Fatalf("cancel order failed: %v", err)
	}

	if cancelled.Status != "cancelled" {
		t.Fatalf("expected status cancelled, got %s", cancelled.Status)
	}
}
