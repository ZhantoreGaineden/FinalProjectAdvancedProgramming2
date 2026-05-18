package postgres

import (
	"context"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/order-service/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order entity.Order) (entity.Order, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return entity.Order{}, err
	}
	defer tx.Rollback(ctx)

	total := 0.0
	for _, item := range order.Items {
		total += item.Price
	}

	query := `
		INSERT INTO orders (user_id, user_email, total_price, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, user_email, total_price, status, created_at
	`

	if order.Status == "" {
		order.Status = "created"
	}

	var created entity.Order
	err = tx.QueryRow(
		ctx,
		query,
		order.UserID,
		order.UserEmail,
		total,
		order.Status,
	).Scan(
		&created.ID,
		&created.UserID,
		&created.UserEmail,
		&created.TotalPrice,
		&created.Status,
		&created.CreatedAt,
	)
	if err != nil {
		return entity.Order{}, err
	}

	created.Items = make([]entity.OrderItem, 0, len(order.Items))

	itemQuery := `
		INSERT INTO order_items (order_id, pet_id, price)
		VALUES ($1, $2, $3)
		RETURNING id, order_id, pet_id, price
	`

	for _, item := range order.Items {
		var createdItem entity.OrderItem
		err = tx.QueryRow(
			ctx,
			itemQuery,
			created.ID,
			item.PetID,
			item.Price,
		).Scan(
			&createdItem.ID,
			&createdItem.OrderID,
			&createdItem.PetID,
			&createdItem.Price,
		)
		if err != nil {
			return entity.Order{}, err
		}

		created.Items = append(created.Items, createdItem)
	}

	if err := tx.Commit(ctx); err != nil {
		return entity.Order{}, err
	}

	return created, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id string) (entity.Order, error) {
	query := `
		SELECT id, user_id, user_email, total_price, status, created_at
		FROM orders
		WHERE id = $1
	`

	var order entity.Order
	err := r.db.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.UserEmail,
		&order.TotalPrice,
		&order.Status,
		&order.CreatedAt,
	)
	if err != nil {
		return entity.Order{}, err
	}

	items, err := r.getItemsByOrderID(ctx, order.ID)
	if err != nil {
		return entity.Order{}, err
	}

	order.Items = items
	return order, nil
}

func (r *OrderRepository) ListByUserID(ctx context.Context, userID string) ([]entity.Order, error) {
	query := `
		SELECT id, user_id, user_email, total_price, status, created_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]entity.Order, 0)

	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.UserEmail,
			&order.TotalPrice,
			&order.Status,
			&order.CreatedAt,
		); err != nil {
			return nil, err
		}

		items, err := r.getItemsByOrderID(ctx, order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id, status string) (entity.Order, error) {
	query := `
		UPDATE orders
		SET status = $2
		WHERE id = $1
		RETURNING id, user_id, user_email, total_price, status, created_at
	`

	var order entity.Order
	err := r.db.QueryRow(ctx, query, id, status).Scan(
		&order.ID,
		&order.UserID,
		&order.UserEmail,
		&order.TotalPrice,
		&order.Status,
		&order.CreatedAt,
	)
	if err != nil {
		return entity.Order{}, err
	}

	items, err := r.getItemsByOrderID(ctx, order.ID)
	if err != nil {
		return entity.Order{}, err
	}

	order.Items = items
	return order, nil
}

func (r *OrderRepository) Cancel(ctx context.Context, id string) (entity.Order, error) {
	return r.UpdateStatus(ctx, id, "cancelled")
}

func (r *OrderRepository) getItemsByOrderID(ctx context.Context, orderID string) ([]entity.OrderItem, error) {
	query := `
		SELECT id, order_id, pet_id, price
		FROM order_items
		WHERE order_id = $1
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.OrderItem, 0)

	for rows.Next() {
		var item entity.OrderItem
		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.PetID,
			&item.Price,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
