package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/pet-service/internal/entity"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

type fakePetRepository struct {
	pets map[string]entity.Pet
}

func newFakePetRepository() *fakePetRepository {
	return &fakePetRepository{
		pets: make(map[string]entity.Pet),
	}
}

func (r *fakePetRepository) Create(ctx context.Context, pet entity.Pet) (entity.Pet, error) {
	pet.ID = "pet-1"
	pet.CreatedAt = time.Now()
	r.pets[pet.ID] = pet
	return pet, nil
}

func (r *fakePetRepository) GetByID(ctx context.Context, id string) (entity.Pet, error) {
	pet, ok := r.pets[id]
	if !ok {
		return entity.Pet{}, errors.New("pet not found")
	}
	return pet, nil
}

func (r *fakePetRepository) List(ctx context.Context, category, status string) ([]entity.Pet, error) {
	result := make([]entity.Pet, 0)
	for _, pet := range r.pets {
		if category != "" && pet.Category != category {
			continue
		}
		if status != "" && pet.Status != status {
			continue
		}
		result = append(result, pet)
	}
	return result, nil
}

func (r *fakePetRepository) Update(ctx context.Context, pet entity.Pet) (entity.Pet, error) {
	_, ok := r.pets[pet.ID]
	if !ok {
		return entity.Pet{}, errors.New("pet not found")
	}
	pet.CreatedAt = time.Now()
	r.pets[pet.ID] = pet
	return pet, nil
}

func (r *fakePetRepository) Delete(ctx context.Context, id string) error {
	delete(r.pets, id)
	return nil
}

func newTestRedis(t *testing.T) *redis.Client {
	t.Helper()

	server, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}

	t.Cleanup(server.Close)

	return redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})
}

func TestPetUsecaseCreatePetSuccess(t *testing.T) {
	ctx := context.Background()
	repo := newFakePetRepository()
	redisClient := newTestRedis(t)

	usecase := NewPetUsecase(repo, redisClient)

	pet, err := usecase.CreatePet(ctx, entity.Pet{
		Name:     "Buddy",
		Category: "dog",
		Breed:    "Golden Retriever",
		Age:      2,
		Price:    500,
		Status:   "available",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if pet.ID == "" {
		t.Fatal("expected pet id to be generated")
	}

	if pet.Name != "Buddy" {
		t.Fatalf("expected pet name Buddy, got %s", pet.Name)
	}
}

func TestPetUsecaseCreatePetRequiresName(t *testing.T) {
	ctx := context.Background()
	repo := newFakePetRepository()
	redisClient := newTestRedis(t)

	usecase := NewPetUsecase(repo, redisClient)

	_, err := usecase.CreatePet(ctx, entity.Pet{
		Category: "dog",
		Price:    500,
	})

	if err == nil {
		t.Fatal("expected error when pet name is empty")
	}
}

func TestPetUsecaseGetPetUsesRepository(t *testing.T) {
	ctx := context.Background()
	repo := newFakePetRepository()
	redisClient := newTestRedis(t)

	usecase := NewPetUsecase(repo, redisClient)

	created, err := usecase.CreatePet(ctx, entity.Pet{
		Name:     "Milo",
		Category: "cat",
		Price:    300,
		Status:   "available",
	})
	if err != nil {
		t.Fatalf("create pet failed: %v", err)
	}

	found, err := usecase.GetPet(ctx, created.ID)
	if err != nil {
		t.Fatalf("get pet failed: %v", err)
	}

	if found.ID != created.ID {
		t.Fatalf("expected pet id %s, got %s", created.ID, found.ID)
	}
}
