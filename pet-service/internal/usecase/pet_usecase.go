package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/pet-service/internal/entity"
	"github.com/redis/go-redis/v9"
)

type PetRepository interface {
	Create(ctx context.Context, pet entity.Pet) (entity.Pet, error)
	GetByID(ctx context.Context, id string) (entity.Pet, error)
	List(ctx context.Context, category, status string) ([]entity.Pet, error)
	Update(ctx context.Context, pet entity.Pet) (entity.Pet, error)
	Delete(ctx context.Context, id string) error
}

type PetUsecase struct {
	repo  PetRepository
	redis *redis.Client
}

func NewPetUsecase(repo PetRepository, redisClient *redis.Client) *PetUsecase {
	return &PetUsecase{
		repo:  repo,
		redis: redisClient,
	}
}

func (u *PetUsecase) CreatePet(ctx context.Context, pet entity.Pet) (entity.Pet, error) {
	if pet.Name == "" {
		return entity.Pet{}, errors.New("pet name is required")
	}
	if pet.Category == "" {
		return entity.Pet{}, errors.New("pet category is required")
	}
	if pet.Status == "" {
		pet.Status = "available"
	}

	created, err := u.repo.Create(ctx, pet)
	if err != nil {
		return entity.Pet{}, err
	}

	u.deleteListCache(ctx)

	return created, nil
}

func (u *PetUsecase) GetPet(ctx context.Context, id string) (entity.Pet, error) {
	if id == "" {
		return entity.Pet{}, errors.New("pet id is required")
	}

	cacheKey := fmt.Sprintf("pet:%s", id)

	cached, err := u.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var pet entity.Pet
		if json.Unmarshal([]byte(cached), &pet) == nil {
			return pet, nil
		}
	}

	pet, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return entity.Pet{}, err
	}

	data, _ := json.Marshal(pet)
	u.redis.Set(ctx, cacheKey, data, 5*time.Minute)

	return pet, nil
}

func (u *PetUsecase) ListPets(ctx context.Context, category, status string) ([]entity.Pet, error) {
	cacheKey := fmt.Sprintf("pets:list:%s:%s", category, status)

	cached, err := u.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var pets []entity.Pet
		if json.Unmarshal([]byte(cached), &pets) == nil {
			return pets, nil
		}
	}

	pets, err := u.repo.List(ctx, category, status)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(pets)
	u.redis.Set(ctx, cacheKey, data, 2*time.Minute)

	return pets, nil
}

func (u *PetUsecase) UpdatePet(ctx context.Context, pet entity.Pet) (entity.Pet, error) {
	if pet.ID == "" {
		return entity.Pet{}, errors.New("pet id is required")
	}
	if pet.Name == "" {
		return entity.Pet{}, errors.New("pet name is required")
	}
	if pet.Category == "" {
		return entity.Pet{}, errors.New("pet category is required")
	}
	if pet.Status == "" {
		pet.Status = "available"
	}

	updated, err := u.repo.Update(ctx, pet)
	if err != nil {
		return entity.Pet{}, err
	}

	u.redis.Del(ctx, fmt.Sprintf("pet:%s", pet.ID))
	u.deleteListCache(ctx)

	return updated, nil
}

func (u *PetUsecase) DeletePet(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("pet id is required")
	}

	if err := u.repo.Delete(ctx, id); err != nil {
		return err
	}

	u.redis.Del(ctx, fmt.Sprintf("pet:%s", id))
	u.deleteListCache(ctx)

	return nil
}

func (u *PetUsecase) deleteListCache(ctx context.Context) {
	keys, err := u.redis.Keys(ctx, "pets:list:*").Result()
	if err == nil && len(keys) > 0 {
		u.redis.Del(ctx, keys...)
	}
}
