package postgres

import (
	"context"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/user-service/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user entity.User) (entity.User, error) {
	query := `
		INSERT INTO users (full_name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, full_name, email, password_hash, created_at
	`

	var created entity.User
	err := r.db.QueryRow(
		ctx,
		query,
		user.FullName,
		user.Email,
		user.PasswordHash,
	).Scan(
		&created.ID,
		&created.FullName,
		&created.Email,
		&created.PasswordHash,
		&created.CreatedAt,
	)

	return created, err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (entity.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, created_at
		FROM users
		WHERE id = $1
	`

	var user entity.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	return user, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`

	var user entity.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	return user, err
}

func (r *UserRepository) Update(ctx context.Context, user entity.User) (entity.User, error) {
	query := `
		UPDATE users
		SET full_name = $2,
		    email = $3
		WHERE id = $1
		RETURNING id, full_name, email, password_hash, created_at
	`

	var updated entity.User
	err := r.db.QueryRow(
		ctx,
		query,
		user.ID,
		user.FullName,
		user.Email,
	).Scan(
		&updated.ID,
		&updated.FullName,
		&updated.Email,
		&updated.PasswordHash,
		&updated.CreatedAt,
	)

	return updated, err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
