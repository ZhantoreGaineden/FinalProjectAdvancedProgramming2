package postgres

import (
	"context"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/pet-service/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PetRepository struct {
	db *pgxpool.Pool
}

func NewPetRepository(db *pgxpool.Pool) *PetRepository {
	return &PetRepository{db: db}
}

func (r *PetRepository) Create(ctx context.Context, pet entity.Pet) (entity.Pet, error) {
	query := `
		INSERT INTO pets (name, category, breed, age, price, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, category, breed, age, price, status, created_at
	`

	var created entity.Pet
	err := r.db.QueryRow(
		ctx,
		query,
		pet.Name,
		pet.Category,
		pet.Breed,
		pet.Age,
		pet.Price,
		pet.Status,
	).Scan(
		&created.ID,
		&created.Name,
		&created.Category,
		&created.Breed,
		&created.Age,
		&created.Price,
		&created.Status,
		&created.CreatedAt,
	)

	return created, err
}

func (r *PetRepository) GetByID(ctx context.Context, id string) (entity.Pet, error) {
	query := `
		SELECT id, name, category, breed, age, price, status, created_at
		FROM pets
		WHERE id = $1
	`

	var pet entity.Pet
	err := r.db.QueryRow(ctx, query, id).Scan(
		&pet.ID,
		&pet.Name,
		&pet.Category,
		&pet.Breed,
		&pet.Age,
		&pet.Price,
		&pet.Status,
		&pet.CreatedAt,
	)

	return pet, err
}

func (r *PetRepository) List(ctx context.Context, category, status string) ([]entity.Pet, error) {
	query := `
		SELECT id, name, category, breed, age, price, status, created_at
		FROM pets
		WHERE ($1 = '' OR category = $1)
		  AND ($2 = '' OR status = $2)
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, category, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pets := make([]entity.Pet, 0)

	for rows.Next() {
		var pet entity.Pet
		if err := rows.Scan(
			&pet.ID,
			&pet.Name,
			&pet.Category,
			&pet.Breed,
			&pet.Age,
			&pet.Price,
			&pet.Status,
			&pet.CreatedAt,
		); err != nil {
			return nil, err
		}

		pets = append(pets, pet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pets, nil
}

func (r *PetRepository) Update(ctx context.Context, pet entity.Pet) (entity.Pet, error) {
	query := `
		UPDATE pets
		SET name = $2,
		    category = $3,
		    breed = $4,
		    age = $5,
		    price = $6,
		    status = $7
		WHERE id = $1
		RETURNING id, name, category, breed, age, price, status, created_at
	`

	var updated entity.Pet
	err := r.db.QueryRow(
		ctx,
		query,
		pet.ID,
		pet.Name,
		pet.Category,
		pet.Breed,
		pet.Age,
		pet.Price,
		pet.Status,
	).Scan(
		&updated.ID,
		&updated.Name,
		&updated.Category,
		&updated.Breed,
		&updated.Age,
		&updated.Price,
		&updated.Status,
		&updated.CreatedAt,
	)

	return updated, err
}

func (r *PetRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM pets WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
